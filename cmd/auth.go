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
		Flags: func(f *grumble.Flags) {
			f.BoolL("open-browser", true, "open qrcode url in default browser (recommended)")
		},
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
			qrURL := ""
			if authResp.QrcodeUrl != nil {
				qrURL = *authResp.QrcodeUrl
			}

			opened := false
			if c.Flags.Bool("open-browser") && qrURL != "" {
				fmt.Println("opening browser for authorization...")
				if err := util.OpenInBrowser(qrURL); err != nil {
					fmt.Println("open browser failed:", err)
				} else {
					opened = true
					fmt.Println("browser opened. please scan/confirm authorization in browser.")
				}
			}

			// 成功打开浏览器时，不再输出终端二维码（避免字体导致二维码不清晰）
			// 打开失败则回退到终端二维码输出
			if !opened {
				fmt.Println("scan qrcode to authorize cli to visit your baidupan:")
				if qrURL != "" {
					util.PrintQrCode2Console(qrURL)
				}
			}

			if qrURL != "" {
				fmt.Printf("\nauthorization url:\n%s\n", qrURL)
			}

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
			// 保存 token 到可执行文件目录，便于下次启动自动加载/刷新
			_ = util.SaveStoredToken(util.StoredToken{
				AccessToken:  derefStr(TokenResp.AccessToken),
				RefreshToken: derefStr(TokenResp.RefreshToken),
				ExpiresAt:    TokenDeadline.Unix(),
			})
			runRefreshToken()
			return e
		},
	}
)

func runRefreshToken() {
	go func() {
		for {
			// 距离过期还有一段时间就睡眠，避免 busy loop
			if TokenDeadline.IsZero() || TokenResp == nil || TokenResp.RefreshToken == nil || *TokenResp.RefreshToken == "" {
				time.Sleep(1 * time.Second)
				continue
			}
			wait := time.Until(TokenDeadline.Add(-5 * time.Second))
			if wait > 0 {
				time.Sleep(wait)
				continue
			}
			req := app.APIClient.AuthApi.OauthTokenRefreshToken(RootContext)
			resp, _, err := req.RefreshToken(*TokenResp.RefreshToken).ClientId(app.Conf.BaiduPan.AppKey).ClientSecret(app.Conf.BaiduPan.SecretKey).Execute()
			if err != nil {
				fmt.Println("refresh token error:", err)
				time.Sleep(3 * time.Second)
				continue
			}
			TokenDeadline = time.Now().Add(time.Second * time.Duration(*resp.ExpiresIn))
			TokenResp.AccessToken = resp.AccessToken
			TokenResp.RefreshToken = resp.RefreshToken
			TokenResp.ExpiresIn = resp.ExpiresIn

			_ = util.SaveStoredToken(util.StoredToken{
				AccessToken:  derefStr(TokenResp.AccessToken),
				RefreshToken: derefStr(TokenResp.RefreshToken),
				ExpiresAt:    TokenDeadline.Unix(),
			})
		}
	}()
}
