package util

import "time"

const DateFormatStr string = "2006-01-02"

func TodayAtMidnight() time.Time {
	today := time.Now()
	return time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Local().Location())
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormatStr)
}
