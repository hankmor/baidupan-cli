package cmd

import (
	"baidupan-cli/app"
	openapi "baidupan-cli/openxpanapi"
	"baidupan-cli/util"
	"context"
	"fmt"
	"time"

	"github.com/desertbit/grumble"
)

// =============================================
// 初始化
// =============================================

func init() {
	app.RegisterCommand(authCmd)
	app.RegisterCommand(capCmd)
	app.RegisterCommand(userInfoCmd)
	app.RegisterCommand(fileListCmd)
	app.RegisterCommand(fileSearchCmd)
	app.RegisterCommand(fileRenameCmd)
	app.RegisterCommand(fileRenameBatchCmd)
	app.RegisterCommand(fileCopyCmd)
	app.RegisterCommand(fileMoveCmd)
	app.RegisterCommand(fileDeleteCmd)
	app.RegisterCommand(cdCmd)

	// 启动时自动加载 token，并在必要时刷新（不与 app 包形成循环依赖）
	app.RegisterInitHook(func(a *grumble.App, flags grumble.FlagMap) error {
		// debug 显式 token 优先
		if t := flags.String("access-token"); t != "" {
			TokenResp = &openapi.OauthTokenDeviceTokenResponse{AccessToken: &t}
			return nil
		}

		st, err := util.LoadStoredToken()
		if err != nil {
			return err
		}
		if st == nil {
			return nil
		}

		// 先用存量 token
		at := st.AccessToken
		TokenResp = &openapi.OauthTokenDeviceTokenResponse{AccessToken: &at, RefreshToken: &st.RefreshToken}
		if st.ExpiresAt > 0 {
			TokenDeadline = time.Unix(st.ExpiresAt, 0)
		}


		// 未过期则直接使用
		if !st.Expired(10 * time.Second) {
			// 启动后台刷新（只要有 refresh_token）
			if st.RefreshToken != "" {
				runRefreshToken()
			}
			return nil
		}

		// 过期则尝试刷新
		if st.RefreshToken == "" {
			return nil
		}
		if app.Conf == nil {
			return fmt.Errorf("config not loaded")
		}
		req := app.APIClient.AuthApi.OauthTokenRefreshToken(context.Background())
		resp, _, err := req.
			RefreshToken(st.RefreshToken).
			ClientId(app.Conf.BaiduPan.AppKey).
			ClientSecret(app.Conf.BaiduPan.SecretKey).
			Execute()
		if err != nil {
			return err
		}
		if resp.AccessToken == nil || *resp.AccessToken == "" {
			return fmt.Errorf("refresh token failed: empty access_token")
		}

		TokenResp.AccessToken = resp.AccessToken
		TokenResp.RefreshToken = resp.RefreshToken
		TokenResp.ExpiresIn = resp.ExpiresIn
		TokenDeadline = time.Now().Add(time.Second * time.Duration(*resp.ExpiresIn))

		_ = util.SaveStoredToken(util.StoredToken{
			AccessToken:  *resp.AccessToken,
			RefreshToken: derefStr(resp.RefreshToken),
			ExpiresAt:    TokenDeadline.Unix(),
		})

		// 后台保持刷新
		runRefreshToken()
		return nil
	})
}
