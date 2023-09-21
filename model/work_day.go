package model

import (
	"database/sql"
	"fmt"
	"log"
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

	workPeriods []WorkPeriod
	timeWorked  *time.Duration
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

func (w *WorkDay) SetWorkPeriods(periods []WorkPeriod) {
	w.workPeriods = periods
	w.timeWorked = nil
}

func (w *WorkDay) TimeWorked() time.Duration {
	if w.timeWorked != nil {
		return *w.timeWorked
	}

	var timeWorked time.Duration
	for _, wp := range w.workPeriods {
		timeWorked += wp.TimeWorked()
	}

	w.timeWorked = &timeWorked
	return timeWorked
}

func (w *WorkDay) TimeRemaining() time.Duration {
	return w.Length() - w.TimeWorked()
}

func (w *WorkDay) EstimatedFinish() time.Time {
	return time.Now().Truncate(time.Minute).Add(w.TimeRemaining())
}

func (w *WorkDay) Length() time.Duration {
	duration, err := time.ParseDuration(fmt.Sprintf("%dm", w.LengthMins))
	if err != nil {
		log.Fatalf("error parsing day length: %v\n", err)
	}

	return duration
}
