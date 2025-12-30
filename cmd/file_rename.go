package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"io"
	pathpkg "path"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/desertbit/grumble"
)

// =============================================
// 文件/文件夹重命名
// =============================================

type fileManagerRenameItem struct {
	Path    string `json:"path"`
	Newname string `json:"newname"`
}

type fileManagerOpInfo struct {
	Errno int    `json:"errno,omitempty"`
	Path  string `json:"path,omitempty"`
}

type fileManagerOpResp struct {
	BaseVo
	RequestId uint64              `json:"request_id,omitempty"`
	TaskId    int64               `json:"taskid,omitempty"`
	Info      []fileManagerOpInfo `json:"info,omitempty"`
}

var fileRenameCmd = &grumble.Command{
	Name:    "rename",
	Aliases: []string{"rn"},
	Help:    "rename a file/folder by path",
	Usage:   "rename --path /a/b.txt --newname c.txt",
	Flags: func(f *grumble.Flags) {
		f.String("d", "dir", "", "base directory when --path is not an absolute path (optional)")
		f.String("p", "path", "", "source absolute path to rename (required)")
		f.String("n", "newname", "", "new name (base name, required)")
		f.Bool("a", "apply", false, "apply rename (default: dry-run)")
		f.Bool("A", "async", false, "whether to execute as async task (default false)")
		f.StringL("ondup", "", "duplication policy (optional, passed to openapi as-is)")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		baseDir := strings.TrimSpace(ctx.Flags.String("dir"))
		srcPath := strings.TrimSpace(ctx.Flags.String("path"))
		newname := strings.TrimSpace(ctx.Flags.String("newname"))
		if srcPath == "" {
			return fmt.Errorf("missing required flag: --path")
		}
		if newname == "" {
			return fmt.Errorf("missing required flag: --newname")
		}
		// 统一路径规则：允许不以 / 开头（自动从根目录补齐），也允许配合 --dir
		if baseDir != "" {
			var err error
			baseDir, err = cleanAbsDir(baseDir)
			if err != nil {
				return fmt.Errorf("invalid --dir: %w", err)
			}
		}
		if baseDir != "" && !strings.HasPrefix(srcPath, "/") {
			srcPath = pathpkg.Join(baseDir, srcPath)
		}
		var err error
		srcPath, err = cleanAbsPath(srcPath)
		if err != nil {
			return fmt.Errorf("invalid --path: %w", err)
		}

		filelist, dstPath, err := buildRenameFilelist(srcPath, newname)
		if err != nil {
			return err
		}

		if !ctx.Flags.Bool("apply") {
			fmt.Printf("rename plan:\n  %s -> %s\n", srcPath, dstPath)
			fmt.Println("dry-run only. add -a/--apply to execute.")
			return nil
		}

		async := int32(0)
		if ctx.Flags.Bool("async") {
			async = 1
		}

		req := app.APIClient.FilemanagerApi.Filemanagerrename(RootContext).
			AccessToken(*TokenResp.AccessToken).
			Async(async).
			Filelist(filelist)

		ondup := strings.TrimSpace(ctx.Flags.String("ondup"))
		if ondup != "" {
			req = req.Ondup(ondup)
		}

		httpResp, err := req.Execute()
		if err != nil {
			return err
		}
		defer httpResp.Body.Close()
		b, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		var resp fileManagerOpResp
		if err := sonic.Unmarshal(b, &resp); err != nil {
			return fmt.Errorf("failed to parse rename response: %w, raw=%s", err, string(b))
		}
		if !resp.Success() {
			return fmt.Errorf("rename failed, errno=%d, raw=%s", resp.Errno, string(b))
		}

		// async=1 时通常会返回 taskid
		if async == 1 && resp.TaskId != 0 {
			fmt.Printf("rename submitted: %s -> %s (taskid=%d)\n", srcPath, dstPath, resp.TaskId)
			return nil
		}

		fmt.Printf("rename success: %s -> %s\n", srcPath, dstPath)
		return nil
	},
}

func buildRenameFilelist(srcPath, newname string) (filelist string, dstPath string, err error) {
	srcPath = strings.TrimSpace(srcPath)
	newname = strings.TrimSpace(newname)
	if srcPath == "" {
		return "", "", fmt.Errorf("invalid path: empty")
	}
	if !strings.HasPrefix(srcPath, "/") {
		return "", "", fmt.Errorf("invalid path: %q, must be absolute path starting with '/'", srcPath)
	}
	if srcPath == "/" {
		return "", "", fmt.Errorf("invalid path: cannot rename root '/'")
	}
	if newname == "" {
		return "", "", fmt.Errorf("invalid newname: empty")
	}
	if strings.Contains(newname, "/") {
		return "", "", fmt.Errorf("invalid newname: must be a base name, but got %q", newname)
	}

	// 百度盘路径一般是 unix 风格，path 包更合适
	cleanSrc := strings.TrimRight(srcPath, "/")
	if cleanSrc == "" {
		cleanSrc = "/"
	}
	if cleanSrc == "/" {
		return "", "", fmt.Errorf("invalid path: cannot rename root '/'")
	}

	dstPath = pathpkg.Join(pathpkg.Dir(cleanSrc), newname)

	item := []fileManagerRenameItem{{Path: cleanSrc, Newname: newname}}
	s, e := sonic.MarshalString(item)
	if e != nil {
		return "", "", e
	}
	return s, dstPath, nil
}
