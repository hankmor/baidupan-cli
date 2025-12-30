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
		f.BoolL("async", false, "whether to execute as async task (default false)")
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
		// 兼容绝对/相对路径：
		// - srcPath 若是相对路径且提供了 baseDir，则 join 后再请求
		// - baseDir 允许不以 "/" 开头（由服务端解释为根目录相对路径）
		if !strings.HasPrefix(srcPath, "/") && baseDir != "" {
			srcPath = pathpkg.Join(baseDir, srcPath)
		}

		filelist, dstPath, err := buildRenameFilelist(srcPath, newname)
		if err != nil {
			return err
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

	dir := pathpkg.Dir(cleanSrc)
	// 相对路径场景下 path.Dir("a.txt") 会返回 "."，这里把它视为“根/当前目录”
	if dir == "." {
		dir = ""
	}
	dstPath = pathpkg.Join(dir, newname)

	item := []fileManagerRenameItem{{Path: cleanSrc, Newname: newname}}
	s, e := sonic.MarshalString(item)
	if e != nil {
		return "", "", e
	}
	return s, dstPath, nil
}
