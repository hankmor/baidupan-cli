package util

import (
	"github.com/mdp/qrterminal/v3"
	"os"
)

// PrintQrCode2Console 在控制台输出url对应的二维码
func PrintQrCode2Console(url string) {
	config := qrterminal.Config{
		Level:      qrterminal.M,
		Writer:     os.Stdout,
		BlackChar:  qrterminal.WHITE,
		WhiteChar:  qrterminal.BLACK,
		QuietZone:  1,
		HalfBlocks: false,
	}
	qrterminal.GenerateWithConfig(url, config)
}
