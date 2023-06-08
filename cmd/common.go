package cmd

import (
	"baidupan-cli/util"
	"fmt"
	"github.com/desertbit/grumble"
)

type BaseVo struct {
	Errno int `json:"errno,omitempty"`
}

func (v *BaseVo) Success() bool {
	return v.Errno == 0
}

func checkAuthorized(ctx *grumble.Context) error {
	test := ctx.Flags.Bool("test")
	if test {
		TokenResp = util.MockAccessToken()
	}
	if TokenResp == nil {
		return fmt.Errorf("not authorized, execute `auth` command to authorize first")
	}
	return nil
}
