package main

import (
	"testing"
	"time"
)

func TestGiveDateWhenExecutedWeekIdGivesFormattedId(t *testing.T) {
	time := time.Date(2021, 4, 10, 0, 0, 0, 0, time.Local)
	actual := WeekIdFromTime(time)

	if actual.String() != "202114" {
		t.Fail()
	}
}

func TestGiveDateWithASingleDigitWeekWhenExecutedWeekIdGivesFormattedId(t *testing.T) {
	time := time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local)
	actual := WeekIdFromTime(time)

	if actual.String() != "202101" {
		t.Fail()
	}
}
