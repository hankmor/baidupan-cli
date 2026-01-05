// Package app is the entrence of this app.
package app

import (
	openapi "baidupan-cli/openxpanapi"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

const version = "0.1.0"

var (
	Conf *Config
	App  = grumble.New(&grumble.Config{
		Name:                  "baidupan-cli",
		Description:           "baidu network disk command line tool",
		Prompt:                "cli » ",
		PromptColor:           color.New(color.FgGreen, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,

		Flags: func(f *grumble.Flags) {
			f.String("c", "config", "./config.yaml", "set system config file")
			// TODO remove
			f.Bool("t", "test", true, "use test mode")
			// debug helper: allow injecting access token directly (avoids relying on expired mock token)
			f.StringL("access-token", "", "override baidupan access token (debug)")
		},
	})
	
	CurrentDir = "/"
)
var APIClient *openapi.APIClient

var initHooks []func(a *grumble.App, flags grumble.FlagMap) error

// RegisterInitHook registers an initialization hook.
// Hooks are executed in registration order before the first command.
func RegisterInitHook(h func(a *grumble.App, flags grumble.FlagMap) error) {
	initHooks = append(initHooks, h)
}

func init() {
	// 默认 hook：加载配置文件（必须最先执行）
	RegisterInitHook(func(a *grumble.App, flags grumble.FlagMap) error {
		if c, err := LoadConf(flags.String("config")); err != nil {
			return err
		} else {
			Conf = c
			return nil
		}
	})

	// 程序执行第一个命令前执行所有 hooks
	App.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		for _, h := range initHooks {
			if err := h(a, flags); err != nil {
				return err
			}
		}
		return nil
	})

	// 创建 api client
	APIClient = openapi.NewAPIClient(openapi.NewConfiguration())

	// 打印
	App.SetPrintASCIILogo(func(a *grumble.App) {
		logo := `
___.          .__    .___                                         .__  .__ 
\_ |__ _____  |__| __| _/_ _____________    ____             ____ |  | |__|
 | __ \\__  \ |  |/ __ |  |  \____ \__  \  /    \   ______ _/ ___\|  | |  |
 | \_\ \/ __ \|  / /_/ |  |  /  |_> > __ \|   |  \ /_____/ \  \___|  |_|  |
 |___  (____  /__\____ |____/|   __(____  /___|  /          \___  >____/__|
     \/     \/        \/     |__|       \/     \/               \/
`
		_, _ = a.Println(logo)
		_, _ = a.Printf("  Version %s\n", version)
		_, _ = a.Println()
	})
}

func RegisterCommand(cmd *grumble.Command) {
	if cmd != nil {
		App.AddCommand(cmd)
	}
}
