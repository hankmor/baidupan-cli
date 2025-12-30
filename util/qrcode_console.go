package util

import (
	"os"

	"github.com/mdp/qrterminal/v3"
)

// PrintQrCode2Console 在控制台输出url对应的二维码
func PrintQrCode2Console(url string) {
	config := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		// QuietZone/HalfBlocks 会显著影响可扫描性：
		// - QuietZone 太小容易扫不出来
		// - HalfBlocks=true 可以提升“像素密度”，更清晰
		QuietZone:  2,
		HalfBlocks: true,
	}
	qrterminal.GenerateWithConfig(url, config)
}
