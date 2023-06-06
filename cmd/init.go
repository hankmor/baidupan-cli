package cmd

import "baidupan-cli/app"

// =============================================
// 初始化
// =============================================

func init() {
	app.RegisterCommand(authCmd)
	app.RegisterCommand(capCmd)
	app.RegisterCommand(userInfoCmd)
	app.RegisterCommand(fileListCmd)
}
