package util

import (
	"fmt"
	"github.com/hankmor/gotools/date"
	"time"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

func ConvReadableSize(size int64) string {
	if size <= 0 {
		return ""
	}
	tb := size / TB
	if tb > 0 {
		if tb < 10 {
			return fmt.Sprintf("%dGB", size/GB)
		} else {
			return fmt.Sprintf("%dTB", tb)
		}
	}
	gb := size / GB
	if gb > 0 {
		if gb <= 5 {
			return fmt.Sprintf("%dMB", size/MB)
		} else {
			return fmt.Sprintf("%dGB", gb)
		}
	}
	mb := size / MB
	if mb > 0 {
		return fmt.Sprintf("%dMB", mb)
	}
	kb := size / KB
	if kb > 0 {
		return fmt.Sprintf("%dKB", kb)
	}
	return fmt.Sprintf("%dB", size)
}

func ConvTimestamp(ts int64) string {
	if ts <= 0 {
		return ""
	}
	tm := time.Unix(ts, 0)
	return date.FmtDateTime(tm)
}
