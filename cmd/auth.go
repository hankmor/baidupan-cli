package cmd

import (
	"baidupan-cli/app"
	openapi "baidupan-cli/openxpanapi"
	"baidupan-cli/util"
	"context"
	"fmt"
	"github.com/desertbit/grumble"
)

var AuthResp *openapi.OauthTokenDeviceCodeResponse

var loginCmd = &grumble.Command{
	Name:     "auth",
	Help:     "authorize cli to visit your baidupan account",
	LongHelp: "scan the given qrcode to authorize cli to visit your baidupan account",
	Run: func(c *grumble.Context) error {
		if AuthResp != nil {
			return fmt.Errorf("already authorized, please do not re-authorize")
		}

		// 使用设备码授权
		fmt.Println("generating qrcode...")
		ctx := context.Background()
		authReq := app.ApiClient.AuthApi.OauthTokenDeviceCode(ctx)
		authResp, _, err := authReq.ClientId(app.Conf.BaiduPan.AppKey).Scope("basic,netdisk").Execute()
		if err != nil {
			return err
		}
		fmt.Println("scan qrcode to authorize cli to visit your baidupan:")
		util.PrintQrCode2Console(*authResp.QrcodeUrl)

		// 轮询获取 accesstoken
		fmt.Println()
		fmt.Println("waiting for authorizing...")
		tokenReq := app.ApiClient.AuthApi.OauthTokenDeviceToken(ctx)
		tokenResp, _, err := tokenReq.Execute()
		if err != nil {
			return err
		}
		fmt.Println(tokenResp)
		return err
	},
}

func init() {
	app.RegisterCommand(loginCmd)
}

type authInfo struct {
}
