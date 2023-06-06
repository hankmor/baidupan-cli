package cmd

import (
	"baidupan-cli/app"
	"baidupan-cli/cmd/vo"
	"baidupan-cli/util"
	"encoding/json"
	"fmt"
	"github.com/desertbit/grumble"
)

// =============================================
// 文件列表查询
// =============================================

const (
	orderByTime = "time"
	orderByName = "name"
	orderBySize = "size"
)

var fileListCmd = &grumble.Command{
	Name:    "fs",
	Aliases: []string{"files"},
	Help:    "show file lists",
	Usage:   "fs",
	Flags: func(f *grumble.Flags) {
		f.String("d", "dir", "/", "the directory to show files in it")
		f.String("o", "order", "name", `order type, support 'time','name' and 'size', default is 'name':
				1. time: sort files by file type first, then sort by modification time
				2. name: sort files by file type first, then sort by file name
				3. size: sort files by file type first, then sort by file size`)
		f.Bool("a", "asc", true, "whether to sort in ascending order")
		f.Bool("f", "only-folder", false, "whether only to query folders")
		f.Bool("e", "show-empty", false, "whether to show empty folder info")
		f.Int("l", "limit", 1000, "the number of queries, default is 1000, and it is recommended that the maximum number not exceed 1000")
		f.Bool("v", "verbose", false, "show verbose info of files")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}
		dir := ctx.Flags.String("dir")
		order := ctx.Flags.String("order")
		asc := ctx.Flags.Bool("asc")
		onlyFolder := ctx.Flags.Bool("only-folder")
		showEmpty := ctx.Flags.Bool("show-empty")
		limit := ctx.Flags.Int("limit")
		if limit <= 0 {
			return fmt.Errorf("invalid parameter: limit must be great than 0, but found %d", limit)
		}
		switch order {
		case orderByTime:
		case orderByName:
		case orderBySize:
		default:
			return fmt.Errorf("invalid parameter: unsupported order type %s", order)
		}
		req := app.ApiClient.FileinfoApi.Xpanfilelist(RootContext)
		if !asc {
			req = req.Desc(1)
		}
		if onlyFolder {
			req = req.Folder("1")
		}
		if showEmpty {
			req = req.Showempty(1)
		}
		resp, _, err := req.Limit(int32(limit)).Dir(dir).Order(order).AccessToken(*TokenResp.AccessToken).Execute()
		if err != nil {
			return err
		}

		var fileListResp vo.FileListResp
		err = json.Unmarshal([]byte(resp), &fileListResp)
		if err != nil {
			return err
		}
		for _, f := range fileListResp.Files {
			fmt.Println(util.Decode2Chinese(f.Path))
		}
		return err
	},
}
