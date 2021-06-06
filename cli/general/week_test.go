package general

import (
	"testing"
	"time"
)

func TestGiveDateWhenExecutedWeekIdGivesFormattedId(t *testing.T) {
	time := time.Date(2021, 4, 10, 0, 0, 0, 0, time.Local)
	actual := WeekIDFromTime(time)

	if actual.String() != "202114" {
		t.Fail()
	}
}

func TestGiveDateWithASingleDigitWeekWhenExecutedWeekIdGivesFormattedId(t *testing.T) {
	time := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
	actual := WeekIDFromTime(time)

	if actual.String() != "202101" {
		t.Fail()
	}
}

func TestGivenWeekIDStringConversionWorksCorrectly(t *testing.T) {
	expected := &WeekID{
		IsoYear: 2021,
		IsoWeek: 20,
	}

	actual, err := WeekIDFromString("202120")
	if err != nil {
		t.Fail()
	}

	if !(expected.IsoYear == actual.IsoYear || expected.IsoWeek == actual.IsoWeek) {
		t.Fail()
	}
}
