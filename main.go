package main

import (
	"baidupan-cli/app"
	_ "baidupan-cli/cmd"
	"os"
	"path/filepath"

	"github.com/desertbit/grumble"
)

func main() {
	// 针对 go run 运行场景：
	// 如果未指定 BAIDUPAN_CLI_TOKEN_DIR，则默认使用当前目录下的 .debug 目录
	// 避免 go run 产生的临时二进制文件找不到 token
	if os.Getenv("BAIDUPAN_CLI_TOKEN_DIR") == "" {
		wd, _ := os.Getwd()
		os.Setenv("BAIDUPAN_CLI_TOKEN_DIR", filepath.Join(wd, ".debug"))
	}
	grumble.Main(app.App)
}
