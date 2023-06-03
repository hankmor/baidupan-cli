package main

import (
	"baidupan-cli/app"
	_ "baidupan-cli/cmd"
	"github.com/desertbit/grumble"
)

func main() {
	grumble.Main(app.App)
}
