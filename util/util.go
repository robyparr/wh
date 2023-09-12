package util

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DateFormatStr string = "2006-01-02"

var exactTimeRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)
var relativeTimeRegex = regexp.MustCompile(`^(-?\d+h(\d+m)?)|(-?\d+m)$`)

func TodayAtMidnight() time.Time {
	today := time.Now()
	return time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Local().Location())
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormatStr)
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
