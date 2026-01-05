package cmd

import (
	"baidupan-cli/app"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/desertbit/grumble"
)

var cdCmd = &grumble.Command{
	Name:  "cd",
	Help:  "change current directory",
	Usage: "cd [PATH]",
	Args: func(a *grumble.Args) {
		a.String("path", "the target directory path", grumble.Default("/"))
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		target := ctx.Args.String("path")

		target = ResolvePath(target)

		// Verify directory exists
		if target != "/" {
			req := app.APIClient.FileinfoApi.Xpanfilelist(RootContext)
			// checking if dir exists by listing it with limit=1
			resp, _, err := req.Dir(target).Limit(1).AccessToken(*TokenResp.AccessToken).Execute()
			if err != nil {
				return err
			}

			var fileListResp FileListResp
			err = sonic.UnmarshalString(resp, &fileListResp)
			if err != nil {
				return err
			}

			if !fileListResp.Success() {
				// If errno is not 0, it likely means directory doesn't exist or other error
				return fmt.Errorf("cannot cd to %s: error code %d (maybe directory does not exist)", target, fileListResp.Errno)
			}
		}

		app.CurrentDir = target
		updatePrompt()
		return nil
	},
}

func updatePrompt() {
	p := app.CurrentDir
	if p == "/" {
		app.App.SetPrompt("cli » ")
	} else {
		app.App.SetPrompt(fmt.Sprintf("cli %s » ", p))
	}
}
