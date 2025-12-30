package cmd

import (
	"baidupan-cli/app"
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"

	"baidupan-cli/util"

	"github.com/bytedance/sonic"
	"github.com/desertbit/grumble"
)

// =============================================
// 批量重命名（正则匹配）
// =============================================

type batchRenamePlanItem struct {
	OldPath string
	OldName string
	NewName string
	NewPath string
	IsDir   bool
}

type renameBatchTarget string

const (
	renameBatchTargetFiles renameBatchTarget = "files"
	renameBatchTargetDirs  renameBatchTarget = "dirs"
	renameBatchTargetAll   renameBatchTarget = "all"
)

var fileRenameBatchCmd = &grumble.Command{
	Name:    "rename-batch",
	Aliases: []string{"rnb", "rb"},
	Help:    "batch rename files/folders in a directory (default: sed-like replace mode)",
	Usage:   "rb --dir /video FIND TO  (or: rb --dir /video --pattern '^(.*)\\.mp4$' --replace '${1}_1080p.mp4')",
	Args: func(a *grumble.Args) {
		// sed-like default: FIND TO
		a.String("find", "find string (positional; optional if --find is set)", grumble.Default(""))
		a.String("to", "replace-to string (positional; optional if --to is set)", grumble.Default(""))
	},
	Flags: func(f *grumble.Flags) {
		f.String("d", "dir", "/", "target directory to scan")
		f.Bool("r", "recurse", false, "whether to scan directory recursively")
		f.StringL("target", "files", "rename target: files|dirs|all (default files)")
		f.Int("l", "limit", 1000, "max number of entries to scan (default 1000)")
		// Full regex mode:
		f.StringL("pattern", "", "regex pattern to match file/folder name (optional)")
		f.StringL("replace", "", "replacement string (supports $1..$n) (optional)")
		// Simple replace mode:
		f.StringL("find", "", "find substring to replace (optional). When set, no need to write full regex pattern")
		f.StringL("to", "", "replacement string for --find (optional)")
		f.BoolL("find-regex", false, "treat --find as regex (default false: literal substring)")

		f.Bool("a", "apply", false, "apply changes (default: dry-run)")
		f.BoolL("async", false, "submit rename as async task")
		f.Int("s", "size", 100, "max items per request when applying (default 100)")
		f.StringL("ondup", "", "duplication policy (optional, passed to openapi as-is)")
		f.Bool("p", "progress", true, "show progress/spinner when applying (default true)")
		f.Bool("c", "continue-on-error", false, "continue processing remaining chunks when error happens (default false)")
		f.Bool("i", "ignore-errors", false, "exit with success even if some items failed (only meaningful with --continue-on-error)")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		dir := strings.TrimSpace(ctx.Flags.String("dir"))
		if dir == "" {
			dir = "/"
		}
		if !strings.HasPrefix(dir, "/") {
			return fmt.Errorf("invalid --dir %q: must start with '/'", dir)
		}
		// prefer flags; fallback to positional args (sed-like)
		find := ctx.Flags.String("find")
		to := ctx.Flags.String("to")
		if strings.TrimSpace(find) == "" {
			find = ctx.Args.String("find")
		}
		if to == "" {
			to = ctx.Args.String("to")
		}

		re, replace, err := buildBatchRenameMatcher(
			ctx.Flags.String("pattern"),
			ctx.Flags.String("replace"),
			find,
			to,
			ctx.Flags.Bool("find-regex"),
		)
		if err != nil {
			return err
		}

		limit := ctx.Flags.Int("limit")
		if limit <= 0 {
			return fmt.Errorf("invalid --limit: %d", limit)
		}

		files, err := listFilesForBatchRename(dir, ctx.Flags.Bool("recurse"), int32(limit))
		if err != nil {
			return err
		}

		target, err := parseRenameBatchTarget(ctx.Flags.String("target"))
		if err != nil {
			return err
		}
		_, plan, err := buildBatchRenamePlan(files, re, replace, target)
		if err != nil {
			return err
		}
		if len(plan) == 0 {
			fmt.Println("no rename candidates matched (or all replacements result in no change).")
			return nil
		}

		// 输出计划（稳定排序）
		sort.Slice(plan, func(i, j int) bool { return plan[i].OldPath < plan[j].OldPath })
		fmt.Printf("matched %d item(s):\n", len(plan))
		for _, p := range plan {
			kind := "FILE"
			if p.IsDir {
				kind = "DIR "
			}
			fmt.Printf("  [%s] %s  ->  %s\n", kind, p.OldPath, p.NewPath)
		}

		if !ctx.Flags.Bool("apply") {
			fmt.Println("\n(dry-run) add --apply to execute.")
			return nil
		}

		async := int32(0)
		if ctx.Flags.Bool("async") {
			async = 1
		}
		showProgress := ctx.Flags.Bool("progress")
		continueOnError := ctx.Flags.Bool("continue-on-error")
		ignoreErrors := ctx.Flags.Bool("ignore-errors")
		if ignoreErrors && !continueOnError {
			return fmt.Errorf("--ignore-errors requires --continue-on-error")
		}
		chunkSize := ctx.Flags.Int("size")
		if chunkSize <= 0 {
			return fmt.Errorf("invalid --size: %d", chunkSize)
		}
		ondup := strings.TrimSpace(ctx.Flags.String("ondup"))

		// Apply order matters when directories are involved:
		// - rename files first (so their old paths are still valid)
		// - then rename directories from deep to shallow (avoid parent rename breaking child paths)
		applyPlan := make([]batchRenamePlanItem, len(plan))
		copy(applyPlan, plan)
		sort.Slice(applyPlan, func(i, j int) bool {
			a, b := applyPlan[i], applyPlan[j]
			if a.IsDir != b.IsDir {
				return !a.IsDir // files first
			}
			da := pathDepth(a.OldPath)
			db := pathDepth(b.OldPath)
			if a.IsDir && da != db {
				return da > db // dirs: deep first
			}
			// stable-ish
			return a.OldPath < b.OldPath
		})
		applyItems := make([]fileManagerRenameItem, 0, len(applyPlan))
		for _, p := range applyPlan {
			applyItems = append(applyItems, fileManagerRenameItem{
				Path:    p.OldPath,
				Newname: p.NewName,
			})
		}

		total := len(applyItems)
		var (
			appliedChunks int
			failedChunks  int
			failedItems   int
		)
		for i := 0; i < total; i += chunkSize {
			end := i + chunkSize
			if end > total {
				end = total
			}
			chunk := applyItems[i:end]

			if showProgress {
				fmt.Printf("\n[%d/%d] applying %d item(s)...\n", i+1, total, len(chunk))
				fmt.Printf("  (hint) if request times out, try smaller --size or enable --async\n")
			}

			filelist, err := sonic.MarshalString(chunk)
			if err != nil {
				return err
			}
			req := app.APIClient.FilemanagerApi.Filemanagerrename(context.Background()).
				AccessToken(*TokenResp.AccessToken).
				Async(async).
				Filelist(filelist)
			if ondup != "" {
				req = req.Ondup(ondup)
			}

			var closeSpin chan struct{}
			if showProgress {
				closeSpin = make(chan struct{})
				util.Spin("processing", closeSpin)
			}
			start := time.Now()
			httpResp, err := req.Execute()
			if showProgress {
				close(closeSpin)
				fmt.Print("\r") // clear spinner line
			}
			if err != nil {
				failedChunks++
				if !continueOnError {
					return err
				}
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, err)
				continue
			}
			b, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				failedChunks++
				if !continueOnError {
					return err
				}
				fmt.Printf("chunk %d-%d/%d read response failed: %v\n", i+1, end, total, err)
				continue
			}

			var resp fileManagerOpResp
			if err := sonic.Unmarshal(b, &resp); err != nil {
				failedChunks++
				e := fmt.Errorf("failed to parse rename response: %w, raw=%s", err, string(b))
				if !continueOnError {
					return e
				}
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, e)
				continue
			}
			if !resp.Success() {
				// errno=12 常见为“部分失败”，errno!=0 但仍可能有部分成功
				itemFails := 0
				for _, info := range resp.Info {
					if info.Errno != 0 {
						itemFails++
						fmt.Printf("  item failed: errno=%d path=%s\n", info.Errno, info.Path)
					}
				}
				if itemFails == 0 {
					itemFails = len(chunk)
				}
				failedItems += itemFails
				failedChunks++

				e := fmt.Errorf("batch rename chunk failed, errno=%d request_id=%d", resp.Errno, resp.RequestId)
				if !continueOnError {
					return fmt.Errorf("%w, raw=%s", e, string(b))
				}
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, e)
				continue
			}
			if async == 1 && resp.TaskId != 0 {
				fmt.Printf("submitted chunk %d-%d/%d (taskid=%d, cost=%s)\n", i+1, end, total, resp.TaskId, time.Since(start).Truncate(10*time.Millisecond))
			} else {
				fmt.Printf("applied chunk %d-%d/%d (cost=%s)\n", i+1, end, total, time.Since(start).Truncate(10*time.Millisecond))
			}
			appliedChunks++
		}

		if continueOnError && (failedChunks > 0 || failedItems > 0) {
			msg := fmt.Sprintf("completed with failures: applied_chunks=%d failed_chunks=%d failed_items~=%d", appliedChunks, failedChunks, failedItems)
			if ignoreErrors {
				fmt.Println(msg)
				return nil
			}
			return fmt.Errorf(msg)
		}

		return nil
	},
}

func pathDepth(p string) int {
	p = strings.TrimSpace(p)
	p = strings.TrimRight(p, "/")
	if p == "" {
		return 0
	}
	return strings.Count(p, "/")
}

func listFilesForBatchRename(dir string, recurse bool, limit int32) ([]*File, error) {
	// 复用现有逻辑（同包可直接用）
	options := NewFileListOptions().Limit(limit)
	var lister FileLister
	if recurse {
		lister = &RecursionFileLister{}
	} else {
		lister = &SimpleFileLister{}
	}
	return lister.List(dir, *options)
}

func buildBatchRenamePlan(files []*File, re *regexp.Regexp, replace string, target renameBatchTarget) ([]fileManagerRenameItem, []batchRenamePlanItem, error) {
	replace = normalizeRegexReplace(replace)
	var items []fileManagerRenameItem
	var plan []batchRenamePlanItem

	// 冲突检测：newPath -> oldPath
	seen := map[string]string{}

	for _, f := range files {
		if f == nil {
			continue
		}
		isDir := f.IsDir == 1
		switch target {
		case renameBatchTargetFiles:
			if isDir {
				continue
			}
		case renameBatchTargetDirs:
			if !isDir {
				continue
			}
		case renameBatchTargetAll:
			// keep all
		default:
			return nil, nil, fmt.Errorf("invalid target: %q", target)
		}

		oldName := f.ServerFilename
		if oldName == "" {
			continue
		}
		if !re.MatchString(oldName) {
			continue
		}
		newName := re.ReplaceAllString(oldName, replace)
		newName = strings.TrimSpace(newName)
		if newName == "" {
			return nil, nil, fmt.Errorf("replacement produced empty name for %q", oldName)
		}
		if strings.Contains(newName, "/") {
			return nil, nil, fmt.Errorf("replacement produced invalid name (contains '/'): %q", newName)
		}
		if newName == oldName {
			continue
		}

		oldPath := strings.TrimSpace(f.Path)
		if oldPath == "" {
			// fallback：尽量拼出来（一般不会走到）
			oldPath = oldName
		}
		newPath := oldPath
		if strings.HasSuffix(oldPath, oldName) {
			newPath = strings.TrimSuffix(oldPath, oldName) + newName
		}

		if prev, ok := seen[newPath]; ok && prev != oldPath {
			return nil, nil, fmt.Errorf("rename conflict: %q and %q would both become %q", prev, oldPath, newPath)
		}
		seen[newPath] = oldPath

		items = append(items, fileManagerRenameItem{
			Path:    oldPath,
			Newname: newName,
		})
		plan = append(plan, batchRenamePlanItem{
			OldPath: oldPath,
			OldName: oldName,
			NewName: newName,
			NewPath: newPath,
			IsDir:   isDir,
		})
	}
	return items, plan, nil
}

func parseRenameBatchTarget(s string) (renameBatchTarget, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", string(renameBatchTargetFiles):
		return renameBatchTargetFiles, nil
	case string(renameBatchTargetDirs):
		return renameBatchTargetDirs, nil
	case string(renameBatchTargetAll):
		return renameBatchTargetAll, nil
	default:
		return "", fmt.Errorf("invalid --target %q: must be one of files|dirs|all", s)
	}
}

// normalizeRegexReplace fixes a common Go regexp replacement pitfall:
// In Go, `$1_foo` is treated as a named group `$1_foo` (often empty),
// while users usually mean `${1}_foo`. We rewrite `$<digits><letter|_>` into `${<digits>}<letter|_>`.
func normalizeRegexReplace(repl string) string {
	repl = strings.TrimSpace(repl)
	if repl == "" {
		return repl
	}
	amb := regexp.MustCompile(`\$(\d+)([A-Za-z_])`)
	return amb.ReplaceAllString(repl, `${$1}$2`)
}

func buildBatchRenameMatcher(pattern, replace, find, to string, findRegex bool) (*regexp.Regexp, string, error) {
	pattern = strings.TrimSpace(pattern)
	replace = normalizeRegexReplace(replace)
	find = strings.TrimSpace(find)
	// keep `to` raw (allow spaces)

	// Full regex mode
	if pattern != "" || replace != "" {
		if pattern == "" {
			return nil, "", fmt.Errorf("missing --pattern (when using --replace)")
		}
		if replace == "" {
			return nil, "", fmt.Errorf("missing --replace (when using --pattern)")
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, "", fmt.Errorf("invalid --pattern: %w", err)
		}
		return re, replace, nil
	}

	// Simple replace mode
	if find == "" && to == "" {
		return nil, "", fmt.Errorf("missing matcher: provide either (--pattern and --replace) or (FIND TO) or (--find and --to)")
	}
	if find == "" {
		return nil, "", fmt.Errorf("missing FIND (positional) or --find")
	}
	if to == "" {
		return nil, "", fmt.Errorf("missing TO (positional) or --to")
	}

	var re *regexp.Regexp
	var err error
	if findRegex {
		re, err = regexp.Compile(find)
		if err != nil {
			return nil, "", fmt.Errorf("invalid --find regex: %w", err)
		}
	} else {
		re = regexp.MustCompile(regexp.QuoteMeta(find))
	}
	return re, to, nil
}
