package testutil

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
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

func AssertOutput(t *testing.T, got *bytes.Buffer, want string) {
	t.Helper()

	out := got.String()
	if out != want {
		t.Errorf("Unexpected output:\ngot: %s\nwant: %s\n", got, want)
	}
}

func AssertWorkDay(t *testing.T, got model.WorkDay, want model.WorkDay) {
	t.Helper()

	AssertAroundTime(t, "Date", got.Date, want.Date)
	AssertAroundTime(t, "CreatedAt", got.CreatedAt, want.CreatedAt)
	AssertAroundTime(t, "UpdatedAt", got.UpdatedAt, want.UpdatedAt)

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

func AssertEqualStructs(t *testing.T, got any, want any) {
	t.Helper()

	gotStructValue := reflect.ValueOf(got)
	wantStructValue := reflect.ValueOf(want)
	structType := gotStructValue.Type()

	mismatches := []string{"Structs are not equal."}
	for i := 0; i < gotStructValue.NumField(); i++ {
		gotFieldValue := gotStructValue.Field(i)
		wantFieldValue := wantStructValue.Field(i)

		if !gotFieldValue.CanInterface() {
			continue
		}
		gotValue := gotFieldValue.Interface()
		wantValue := wantFieldValue.Interface()

		var isEqual bool
		switch gotFieldValue.Type().Name() {
		case "Time":
			isEqual = isAroundTime(gotValue.(time.Time), wantValue.(time.Time))
		case "NullTime":
			gotNullTime := gotValue.(sql.NullTime)
			wantNullTime := wantValue.(sql.NullTime)

			isEqual = gotNullTime.Valid == wantNullTime.Valid && isAroundTime(gotNullTime.Time, wantNullTime.Time)
		default:
			isEqual = gotValue == wantValue
		}

		if !isEqual {
			fieldInfo := structType.Field(i)
			mismatches = append(mismatches, fmt.Sprintf("%s: got '%v', want '%v'", fieldInfo.Name, gotValue, wantValue))
		}
	}

	if len(mismatches) > 1 {
		t.Errorf(strings.Join(mismatches, "\n\t"))
	}
}

func AssertAroundTime(t *testing.T, label string, got time.Time, want time.Time) {
	t.Helper()

	if !isAroundTime(got, want) {
		t.Errorf("%s: Expected %v to be around %v.\n", label, got, want)
	}
}

func isAroundTime(t time.Time, aroundT time.Time) bool {
	min := aroundT.Add(-1 * time.Second)
	max := aroundT.Add(1 * time.Second)

	return t.After(min) && t.Before(max)
}
