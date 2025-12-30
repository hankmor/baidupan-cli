package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"io"

	"github.com/desertbit/grumble"
)

// =============================================
// 删除（rm/del/delete）
// =============================================

type fileManagerDeleteItem struct {
	Path string `json:"path"`
}

var fileDeleteCmd = &grumble.Command{
	Name:    "rm",
	Aliases: []string{"del", "delete"},
	Help:    "delete file(s)/folder(s) by path (default: dry-run; use -a/--apply to execute)",
	Usage:   "rm PATH...            (dry-run)\nrm -a PATH...         (apply deletion)",
	Args: func(a *grumble.Args) {
		a.StringList("paths", "path(s) to delete (path auto starts from '/')", grumble.Default([]string{}))
	},
	Flags: func(f *grumble.Flags) {
		f.Bool("r", "recursive", false, "compat flag: delete directory recursively (baidupan deletes directories by default)")
		f.Bool("a", "apply", false, "apply deletion (default: dry-run)")
		f.Bool("A", "async", false, "submit as async task")
		f.Int("s", "size", 100, "max items per request (default 100)")
		f.Bool("p", "progress", true, "show progress when executing (default true)")
		f.Bool("c", "continue-on-error", false, "continue processing remaining chunks when error happens (default false)")
		f.Bool("i", "ignore-errors", false, "exit with success even if some items failed (only meaningful with --continue-on-error)")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		paths := ctx.Args.StringList("paths")
		if len(paths) == 0 {
			return fmt.Errorf("missing PATH(s)")
		}

		items := make([]fileManagerDeleteItem, 0, len(paths))
		for _, p := range paths {
			src, err := cleanAbsPath(p)
			if err != nil {
				return err
			}
			items = append(items, fileManagerDeleteItem{Path: src})
		}

		_ = ctx.Flags.Bool("recursive") // compat only

		apply := ctx.Flags.Bool("apply")
		if !apply {
			fmt.Printf("delete plan (%d item(s)):\n", len(items))
			for _, it := range items {
				fmt.Println("  " + it.Path)
			}
			if !apply {
				fmt.Println("dry-run only. add -a/--apply to execute.")
			}
			return nil
		}

		async := int32(0)
		if ctx.Flags.Bool("async") {
			async = 1
		}

		chunkSize := ctx.Flags.Int("size")
		showProgress := ctx.Flags.Bool("progress")
		continueOnError := ctx.Flags.Bool("continue-on-error")
		ignoreErrors := ctx.Flags.Bool("ignore-errors")

		return applyFileManagerChunks(
			"delete",
			items,
			chunkSize,
			showProgress,
			continueOnError,
			ignoreErrors,
			func(filelist string) ([]byte, error) {
				req := app.APIClient.FilemanagerApi.Filemanagerdelete(RootContext).
					AccessToken(*TokenResp.AccessToken).
					Async(async).
					Filelist(filelist)
				// delete API also has ondup in generator, but it's not meaningful; keep it unused.
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
