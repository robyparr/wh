package cmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"
)

func TestRunStartCmd(t *testing.T) {
	type testCase struct {
		title      string
		timeStr    string
		wantPeriod model.WorkPeriod
	}

	today := util.TodayAtMidnight()
	testCases := []testCase{
		{
			title: "start default",
			wantPeriod: model.WorkPeriod{
				Id:        1,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: false},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			title:   "start with time arg",
			timeStr: "09:00",
			wantPeriod: model.WorkPeriod{
				Id:        1,
				WorkDayId: 1,
				StartAt:   today.Add(9 * time.Hour),
				EndAt:     sql.NullTime{Valid: false},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			title:   "start with relative time arg",
			timeStr: "-1h30m",
			wantPeriod: model.WorkPeriod{
				Id:        1,
				WorkDayId: 1,
				StartAt:   time.Now().Add(-90 * time.Minute),
				EndAt:     sql.NullTime{Valid: false},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			out := &bytes.Buffer{}
			repo := testutil.NewRepo(t)
			workDay, err := repo.CreateWorkDay(model.NewWorkDay(today))
			testutil.AssertNoErr(t, err)

			err = runStartCmd(out, repo, tc.timeStr)
			testutil.AssertNoErr(t, err)
			testutil.AssertOutput(t, out, fmt.Sprintf("Started tracking time on work day #1 (%s).\n", util.FormatDate(today)))

			gotPeriods, err := repo.GetWorkPeriods(workDay)
			testutil.AssertNoErr(t, err)

			if len(gotPeriods) != 1 {
				t.Fatalf("Expected work day to have 1 period, has %d", len(gotPeriods))
			}

			testutil.AssertEqualStructs(t, gotPeriods[0], tc.wantPeriod)
		})
	}
}

func TestRunStartCmdNoWorkDay(t *testing.T) {
	repo := testutil.NewRepo(t)
	out := &bytes.Buffer{}

	err := runStartCmd(out, repo, "")
	testutil.AssertNoErr(t, err)

	midnight := util.TodayAtMidnight()
	testutil.AssertOutput(t, out, fmt.Sprintf("Started tracking time on NEW work day #1 (%s).\n", util.FormatDate(midnight)))

	gotWorkDay, err := repo.GetWorkDayByDate(midnight)
	testutil.AssertNoErr(t, err)
	testutil.AssertEqualStructs(t, gotWorkDay, model.WorkDay{
		Id:        1,
		Date:      midnight,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	gotWorkPeriods, err := repo.GetWorkPeriods(gotWorkDay)
	testutil.AssertNoErr(t, err)

	if len(gotWorkPeriods) != 1 {
		t.Fatalf("Expected 1 work period, got %d", len(gotWorkPeriods))
	}

	testutil.AssertEqualStructs(t, gotWorkPeriods[0], model.WorkPeriod{
		Id:        1,
		WorkDayId: gotWorkDay.Id,
		StartAt:   time.Now(),
		EndAt:     sql.NullTime{Valid: false},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}
