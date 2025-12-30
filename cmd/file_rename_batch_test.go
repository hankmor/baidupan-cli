package cmd

import (
	"regexp"
	"testing"
)

func TestBuildBatchRenamePlan_Basic(t *testing.T) {
	re := regexp.MustCompile(`^(.*)\.mp4$`)
	files := []*File{
		{Path: "/a/1.mp4", ServerFilename: "1.mp4", IsDir: 0},
		{Path: "/a/2.txt", ServerFilename: "2.txt", IsDir: 0},
	}
	items, plan, err := buildBatchRenamePlan(files, re, `$1_new.mp4`, renameBatchTargetFiles)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(items) != 1 || len(plan) != 1 {
		t.Fatalf("unexpected len: items=%d plan=%d", len(items), len(plan))
	}
	if items[0].Path != "/a/1.mp4" || items[0].Newname != "1_new.mp4" {
		t.Fatalf("unexpected item: %+v", items[0])
	}
	if plan[0].NewPath != "/a/1_new.mp4" {
		t.Fatalf("unexpected new path: %s", plan[0].NewPath)
	}
}

func TestBuildBatchRenamePlan_Conflict(t *testing.T) {
	re := regexp.MustCompile(`^a.*\.txt$`)
	files := []*File{
		{Path: "/a/a1.txt", ServerFilename: "a1.txt", IsDir: 0},
		{Path: "/a/a2.txt", ServerFilename: "a2.txt", IsDir: 0},
	}
	_, _, err := buildBatchRenamePlan(files, re, `same.txt`, renameBatchTargetFiles)
	if err == nil {
		t.Fatalf("expected conflict error")
	}
}

func TestBuildBatchRenamePlan_Target(t *testing.T) {
	re := regexp.MustCompile(`^old$`)
	files := []*File{
		{Path: "/a/old", ServerFilename: "old", IsDir: 1},
	}
	_, plan0, err := buildBatchRenamePlan(files, re, `new`, renameBatchTargetFiles)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(plan0) != 0 {
		t.Fatalf("expected dirs excluded when target=files")
	}
	_, plan1, err := buildBatchRenamePlan(files, re, `new`, renameBatchTargetDirs)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(plan1) != 1 {
		t.Fatalf("expected dirs included when target=dirs")
	}
}

func TestBuildBatchRenameMatcher_SimpleLiteral(t *testing.T) {
	re, repl, err := buildBatchRenameMatcher("", "", "设计", "分析", false)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repl != "分析" {
		t.Fatalf("unexpected repl: %q", repl)
	}
	if got := re.ReplaceAllString("UML设计图", repl); got != "UML分析图" {
		t.Fatalf("unexpected replace result: %q", got)
	}
}

func TestBuildBatchRenameMatcher_FullRegexRequiredPair(t *testing.T) {
	_, _, err := buildBatchRenameMatcher("a", "", "", "", false)
	if err == nil {
		t.Fatalf("expected err")
	}
	_, _, err = buildBatchRenameMatcher("", "b", "", "", false)
	if err == nil {
		t.Fatalf("expected err")
	}
}
