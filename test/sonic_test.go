package test

import (
	"baidupan-cli/cmd"
	"fmt"
	"github.com/bytedance/sonic"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	s := "{\"cursor\":7,\"errmsg\":\"succ\",\"errno\":0,\"has_more\":0,\"list\":[{\"category\":6,\"fs_id\":94161225432639,\"isdir\":1,\"local_ctime\":1419989944,\"local_mtime\":1419989944,\"md5\":\"\",\"path\":\"/00-电子书/工具软件/visio\",\"server_ctime\":1419989944,\"server_filename\":\"visio\",\"server_mtime\":1640073745,\"size\":0},{\"category\":6,\"fs_id\":796579127537769,\"isdir\":1,\"local_ctime\":1397729790,\"local_mtime\":1397729790,\"md5\":\"\",\"path\":\"/00-电子书/工具软件/uml\",\"server_ctime\":1397729790,\"server_filename\":\"uml\",\"server_mtime\":1640073745,\"size\":0},{\"category\":4,\"fs_id\":335604055992110,\"isdir\":0,\"local_ctime\":1397729793,\"local_mtime\":1397729793,\"md5\":\"3bfb0fd3334a1d76b9c20576626bb682\",\"path\":\"/00-电子书/工具软件/uml/[大象Thinking.in.UML].ThinkingInUML.pdf\",\"server_ctime\":1397729794,\"server_filename\":\"[大象Thinking.in.UML].ThinkingInUML.pdf\",\"server_mtime\":1640073745,\"size\":61939320},{\"category\":4,\"fs_id\":855119404237552,\"isdir\":0,\"local_ctime\":1397729789,\"local_mtime\":1397729789,\"md5\":\"7a9904d49fdb9ab8cbbd18af0bf320c7\",\"path\":\"/00-电子书/工具软件/uml/uml精粹(第三版)标准对象建模语言简明指南.pdf\",\"server_ctime\":1397729790,\"server_filename\":\"uml精粹(第三版)标准对象建模语言简明指南.pdf\",\"server_mtime\":1640073745,\"size\":81144580},{\"category\":4,\"fs_id\":999612811922406,\"isdir\":0,\"local_ctime\":1420532396,\"local_mtime\":1420532396,\"md5\":\"631ac4d30cc98c1d1b98bae7020fd9c2\",\"path\":\"/00-电子书/工具软件/uml/[大家网]UML用户指南(第2版)[www.TopSage.com].pdf\",\"server_ctime\":1420532404,\"server_filename\":\"[大家网]UML用户指南(第2版)[www.TopSage.com].pdf\",\"server_mtime\":1640073745,\"size\":144056256},{\"category\":4,\"fs_id\":1005862119845709,\"isdir\":0,\"local_ctime\":1399775765,\"local_mtime\":1399775765,\"md5\":\"a2484e4945742ed2041b9613753bfcc2\",\"path\":\"/00-电子书/工具软件/uml/[www.java1234.com]UML系统建模与分析设计.刁成嘉.pdf\",\"server_ctime\":1399775766,\"server_filename\":\"[www.java1234.com]UML系统建模与分析设计.刁成嘉.pdf\",\"server_mtime\":1640073745,\"size\":148521405},{\"category\":4,\"fs_id\":894633045476667,\"isdir\":0,\"local_ctime\":1419989963,\"local_mtime\":1419989963,\"md5\":\"584f7398fc353387a570d54c94ac00aa\",\"path\":\"/00-电子书/工具软件/visio/visio 2010图形设计从新手到高手.pdf\",\"server_ctime\":1419989965,\"server_filename\":\"visio 2010图形设计从新手到高手.pdf\",\"server_mtime\":1640073745,\"size\":43006680}],\"request_id\":\"222714725402945750\"}"
	var recursionFileResp cmd.RecursionFileResp
	err := sonic.UnmarshalString(s, &recursionFileResp)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Println(recursionFileResp)
	}
}
