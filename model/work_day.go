package model

import (
	"database/sql"
	"time"

	"github.com/robyparr/wh/util"
)

type WorkDay struct {
	Id         int
	Date       time.Time
	LengthMins int `db:"length_mins"`
	Note       sql.NullString
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

const DefaultDayLengthMins int = 7.5 * 60

func NewWorkDay(date time.Time) WorkDay {
	return WorkDay{
		Date:       date,
		LengthMins: DefaultDayLengthMins,
	}
}

func NewWorkDayToday() WorkDay {
	return NewWorkDay(util.TodayAtMidnight())
}

func (w *WorkDay) SetNote(note string) {
	w.Note = sql.NullString{String: note, Valid: true}
}
