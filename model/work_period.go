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
