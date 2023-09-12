package cmd

import (
	"bytes"
	"database/sql"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"
)

func TestRunStopCmd(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)

	midnight := util.TodayAtMidnight()
	workDay, err := repo.CreateWorkDay(model.NewWorkDay(midnight))
	testutil.AssertNoErr(t, err)

	testCases := []struct {
		name       string
		timeStr    string
		note       string
		wantPeriod model.WorkPeriod
	}{
		{
			name: "default",
			wantPeriod: model.WorkPeriod{
				Id:        1,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: true, Time: time.Now()},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name:    "with exact time arg",
			timeStr: "17:00",
			wantPeriod: model.WorkPeriod{
				Id:        2,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: true, Time: midnight.Add(17 * time.Hour)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name:    "with relative time arg",
			timeStr: "1h30m",
			wantPeriod: model.WorkPeriod{
				Id:        3,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: true, Time: time.Now().Add(90 * time.Minute)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name:    "with past relative time arg",
			timeStr: "-1h30m",
			wantPeriod: model.WorkPeriod{
				Id:        4,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: true, Time: time.Now().Add(-90 * time.Minute)},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "with note flag",
			note: "This is a note.",
			wantPeriod: model.WorkPeriod{
				Id:        5,
				WorkDayId: 1,
				StartAt:   time.Now(),
				EndAt:     sql.NullTime{Valid: true, Time: time.Now()},
				Note:      sql.NullString{Valid: true, String: "This is a note."},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err = repo.CreateWorkPeriod(model.NewWorkPeriod(workDay))
			testutil.AssertNoErr(t, err)

			err = runStopCmd(out, repo, tc.timeStr, tc.note)
			testutil.AssertNoErr(t, err)

			got := out.String()
			want := ""
			if got != want {
				t.Errorf("got '%s', want '%s'", got, want)
			}

			periods, err := repo.GetWorkPeriods(workDay)
			testutil.AssertNoErr(t, err)

			gotPeriod := periods[len(periods)-1]
			testutil.AssertEqualStructs(t, gotPeriod, tc.wantPeriod)
		})
	}
}

func TestRunStopCmdNoWorkDay(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)

	err := runStopCmd(out, repo, "", "")
	testutil.AssertNoErr(t, err)

	got := out.String()
	want := "Unable to find an ongoing work period.\n"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestRunStopCmdNoOpenPeriod(t *testing.T) {
	out := &bytes.Buffer{}
	repo := testutil.NewRepo(t)

	workDay, err := repo.CreateWorkDay(model.NewWorkDay(util.TodayAtMidnight()))
	testutil.AssertNoErr(t, err)

	period := model.NewWorkPeriod(workDay)
	period.EndAt = sql.NullTime{Valid: true, Time: time.Now()}

	_, err = repo.CreateWorkPeriod(period)
	testutil.AssertNoErr(t, err)

	err = runStopCmd(out, repo, "", "")
	testutil.AssertNoErr(t, err)

	got := out.String()
	want := "Unable to find an ongoing work period.\n"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}

	periods, err := repo.GetWorkPeriods(workDay)
	testutil.AssertNoErr(t, err)

	gotPeriodCount := len(periods)
	if gotPeriodCount != 1 {
		t.Errorf("got %d work periods, want 1", gotPeriodCount)
	}
}
