package cmd

import (
	"baidupan-cli/app"
	"baidupan-cli/util"
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/liushuochen/gotable"
)

// =============================================
// 网盘容量查询
// =============================================

var capCmd = &grumble.Command{
	Name:     "cap",
	Help:     "show capacity",
	LongHelp: "show capacity of your baidupan",
	Flags: func(f *grumble.Flags) {
		f.Bool("e", "expire", false, "whether to show expire information")
		f.Bool("f", "free", false, "whether to show free information")
		f.Bool("H", "human-readable", false, "whether to show information as human readable")
	},
	Run: func(c *grumble.Context) error {
		if err := checkAuthorized(c); err != nil {
			return err
		}

		table, err := gotable.Create("Total", "Used", "Free", "Expire In 7 Days")
		if err != nil {
			return err
		}
		req := app.APIClient.UserinfoApi.Apiquota(RootContext)
		humanReadable := c.Flags.Bool("human-readable")

		if c.Flags.Bool("expire") {
			req = req.Checkexpire(1)
		}
		if c.Flags.Bool("free") {
			req = req.Checkfree(1)
		}
		resp, _, err := req.AccessToken(*TokenResp.AccessToken).Execute()
		free := resp.Free
		expire := resp.Expire
		var totalstr, usedstr, freestr, expirestr string
		if free != nil {
			freestr = util.Int64ToStr(*resp.Free)
			if humanReadable {
				freestr = util.ConvReadableSize(*resp.Free)
			}
		}
		if expire != nil {
			expirestr = getExpireIn7Days(*expire)
		}
		if humanReadable {
			totalstr = util.ConvReadableSize(*resp.Total)
			usedstr = util.ConvReadableSize(*resp.Used)
		} else {
			totalstr = util.Int64ToStr(*resp.Total)
			usedstr = util.Int64ToStr(*resp.Used)
		}
		table.AddRow([]string{totalstr, usedstr, freestr, expirestr})
		fmt.Println(table)
		return err
	},
}

func getExpireIn7Days(b bool) string {
	if b {
		return Yes
	}
	return No
}
