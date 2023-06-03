package app

import (
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

const version = "0.1.0"

var Conf *Config
var App = grumble.New(&grumble.Config{
	Name:                  "baidupan-cli",
	Description:           "百度网盘命令行工具",
	Prompt:                "cli » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.String("c", "config", "./config.yaml", "设置系统配置文件")
	},
})

func init() {
	// 程序执行第一个命令前执行此函数，用于加载配置文件
	App.OnInit(func(a *grumble.App, flags grumble.FlagMap) error {
		if c, err := LoadConf(flags.String("config")); err != nil {
			return err
		} else {
			Conf = c
			return nil
		}
	})

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
