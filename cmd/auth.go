package cmd

import (
	"baidupan-cli/app"
	"fmt"
	"github.com/desertbit/grumble"
)

var loginCmd = &grumble.Command{
	Name:     "auth",
	Help:     "authorize cli to visit your baidupan account",
	LongHelp: "scan the given qrcode to authorize cli to visit your baidupan account",
	Run: func(c *grumble.Context) error {
		fmt.Println(c.Flags.String("config"))
		fmt.Println(app.Conf)
		return nil
	},
}

func init() {
	app.RegisterCommand(loginCmd)
}
