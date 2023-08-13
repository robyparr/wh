package util_test

import (
	"testing"
	"time"

	"github.com/robyparr/wh/util"
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
