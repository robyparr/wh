package model_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
)

func TestNewWorkDay(t *testing.T) {
	date := time.Date(2023, 8, 1, 0, 0, 0, 0, time.Local)
	got := model.NewWorkDay(date)
	want := model.WorkDay{
		Date:       date,
		LengthMins: 7.5 * 60,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

func TestSetNote(t *testing.T) {
	wd := model.NewWorkDay(util.TodayAtMidnight())
	wd.SetNote("Hello!")

	got := wd.Note
	want := sql.NullString{String: "Hello!", Valid: true}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
