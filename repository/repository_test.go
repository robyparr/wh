package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/util"
	"github.com/robyparr/wh/util/testutil"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateWorkDay(t *testing.T) {
	repo := testutil.NewRepo(t)

	workDay := model.WorkDay{
		Date:       time.Date(2023, 8, 11, 0, 0, 0, 0, time.Local),
		LengthMins: 7.5 * 60,
		Note:       sql.NullString{String: "This is a note", Valid: true},
	}

	got, err := repo.CreateWorkDay(workDay)
	testutil.AssertNoErr(t, err)

	want := workDay
	want.Id = 1
	want.CreatedAt = time.Now()
	want.UpdatedAt = time.Now()

	testutil.AssertWorkDay(t, got, want)

	gotFromDb, err := repo.GetWorkDayByDate(workDay.Date)
	testutil.AssertNoErr(t, err)
	testutil.AssertWorkDay(t, gotFromDb, want)
}

func TestGetWorkDayByDate(t *testing.T) {
	repo := testutil.NewRepo(t)
	date := time.Date(2023, 8, 11, 0, 0, 0, 0, time.Local)

	t.Run("No results", func(t *testing.T) {
		got, err := repo.GetWorkDayByDate(date)
		testutil.AssertNoErr(t, err)

		var emptyWorkDay model.WorkDay
		if got != emptyWorkDay {
			t.Errorf("Expected empty workday, got %+v\n", got)
		}
	})

	t.Run("Found a work day", func(t *testing.T) {
		want, err := repo.CreateWorkDay(model.WorkDay{Date: date})
		testutil.AssertNoErr(t, err)

		got, err := repo.GetWorkDayByDate(date)
		testutil.AssertNoErr(t, err)
		testutil.AssertWorkDay(t, got, want)
	})
}

func TestGetWorkDayCount(t *testing.T) {
	repo := testutil.NewRepo(t)

	t.Run("No work days", func(t *testing.T) {
		got, err := repo.GetWorkDayCount()
		testutil.AssertNoErr(t, err)

		if got != 0 {
			t.Errorf("Expected 0, got %d\n", got)
		}
	})

	t.Run("2 work days", func(t *testing.T) {
		_, err := repo.CreateWorkDay(model.WorkDay{Date: time.Date(2023, 8, 11, 0, 0, 0, 0, time.Local)})
		testutil.AssertNoErr(t, err)

		_, err = repo.CreateWorkDay(model.WorkDay{Date: time.Date(2023, 8, 12, 0, 0, 0, 0, time.Local)})
		testutil.AssertNoErr(t, err)

		got, err := repo.GetWorkDayCount()
		testutil.AssertNoErr(t, err)

		if got != 2 {
			t.Errorf("Expected 2, got %d\n", got)
		}
	})
}

func TestCreateWorkPeriod(t *testing.T) {
	repo := testutil.NewRepo(t)

	workDay, err := repo.CreateWorkDay(model.WorkDay{Date: util.TodayAtMidnight()})
	testutil.AssertNoErr(t, err)

	workPeriod := model.WorkPeriod{
		WorkDayId: workDay.Id,
		StartAt:   util.TodayAtMidnight().Add(9 * time.Hour),
		EndAt:     sql.NullTime{Time: util.TodayAtMidnight().Add(10 * time.Hour), Valid: true},
	}

	got, err := repo.CreateWorkPeriod(workPeriod)
	testutil.AssertNoErr(t, err)

	want := workPeriod
	want.Id = 1
	want.CreatedAt = time.Now()
	want.UpdatedAt = time.Now()

	testutil.AssertEqualStructs(t, got, want)

	wpsFromDb, err := repo.GetWorkPeriods(workDay)
	testutil.AssertNoErr(t, err)

	if len(wpsFromDb) != 1 {
		t.Errorf("Expected work day to have 1 work period but has %d", len(wpsFromDb))
	}
	testutil.AssertEqualStructs(t, got, wpsFromDb[0])
}

func TestGetWorkPeriods(t *testing.T) {
	repo := testutil.NewRepo(t)
	workDay, err := repo.CreateWorkDay(model.WorkDay{Date: util.TodayAtMidnight()})
	testutil.AssertNoErr(t, err)

	t.Run("no work periods", func(t *testing.T) {
		gotWorkPeriods, err := repo.GetWorkPeriods(workDay)
		testutil.AssertNoErr(t, err)

		if len(gotWorkPeriods) != 0 {
			t.Errorf("Expected 2 work periods, got %d", len(gotWorkPeriods))
		}
	})

	t.Run("2 work periods", func(t *testing.T) {
		workPeriod1, err := repo.CreateWorkPeriod(model.WorkPeriod{WorkDayId: workDay.Id, StartAt: time.Now()})
		testutil.AssertNoErr(t, err)

		workPeriod2, err := repo.CreateWorkPeriod(model.WorkPeriod{WorkDayId: workDay.Id, StartAt: time.Now()})
		testutil.AssertNoErr(t, err)

		gotWorkPeriods, err := repo.GetWorkPeriods(workDay)
		testutil.AssertNoErr(t, err)

		if len(gotWorkPeriods) != 2 {
			t.Errorf("Expected 2 work periods, got %d", len(gotWorkPeriods))
		}

		testutil.AssertEqualStructs(t, workPeriod1, gotWorkPeriods[0])
		testutil.AssertEqualStructs(t, workPeriod2, gotWorkPeriods[1])
	})
}
