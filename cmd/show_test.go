package cmd

import (
	"bytes"
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

		compareShowOutput(
			t,
			out.String(),
			`
September 01, 2023 (Fri)
========================

Work Day: 		7h30m
Time Worked:		0m
Time Remaining:		7h30m
Estimated Finish:	:estFinish
Note:			This is a note.


WORK PERIODS
ID	START			END			TIME WORKED	NOTE

`,
			map[string]string{
				"estFinish": util.FormatDateTime(time.Now().Add(7 * time.Hour).Add(30 * time.Minute)),
			},
		)
	})

	t.Run("with work periods", func(t *testing.T) {
		wp1 := model.NewWorkPeriod(wd)
		wp1.StartAt = time.Date(2023, 9, 1, 9, 0, 0, 0, time.Local)
		wp1.SetEndAt(wp1.StartAt.Add(1 * time.Hour))
		wp1.SetNote("Period note.")

		_, err = repo.CreateWorkPeriod(wp1)
		testutil.AssertNoErr(t, err)

		wp2 := model.NewWorkPeriod(wd)
		wp2.StartAt = time.Date(2023, 9, 1, 10, 0, 0, 0, time.Local)
		wp2.SetEndAt(wp2.StartAt.Add(30 * time.Minute))
		_, err = repo.CreateWorkPeriod(wp2)
		testutil.AssertNoErr(t, err)

		out := &bytes.Buffer{}
		err = runShowCmd(out, repo, "2023-09-01")
		testutil.AssertNoErr(t, err)

		compareShowOutput(
			t,
			out.String(),
			`
September 01, 2023 (Fri)
========================

Work Day: 		7h30m
Time Worked:		1h30m
Time Remaining:		6h0m
Estimated Finish:	:estFinish
Note:			This is a note.


WORK PERIODS
ID	START			END			TIME WORKED	NOTE
1	2023-09-01 9:00 AM	2023-09-01 10:00 AM 	1h0m		Period note.
2	2023-09-01 10:00 AM	2023-09-01 10:30 AM 	30m	:tab

`,
			map[string]string{
				"estFinish":    util.FormatDateTime(time.Now().Add(6 * time.Hour)),
				"p2TimeWorked": util.FormatDuration(wp2.TimeWorked()),
				"tab":          "	",
			},
		)
	})
}

func compareShowOutput(t *testing.T, got string, want string, replacements map[string]string) {
	want = strings.TrimPrefix(want, "\n")
	for k, v := range replacements {
		want = strings.ReplaceAll(want, ":"+k, v)
	}

	if got != want {
		t.Errorf("got `%s`, want `%s`", got, want)
	}
}
