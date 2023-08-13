package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
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
