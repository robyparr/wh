package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"
)

func TestRunShowCmdNoWorkday(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)

	err := runShowCmd(out, repo, "2023-09-01")
	testutil.AssertNoErr(t, err)

	got := out.String()
	want := "No work day for 2023-09-01 yet.\n"
	if got != want {
		t.Errorf("got `%s`, want `%s`", got, want)
	}
}

func TestRunShowCmd(t *testing.T) {
	repo := testutil.NewRepo(t)

	date := time.Date(2023, 9, 1, 0, 0, 0, 0, time.Local)
	wd := model.NewWorkDay(date)
	wd.SetNote("This is a note.")

	wd, err := repo.CreateWorkDay(wd)
	testutil.AssertNoErr(t, err)

	t.Run("without work periods", func(t *testing.T) {
		out := &bytes.Buffer{}
		err = runShowCmd(out, repo, "2023-09-01")
		testutil.AssertNoErr(t, err)

		wantEstFinish := time.Now().Add(7 * time.Hour).Add(30 * time.Minute)
		compareShowOutput(
			t,
			out.String(),
			`
September 01, 2023 (Fri)
========================

Work Day: 				7h30m
Time Worked:			0m
Time Remaining:		7h30m
Estimated Finish: %s
Note:							This is a note.
`,
			wantEstFinish,
		)
	})

	t.Run("with work periods", func(t *testing.T) {
		wp := model.NewWorkPeriod(wd)
		wp.SetEndAt(wp.StartAt.Add(1 * time.Hour))

		_, err = repo.CreateWorkPeriod(wp)
		testutil.AssertNoErr(t, err)

		wp = model.NewWorkPeriod(wd)
		wp.StartAt = time.Now().Add(2 * time.Hour)
		wp.SetEndAt(wp.StartAt.Add(30 * time.Minute))
		_, err = repo.CreateWorkPeriod(wp)
		testutil.AssertNoErr(t, err)

		out := &bytes.Buffer{}
		err = runShowCmd(out, repo, "2023-09-01")
		testutil.AssertNoErr(t, err)

		wantEstFinish := time.Now().Add(6 * time.Hour)
		compareShowOutput(
			t,
			out.String(),
			`
September 01, 2023 (Fri)
========================

Work Day: 				7h30m
Time Worked:			1h30m
Time Remaining:		6h0m
Estimated Finish: %s
Note:							This is a note.
`,
			wantEstFinish,
		)
	})
}

func compareShowOutput(t *testing.T, got string, want string, wantEstFinish time.Time) {
	want = fmt.Sprintf(want, util.FormatDateTime(wantEstFinish))
	want = strings.TrimPrefix(want, "\n")

	if got != want {
		t.Errorf("got `%s`, want `%s`", got, want)
	}
}
