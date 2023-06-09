package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/hankmor/gotools/conv"
	"github.com/liushuochen/gotable"
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
		table, err := gotable.Create("ID", "Account", "Nickname", "Type", "Avatar")
		if err != nil {
			return err
		}
		req := app.ApiClient.UserinfoApi.Xpannasuinfo(RootContext)
		resp, _, err := req.AccessToken(*TokenResp.AccessToken).Execute()
		table.AddRow([]string{conv.Int64ToStr(int64(*resp.Uk)), *resp.BaiduName, *resp.NetdiskName, getTypeName(*resp.VipType), *resp.AvatarUrl})
		fmt.Println(table)
		return nil
	},
}

func getTypeName(t int32) string {
	switch t {
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	default:
		return "普通用户"
	}
}
