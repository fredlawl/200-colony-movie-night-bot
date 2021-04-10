package main

import (
	"testing"
	"time"
)

func TestGivenASuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 5, 0, 0, 0, 0, &settings.localization)

	expected := SUGGESTING
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 2 {
		t.Fail()
	}
}

func TestGivenAEndSuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 7, 0, 0, 0, 0, &settings.localization)

	expected := SUGGESTING
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAVotingDateStateIsInVoting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 8, 0, 0, 0, 0, &settings.localization)

	expected := VOTING
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAMovieNightDateStateIsInMovienight(t *testing.T) {
	cfg := DefaultConfiguration()
	settings := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 9, 0, 0, 0, 0, &settings.localization)

	expected := MOVIENIGHT
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenASleepDateStateIsInSleep(t *testing.T) {
	cfg := DefaultConfiguration()
	settings := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 10, 0, 0, 0, 0, &settings.localization)

	expected := SLEEP
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}
