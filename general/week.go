package general

import (
	"fmt"
	"strconv"
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

func WeekIDFromString(id string) (*WeekID, error) {
	isoYear, err := strconv.Atoi(id[0:4])
	if err != nil {
		return nil, err
	}

	isoWeek, err := strconv.Atoi(id[4:])
	if err != nil {
		return nil, err
	}

	return &WeekID{
		IsoYear: isoYear,
		IsoWeek: isoWeek,
	}, nil
}
