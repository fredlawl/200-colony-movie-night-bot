package main

import (
	"testing"
	"time"
)

func TestGivenASuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 5, 0, 0, 0, 0, &settings.localization)

	expected := Suggesting
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 2 {
		t.Fail()
	}
}

func TestGivenAEndSuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 7, 0, 0, 0, 0, &settings.localization)

	expected := Suggesting
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAVotingDateStateIsInVoting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 8, 0, 0, 0, 0, &settings.localization)

	expected := Voting
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAMovieNightDateStateIsInMovienight(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 9, 0, 0, 0, 0, &settings.localization)

	expected := MovieNight
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}

func TestGivenASleepDateStateIsInSleep(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 10, 0, 0, 0, 0, &settings.localization)

	expected := Sleep
	actual := calculatePeriod(cfg, now)

	if expected != actual.name || actual.daysLeft != 0 {
		t.Fail()
	}
}
