package main

import (
	"fmt"
	"time"
)

type WeekID struct {
	IsoYear int
	IsoWeek int
}

func (w WeekID) String() string {
	return fmt.Sprintf("%d%02d", w.IsoYear, w.IsoWeek)
}

func WeekIDFromTime(t time.Time) WeekID {
	isoYear, isoWeek := t.ISOWeek()
	return WeekID{
		IsoYear: isoYear,
		IsoWeek: isoWeek,
	}
}
