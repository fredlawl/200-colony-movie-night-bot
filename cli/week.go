package main

import (
	"fmt"
	"time"
)

type WeekId struct {
	IsoYear int
	IsoWeek int
}

func (w WeekId) String() string {
	return fmt.Sprintf("%d%02d", w.IsoYear, w.IsoWeek)
}

func WeekIdFromTime(t time.Time) WeekId {
	isoYear, isoWeek := t.ISOWeek()
	return WeekId{
		IsoYear: isoYear,
		IsoWeek: isoWeek,
	}
}
