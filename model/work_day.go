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

const defaultLengthMins = 7.5 * 60

func NewWorkDay() WorkDay {
	return WorkDay{
		Date:       util.TodayAtMidnight(),
		LengthMins: defaultLengthMins,
	}
}
