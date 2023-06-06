package util

import (
	"strconv"
	"strings"
)

func Encode2Unicode(s string) string {
	return strconv.QuoteToASCII(s)
}

func Decode2Chinese(unicodeStr string) string {
	textQuoted := strconv.Quote(unicodeStr)
	str, err := strconv.Unquote(strings.Replace(textQuoted, `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}
