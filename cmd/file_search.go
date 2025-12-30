package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/desertbit/grumble"
)

// =============================================
// 文件搜索
// =============================================

type FileSearchResp struct {
	BaseVo
	HasMore   int     `json:"has_more,omitempty"`
	RequestId uint64  `json:"request_id,omitempty"`
	Files     []*File `json:"list,omitempty"`
}

var fileSearchCmd = &grumble.Command{
	Name:    "search",
	Aliases: []string{"find"},
	Help:    "search files/folders by keyword",
	Usage:   "search --key KEYWORD [--dir /] [--recurse] [--limit N]",
	Flags: func(f *grumble.Flags) {
		f.String("k", "key", "", "search keyword (required)")
		f.String("d", "dir", "/", "search directory (absolute path, default /)")
		f.Bool("r", "recurse", true, "search recursively (default true)")
		f.Int("l", "limit", 200, "max results to return (default 200)")
		f.Int("n", "page-size", 100, "results per page (default 100)")

		f.BoolL("only-folder", false, "only show folders")
		f.BoolL("only-files", false, "only show files")
		f.Bool("v", "verbose", false, "show verbose info")
		f.Bool("H", "human-readable", false, "show human readable size/time")
		f.Bool("g", "show-form", false, "show output in form style (requires --verbose)")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		key := strings.TrimSpace(ctx.Flags.String("key"))
		if key == "" {
			return fmt.Errorf("missing required flag: --key")
		}
		dir := strings.TrimSpace(ctx.Flags.String("dir"))
		if dir == "" {
			dir = "/"
		}
		if !strings.HasPrefix(dir, "/") {
			return fmt.Errorf("invalid --dir %q: must start with '/'", dir)
		}

		limit := ctx.Flags.Int("limit")
		if limit <= 0 {
			return fmt.Errorf("invalid --limit: %d", limit)
		}
		pageSize := ctx.Flags.Int("page-size")
		if pageSize <= 0 {
			return fmt.Errorf("invalid --page-size: %d", pageSize)
		}
		if pageSize > limit {
			pageSize = limit
		}

		onlyFolder := ctx.Flags.Bool("only-folder")
		onlyFiles := ctx.Flags.Bool("only-files")
		if onlyFolder && onlyFiles {
			return fmt.Errorf("both giving \"only-folder\" and \"only-files\" tags are ambiguous and not allowed")
		}

		rec := "0"
		if ctx.Flags.Bool("recurse") {
			rec = "1"
		}

		var out []*File
		page := 1
		for len(out) < limit {
			req := app.APIClient.FileinfoApi.Xpanfilesearch(RootContext).
				AccessToken(*TokenResp.AccessToken).
				Key(key).
				Dir(dir).
				Recursion(rec).
				Num(strconv.Itoa(pageSize)).
				Page(strconv.Itoa(page))

			respStr, _, err := req.Execute()
			if err != nil {
				return err
			}
			var resp FileSearchResp
			if err := sonic.UnmarshalString(respStr, &resp); err != nil {
				return fmt.Errorf("failed to parse search response: %w, raw=%s", err, respStr)
			}
			if !resp.Success() {
				return fmt.Errorf("error code: %d", resp.Errno)
			}

			if len(resp.Files) == 0 {
				break
			}

			for _, f := range resp.Files {
				if f == nil {
					continue
				}
				if onlyFolder && f.IsDir != 1 {
					continue
				}
				if onlyFiles && f.IsDir != 0 {
					continue
				}
				out = append(out, f)
				if len(out) >= limit {
					break
				}
			}

			// has_more: 0 no more, 1 has more
			if resp.HasMore == 0 {
				break
			}
			page++
		}

		verbose := ctx.Flags.Bool("verbose")
		humanReadable := ctx.Flags.Bool("human-readable")
		showForm := ctx.Flags.Bool("show-form")
		return (&SimpleFileLister{}).Print(out, FilePrinterOption{Verbose: verbose, HumanReadable: humanReadable, ShowForm: showForm})
	},
}
