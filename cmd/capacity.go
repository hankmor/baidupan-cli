package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"github.com/desertbit/grumble"
)

// =============================================
// 网盘容量查询
// =============================================

var capCmd = &grumble.Command{
	Name:     "cap",
	Help:     "show capacity",
	LongHelp: "show capacity of your baidupan",
	Flags: func(f *grumble.Flags) {
		f.Bool("e", "expire", false, "show expire information")
		f.Bool("f", "free", false, "show free information")
	},
	Run: func(c *grumble.Context) error {
		if err := checkAuthorized(c); err != nil {
			return err
		}
		req := app.ApiClient.UserinfoApi.Apiquota(RootContext)
		if c.Flags.Bool("expire") {
			req = req.Checkexpire(1)
		}
		if c.Flags.Bool("free") {
			req = req.Checkfree(1)
		}
		resp, _, err := req.AccessToken(*TokenResp.AccessToken).Execute()
		b, err := resp.MarshalJSON()
		fmt.Println(string(b))
		return err
	},
}
