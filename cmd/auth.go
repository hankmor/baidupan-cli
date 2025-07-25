// Package cmd show many supported Commands.
package cmd

import (
	"baidupan-cli/app"
	openapi "baidupan-cli/openxpanapi"
	"baidupan-cli/util"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/desertbit/grumble"
)

// =============================================
// 网盘访问授权
// =============================================

var (
	AuthResp      *openapi.OauthTokenDeviceCodeResponse
	TokenResp     *openapi.OauthTokenDeviceTokenResponse
	TokenDeadline time.Time
	RootContext   = context.Background()
	authCmd       = &grumble.Command{
		Name:     "auth",
		Help:     "authorize cli to visit your baidupan account",
		LongHelp: "scan the given qrcode to authorize cli to visit your baidupan account",
		Run: func(c *grumble.Context) error {
			if AuthResp != nil {
				return fmt.Errorf("already authorized, please do not re-authorize")
			}

			// 使用设备码授权
			fmt.Println("generating qrcode...")
			authReq := app.APIClient.AuthApi.OauthTokenDeviceCode(RootContext)
			authResp, _, err := authReq.
				ClientId(app.Conf.BaiduPan.AppKey).
				Scope("basic,netdisk").
				Execute()
			if err != nil {
				return err
			}
			AuthResp = &authResp
			fmt.Println("scan qrcode to authorize cli to visit your baidupan:")
			util.PrintQrCode2Console(*authResp.QrcodeUrl)

			interval := authResp.Interval
			expireIn := authResp.ExpiresIn
			deadline := time.Now().Add(time.Second * time.Duration(*expireIn-5)) // 5秒的冗余时间

			// 轮询获取 accesstoken
			fmt.Println()
			closeSpin := make(chan struct{})
			var e error
			util.Spin("waiting for authorizing...", closeSpin)
			for {
				tokenReq := app.APIClient.AuthApi.OauthTokenDeviceToken(RootContext)
				tokenResp, tokenHttpResp, err := tokenReq.
					Code(*authResp.DeviceCode).
					ClientId(app.Conf.BaiduPan.AppKey).
					ClientSecret(app.Conf.BaiduPan.SecretKey).
					Execute()
				if err != nil {
					// 400 时是等待授权
					if tokenHttpResp.StatusCode != http.StatusBadRequest {
						e = err
						break
					}
				}
				if tokenResp.AccessToken != nil {
					TokenResp = &tokenResp
					break
				}
				time.Sleep(time.Second * time.Duration(*interval))
				if time.Now().After(deadline) {
					e = fmt.Errorf("authrization expired, try it agagin")
					break
				}
			}
			close(closeSpin)
			fmt.Println("\nauthorize success!")
			TokenDeadline = time.Now().Add(time.Second * time.Duration(*TokenResp.ExpiresIn))
			runRefreshToken()
			return e
		},
	}
)

func runRefreshToken() {
	go func() {
		for {
			// 未过期，检测时间有 5 秒的冗余时间
			if !time.Now().Add(time.Second * 5).After(TokenDeadline) {
				continue
			}
			req := app.APIClient.AuthApi.OauthTokenRefreshToken(RootContext)
			resp, _, err := req.RefreshToken(*TokenResp.RefreshToken).ClientId(app.Conf.BaiduPan.AppKey).ClientSecret(app.Conf.BaiduPan.SecretKey).Execute()
			if err != nil {
				fmt.Println("refresh token error:", err)
				continue
			}
			TokenDeadline = time.Now().Add(time.Second * time.Duration(*resp.ExpiresIn))
			TokenResp.AccessToken = resp.AccessToken
			TokenResp.RefreshToken = resp.RefreshToken
			TokenResp.ExpiresIn = resp.ExpiresIn
			time.Sleep(time.Second * 1)
		}
	}()
}
