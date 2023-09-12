package util_test

import (
	"testing"
	"time"

	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"
)

func TestTodayAtMidnight(t *testing.T) {
	got := util.TodayAtMidnight()

	now := time.Now()
	want := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if got != want {
		t.Errorf("got '%v', want '%v'", got, want)
	}
}

func TestFormatDate(t *testing.T) {
	date := time.Date(2023, 8, 13, 9, 30, 0, 0, time.Local)
	got := util.FormatDate(date)
	want := "2023-08-13"

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestParseTimeString(t *testing.T) {
	midnight := util.TodayAtMidnight()

	testCases := []struct {
		name  string
		input string
		want  time.Time
	}{
		{name: "empty string", input: "", want: time.Now()},
		{name: "exact time", input: "09:30", want: midnight.Add(9 * time.Hour).Add(30 * time.Minute)},
		{name: "exact time afternoon", input: "13:00", want: midnight.Add(13 * time.Hour)},
		{name: "relative time", input: "1h30m", want: time.Now().Add(90 * time.Minute)},
		{name: "relative time mins", input: "30m", want: time.Now().Add(30 * time.Minute)},
		{name: "relative time past", input: "-30m", want: time.Now().Add(-30 * time.Minute)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := util.ParseTimeString(tc.input)
			testutil.AssertNoErr(t, err)
			testutil.AssertAroundTime(t, "result", got, tc.want)
		})
	}
}
