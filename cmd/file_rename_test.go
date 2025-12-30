package cmd

import (
	"strings"
	"testing"
)

func TestBuildRenameFilelist(t *testing.T) {
	filelist, dst, err := buildRenameFilelist("/a/b.txt", "c.txt")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if dst != "/a/c.txt" {
		t.Fatalf("unexpected dst: %s", dst)
	}
	if filelist == "" {
		t.Fatalf("expected non-empty filelist")
	}
	// 只做关键字段断言，避免对 JSON key 顺序过于敏感
	if !(containsAll(filelist, []string{`"path"`, `"/a/b.txt"`, `"newname"`, `"c.txt"`})) {
		t.Fatalf("unexpected filelist: %s", filelist)
	}
}

func TestBuildRenameFilelist_Invalid(t *testing.T) {
	cases := []struct {
		path    string
		newname string
	}{
		{"", "a"},
		{"/", "a"},
		{"a.txt", "b.txt"},
		{"/a", ""},
		{"/a", "b/c"},
	}
	for _, c := range cases {
		_, _, err := buildRenameFilelist(c.path, c.newname)
		if err == nil {
			t.Fatalf("expected err for path=%q newname=%q", c.path, c.newname)
		}
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
