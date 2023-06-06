package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"github.com/desertbit/grumble"
)

// =============================================
// 用户信息查询
// =============================================

var userInfoCmd = &grumble.Command{
	Name:     "userinfo",
	Help:     "show user info",
	LongHelp: "show user info of your baidupan account",
	Run: func(c *grumble.Context) error {
		if err := checkAuthorized(c); err != nil {
			return err
		}
		req := app.ApiClient.UserinfoApi.Xpannasuinfo(RootContext)
		resp, _, err := req.AccessToken(*TokenResp.AccessToken).Execute()
		b, err := resp.MarshalJSON()
		fmt.Println(string(b))
		return err
	},
}
