package model

import (
	"database/sql"
	"time"
)

type WorkPeriod struct {
	Id        int
	WorkDayId int          `db:"work_day_id"`
	StartAt   time.Time    `db:"start_at"`
	EndAt     sql.NullTime `db:"end_at"`
	Note      sql.NullString
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewWorkPeriod(workDay WorkDay) WorkPeriod {
	return WorkPeriod{
		WorkDayId: workDay.Id,
		StartAt:   time.Now(),
	}
}

func (wp *WorkPeriod) SetEndAt(t time.Time) {
	wp.EndAt = sql.NullTime{Valid: !t.IsZero(), Time: t}
}

func (wp *WorkPeriod) SetNote(str string) {
	wp.Note = sql.NullString{Valid: str != "", String: str}
}
