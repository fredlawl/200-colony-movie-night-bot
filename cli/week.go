package main

import (
	"fmt"
	"time"
)

type WeekId struct {
	isoYear int
	isoWeek int
}

func (w WeekId) String() string {
	return fmt.Sprintf("%d%02d", w.isoYear, w.isoWeek)
}

func WeekIdFromTime(t time.Time) WeekId {
	isoYear, isoWeek := t.ISOWeek()
	return WeekId{
		isoYear: isoYear,
		isoWeek: isoWeek,
	}
}
