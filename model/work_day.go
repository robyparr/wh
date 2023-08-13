package model

import (
	"database/sql"
	"time"
)

type WorkDay struct {
	Id         int
	Date       time.Time
	LengthMins int `db:"length_mins"`
	Note       sql.NullString
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

const defaultLengthMins = 7.5 * 60

func NewWorkDay(date time.Time) WorkDay {
	return WorkDay{
		Date:       date,
		LengthMins: defaultLengthMins,
	}
}
