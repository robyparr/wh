package model_test

import (
	"reflect"
	"testing"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
)

func TestNewWorkDay(t *testing.T) {
	got := model.NewWorkDay()
	want := model.WorkDay{
		Date:       util.TodayAtMidnight(),
		LengthMins: 7.5 * 60,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v\n", got, want)
	}
}
