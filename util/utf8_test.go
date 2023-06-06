package util

import (
	"fmt"
	"testing"
)

func TestDecode2Chinese(t *testing.T) {
	s := "这是一个测试的中文字符串，this is a chinese string for testing"
	ret := Encode2Unicode(s)
	fmt.Println(ret)
	ret, _ = Decode2Chinese(ret)
	fmt.Println(ret)
}

func TestDecode(t *testing.T) {
	s := "\\/00-\\u7535\\u5b50\\u4e66\\/\\u9ad8\\u65b0\\u6280\\u672f\\/[\\u6e38\\u620f\\u5f00\\u53d1\\u4e2d\\u7684\\u4eba\\u5de5\\u667a\\u80fd].pdf"
	ret, _ := Decode2Chinese(s)
	fmt.Println(ret)
}
