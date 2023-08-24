package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robyparr/wh/model"
)

const schema string = `
	CREATE TABLE IF NOT EXISTS work_days (
		id 					INTEGER PRIMARY KEY,
		date 				DATETIME NOT NULL,
		length_mins INTEGER NOT NULL,
		note 				TEXT,
		created_at 	DATETIME NOT NULL,
		updated_at 	DATETIME NOT NULL
	);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_work_days_date on work_days(date);

	CREATE TABLE IF NOT EXISTS work_periods (
		id					INTEGER PRIMARY KEY,
		work_day_id	INTEGER NOT NULL,
		start_at		DATETIME NOT NULL,
		end_at			DATETIME,
		note				TEXT,
		created_at	DATETIME NOT NULL,
		updated_at	DATETIME NOT NULL,

		FOREIGN KEY(work_day_id) REFERENCES work_days(id)
	);
`

const DefaultDatabasePath string = "./db.sqlite"

func NewRepo(filepath string) (*Repo, error) {
	db, err := sqlx.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	return &Repo{
		db: db,
	}, nil
}

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) CreateWorkDay(workDay model.WorkDay) (model.WorkDay, error) {
	now := time.Now()
	workDay.CreatedAt = now
	workDay.UpdatedAt = now

	result, err := r.db.NamedExec(`
		INSERT INTO work_days (date, length_mins, note, created_at, updated_at)
		VALUES (:date, :length_mins, :note, :created_at, :updated_at)
	`, workDay)

	if err != nil {
		return model.WorkDay{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.WorkDay{}, err
	}

	workDay.Id = int(id)

	return workDay, nil
}

func (r *Repo) GetWorkDayByDate(date time.Time) (model.WorkDay, error) {
	var workDay model.WorkDay
	if err := r.db.Get(&workDay, "SELECT * FROM work_days WHERE date = ?", date); err != nil {
		if err == sql.ErrNoRows {
			return model.WorkDay{}, nil
		}

		return model.WorkDay{}, err
	}

	return workDay, nil
}

func (r *Repo) GetWorkDayCount() (int, error) {
	var count int
	if err := r.db.Get(&count, "SELECT COUNT(*) FROM work_days;"); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repo) CreateWorkPeriod(period model.WorkPeriod) (model.WorkPeriod, error) {
	now := time.Now()
	period.CreatedAt = now
	period.UpdatedAt = now

	result, err := r.db.NamedExec(`
		INSERT INTO work_periods (work_day_id, start_at, end_at, created_at, updated_at)
		VALUES (:work_day_id, :start_at, :end_at, :created_at, :updated_at)
	`, period)

	if err != nil {
		return model.WorkPeriod{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.WorkPeriod{}, nil
	}

	period.Id = int(id)
	return period, nil
}

func (r *Repo) GetWorkPeriods(workDay model.WorkDay) ([]model.WorkPeriod, error) {
	var periods []model.WorkPeriod
	if err := r.db.Select(&periods, "SELECT * FROM work_periods WHERE work_day_id = ?", workDay.Id); err != nil {
		return []model.WorkPeriod{}, err
	}

	return periods, nil
}
