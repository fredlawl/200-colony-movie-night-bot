package main

import (
	"testing"
	"time"
)

func TestGivenASuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := CreateAppConfig()
	now := time.Date(2021, 4, 5, 0, 0, 0, 0, &cfg.timeLocation)

	expected := SUGGESTING
	actual, daysLeft := calculateState(*cfg, now)

	if expected != actual || daysLeft != 2 {
		t.Fail()
	}
}

func TestGivenAEndSuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := CreateAppConfig()
	now := time.Date(2021, 4, 7, 0, 0, 0, 0, &cfg.timeLocation)

	expected := SUGGESTING
	actual, daysLeft := calculateState(*cfg, now)

	if expected != actual || daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAVotingDateStateIsInVoting(t *testing.T) {
	cfg := CreateAppConfig()
	now := time.Date(2021, 4, 8, 0, 0, 0, 0, &cfg.timeLocation)

	expected := VOTING
	actual, daysLeft := calculateState(*cfg, now)

	if expected != actual || daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAMovieNightDateStateIsInMovienight(t *testing.T) {
	cfg := CreateAppConfig()
	now := time.Date(2021, 4, 9, 0, 0, 0, 0, &cfg.timeLocation)

	expected := MOVIENIGHT
	actual, daysLeft := calculateState(*cfg, now)

	if expected != actual || daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenASleepDateStateIsInSleep(t *testing.T) {
	cfg := CreateAppConfig()
	now := time.Date(2021, 4, 10, 0, 0, 0, 0, &cfg.timeLocation)

	expected := SLEEP
	actual, daysLeft := calculateState(*cfg, now)

	if expected != actual || daysLeft != 0 {
		t.Fail()
	}
}
