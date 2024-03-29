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

func TestRunAddCmd(t *testing.T) {
	type testCase struct {
		title           string
		dateStr         string
		lengthStr       string
		note            string
		expectedOutput  string
		expectedWorkDay model.WorkDay
	}

	date := time.Date(2023, 8, 13, 0, 0, 0, 0, time.Local)
	var defaultLength int = 7.5 * 60

	testCases := []testCase{
		{
			title:          "add default",
			expectedOutput: fmt.Sprintf("Added work day #1 on %v\n", util.FormatDate(time.Now())),
			expectedWorkDay: model.WorkDay{
				Id:         1,
				Date:       util.TodayAtMidnight(),
				LengthMins: defaultLength,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			title:          "add with date arg",
			dateStr:        "2023-08-13",
			expectedOutput: "Added work day #1 on 2023-08-13\n",
			expectedWorkDay: model.WorkDay{
				Id:         1,
				Date:       date,
				LengthMins: defaultLength,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			title:          "add with length arg",
			dateStr:        "2023-08-13",
			lengthStr:      "1h30m",
			expectedOutput: "Added work day #1 on 2023-08-13\n",
			expectedWorkDay: model.WorkDay{
				Id:         1,
				Date:       date,
				LengthMins: 90,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			title:          "add with note arg",
			dateStr:        "2023-08-13",
			note:           "This is a note.",
			expectedOutput: "Added work day #1 on 2023-08-13\n",
			expectedWorkDay: model.WorkDay{
				Id:         1,
				Date:       date,
				Note:       sql.NullString{String: "This is a note.", Valid: true},
				LengthMins: defaultLength,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			out := &bytes.Buffer{}
			repo := testutil.NewRepo(t)

			err := runAddCmd(out, repo, tc.dateStr, tc.lengthStr, tc.note)
			testutil.AssertNoErr(t, err)
			testutil.AssertOutput(t, out, tc.expectedOutput)

			gotRecord, err := repo.GetWorkDayByDate(tc.expectedWorkDay.Date)
			testutil.AssertNoErr(t, err)
			testutil.AssertWorkDay(t, gotRecord, tc.expectedWorkDay)
		})
	}
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

	err := runAddCmd(out, repo, "", "", "")
	testutil.AssertNoErr(t, err)
	testutil.AssertOutput(t, out, fmt.Sprintf("Work day on %s already exists.\n", util.FormatDate(today)))

	gotCount, err := repo.GetWorkDayCount()
	testutil.AssertNoErr(t, err)

	wantCount := 1
	if gotCount != wantCount {
		t.Errorf("Unexpected work day count; got %d, want %d\n", gotCount, wantCount)
	}
}
