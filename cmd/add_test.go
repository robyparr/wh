package cmd

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"
)

func TestRunAddCmd(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)

	err := runAddCmd(out, repo)
	testutil.AssertNoErr(t, err)
	assertOutput(t, out, fmt.Sprintf("Added work day #1 on %v\n", util.FormatDate(time.Now())))

	gotRecord, err := repo.GetWorkDayByDate(util.TodayAtMidnight())
	testutil.AssertNoErr(t, err)
	testutil.AssertWorkDay(t, gotRecord, model.WorkDay{
		Id:         1,
		Date:       util.TodayAtMidnight(),
		LengthMins: 7.5 * 60,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
}

func TestRunAddCmdExistingDay(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)
	today := util.TodayAtMidnight()
	repo.CreateWorkDay(model.WorkDay{
		Id:         1,
		Date:       today,
		LengthMins: 7.5 * 60,
	})

	err := runAddCmd(out, repo)
	testutil.AssertNoErr(t, err)
	assertOutput(t, out, fmt.Sprintf("Work day on %s already exists.\n", util.FormatDate(today)))

	gotCount, err := repo.GetWorkDayCount()
	testutil.AssertNoErr(t, err)

	wantCount := 1
	if gotCount != wantCount {
		t.Errorf("Unexpected work day count; got %d, want %d\n", gotCount, wantCount)
	}
}

func assertOutput(t *testing.T, got *bytes.Buffer, want string) {
	t.Helper()

	out := got.String()
	if out != want {
		t.Errorf("Unexpected output:\ngot: %s\nwant: %s\n", got, want)
	}
}
