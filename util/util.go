package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DateFormatStr string = "2006-01-02"

var exactTimeRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)
var relativeTimeRegex = regexp.MustCompile(`^(-?\d+h(\d+m)?)|(-?\d+m)$`)

func TodayAtMidnight() time.Time {
	return timeAtMidnight(time.Now())
}

func timeAtMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormatStr)
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 3:04 PM")
}

func ParseTimeString(str string) (time.Time, error) {
	startAt := time.Now()

	switch {
	case exactTimeRegex.MatchString(str):
		duration, err := parseExactTimeString(str)
		if err != nil {
			return time.Time{}, err
		}

		startAt = TodayAtMidnight().Add(duration)
	case relativeTimeRegex.MatchString(str):
		duration, err := time.ParseDuration(str)
		if err != nil {
			return time.Time{}, err
		}

		startAt = startAt.Add(duration)
	}

	return startAt, nil
}

func parseExactTimeString(str string) (time.Duration, error) {
	timeStrParts := strings.Split(str, ":")
	hour, err := strconv.Atoi(timeStrParts[0])
	if err != nil {
		return 0, nil
	}
	min, err := strconv.Atoi(timeStrParts[1])
	if err != nil {
		return 0, err
	}

	totalMinutes := (hour * 60) + min
	return time.Duration(totalMinutes) * time.Minute, nil
}

func ParseDateString(str string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, err
	}

	return timeAtMidnight(date), nil
}

func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0m"
	}

	d = d.Truncate(time.Minute)
	return strings.Replace(d.String(), "0s", "", 1)
}

func Underline(str string) string {
	underline := strings.Repeat("=", len(str))
	return fmt.Sprintf("%s\n%s", str, underline)
}
