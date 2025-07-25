package util

import "time"

func FmtDateTime(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}
