package cmd

import (
	openapi "baidupan-cli/openxpanapi"
	"baidupan-cli/util"
	"fmt"
	"strings"

	"github.com/desertbit/grumble"
)

type BaseVo struct {
	Errno int `json:"errno,omitempty"`
}

func (v *BaseVo) Success() bool {
	return v.Errno == 0
}

func checkAuthorized(ctx *grumble.Context) error {
	// 优先使用显式传入的 access token，便于 debug（不会被 test mode 的 mock 覆盖）
	if TokenResp == nil {
		if t := strings.TrimSpace(ctx.Flags.String("access-token")); t != "" {
			TokenResp = &openapi.OauthTokenDeviceTokenResponse{AccessToken: &t}
		}
	}

	test := ctx.Flags.Bool("test")
	if test && TokenResp == nil {
		TokenResp = util.MockAccessToken()
	}

	if TokenResp == nil || TokenResp.AccessToken == nil || strings.TrimSpace(*TokenResp.AccessToken) == "" {
		return fmt.Errorf("not authorized, execute `auth` command to authorize first")
	}
	return nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
