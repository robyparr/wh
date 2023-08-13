package util

import "time"

func TodayAtMidnight() time.Time {
	today := time.Now()
	return time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Local().Location())
}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
