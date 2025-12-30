package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"io"
	pathpkg "path"
	"strings"

	"github.com/desertbit/grumble"
)

// =============================================
// 移动（mv/move）
// =============================================

var fileMoveCmd = &grumble.Command{
	Name:    "mv",
	Aliases: []string{"move"},
	Help:    "move file(s)/folder(s) to a destination directory (default: dry-run; use -a/--apply to execute)",
	Usage:   "mv SRC DEST  |  mv SRC1 SRC2... DESTDIR",
	Args: func(a *grumble.Args) {
		a.StringList("args", "SRC... DEST (path auto starts from '/')", grumble.Default([]string{}))
	},
	Flags: func(f *grumble.Flags) {
		f.Bool("a", "apply", false, "apply move (default: dry-run)")
		f.Bool("A", "async", false, "submit as async task")
		f.Int("s", "size", 100, "max items per request (default 100)")
		f.StringL("ondup", "", "duplication policy (optional, passed to openapi as-is)")
		f.Bool("p", "progress", true, "show progress when executing (default true)")
		f.Bool("c", "continue-on-error", false, "continue processing remaining chunks when error happens (default false)")
		f.Bool("i", "ignore-errors", false, "exit with success even if some items failed (only meaningful with --continue-on-error)")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		rawArgs := ctx.Args.StringList("args")
		if len(rawArgs) == 0 {
			return fmt.Errorf("missing arguments: mv SRC DEST  |  mv SRC1 SRC2... DESTDIR")
		}

		var (
			srcArgs []string
			destDir string
			newname string
		)

		if len(rawArgs) < 2 {
			return fmt.Errorf("missing DEST: use `mv SRC DEST` or provide multiple SRC then DESTDIR")
		}
		destArg := rawArgs[len(rawArgs)-1]
		srcArgs = rawArgs[:len(rawArgs)-1]

		// 多源：DEST 必须是目录
		if len(srcArgs) > 1 {
			d, err := cleanAbsDir(destArg)
			if err != nil {
				return fmt.Errorf("invalid DESTDIR: %w", err)
			}
			destDir = d
		} else {
			// 单源：DEST 既可以是目录，也可以是文件路径（重命名）
			if strings.HasSuffix(strings.TrimSpace(destArg), "/") || strings.TrimSpace(destArg) == "/" {
				d, err := cleanAbsDir(destArg)
				if err != nil {
					return fmt.Errorf("invalid DESTDIR: %w", err)
				}
				destDir = d
			} else {
				dstFull, err := cleanAbsPath(destArg)
				if err != nil {
					return fmt.Errorf("invalid DEST: %w", err)
				}
				destDir = pathpkg.Dir(dstFull)
				newname = pathpkg.Base(dstFull)
			}
		}

		items := make([]fileManagerCopyMoveItem, 0, len(srcArgs))
		planLines := make([]string, 0, len(srcArgs))
		for _, p := range srcArgs {
			src, err := cleanAbsPath(p)
			if err != nil {
				return err
			}
			item := fileManagerCopyMoveItem{
				Path: src,
				Dest: destDir,
			}
			dstName := pathpkg.Base(src)
			if newname != "" && len(srcArgs) == 1 {
				item.Newname = newname
				dstName = newname
			}
			items = append(items, item)
			planLines = append(planLines, fmt.Sprintf("%s -> %s", src, pathpkg.Join(destDir, dstName)))
		}

		apply := ctx.Flags.Bool("apply")
		if !apply {
			fmt.Printf("move plan (%d item(s)):\n", len(planLines))
			for _, line := range planLines {
				fmt.Println("  " + line)
			}
			fmt.Println("dry-run only. add -a/--apply to execute.")
			return nil
		}

		async := int32(0)
		if ctx.Flags.Bool("async") {
			async = 1
		}

		ondup := strings.TrimSpace(ctx.Flags.String("ondup"))
		chunkSize := ctx.Flags.Int("size")
		showProgress := ctx.Flags.Bool("progress")
		continueOnError := ctx.Flags.Bool("continue-on-error")
		ignoreErrors := ctx.Flags.Bool("ignore-errors")

		return applyFileManagerChunks(
			"move",
			items,
			chunkSize,
			showProgress,
			continueOnError,
			ignoreErrors,
			func(filelist string) ([]byte, error) {
				req := app.APIClient.FilemanagerApi.Filemanagermove(RootContext).
					AccessToken(*TokenResp.AccessToken).
					Async(async).
					Filelist(filelist)
				if ondup != "" {
					req = req.Ondup(ondup)
				}
				httpResp, err := req.Execute()
				if err != nil {
					return nil, err
				}
				defer httpResp.Body.Close()
				return io.ReadAll(httpResp.Body)
			},
		)
	},
}
