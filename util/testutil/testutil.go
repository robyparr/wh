package testutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/repository"
)

func NewRepo(t *testing.T) *repository.Repo {
	repo, err := repository.NewRepo(":memory:")
	AssertNoErr(t, err)

	return repo
}

func AssertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
}

func AssertWorkDay(t *testing.T, got model.WorkDay, want model.WorkDay) {
	t.Helper()

	assertAroundTime(t, got.Date, want.Date)
	assertAroundTime(t, got.CreatedAt, want.CreatedAt)
	assertAroundTime(t, got.UpdatedAt, want.UpdatedAt)

	got.Date = time.Time{}
	got.CreatedAt = time.Time{}
	got.UpdatedAt = time.Time{}

	want.Date = time.Time{}
	want.CreatedAt = time.Time{}
	want.UpdatedAt = time.Time{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %+v\nwant %+v", got, want)
	}
}

func assertAroundTime(t *testing.T, got time.Time, want time.Time) {
	t.Helper()

	min := want.Add(-1 * time.Second)
	max := want.Add(1 * time.Second)

	if got.Before(min) || got.After(max) {
		t.Errorf("Expected %v to be around %v.\n", got, want)
	}
}
