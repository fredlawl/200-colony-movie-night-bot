package general

import (
	"testing"
	"time"
)

func TestGivenASuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 5, 0, 0, 0, 0, &settings.Localization)

	expected := Suggesting
	actual := calculatePeriod(cfg, now)

	if expected != actual.Name || actual.DaysLeft != 2 {
		t.Fail()
	}
}

func TestGivenAEndSuggestionPeriodDateStateIsInSuggesting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 7, 0, 0, 0, 0, &settings.Localization)

	expected := Suggesting
	actual := calculatePeriod(cfg, now)

	if expected != actual.Name || actual.DaysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAVotingDateStateIsInVoting(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 8, 0, 0, 0, 0, &settings.Localization)

	expected := Voting
	actual := calculatePeriod(cfg, now)

	if expected != actual.Name || actual.DaysLeft != 0 {
		t.Fail()
	}
}

func TestGivenAMovieNightDateStateIsInMovienight(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 9, 0, 0, 0, 0, &settings.Localization)

	expected := MovieNight
	actual := calculatePeriod(cfg, now)

	if expected != actual.Name || actual.DaysLeft != 0 {
		t.Fail()
	}
}

func TestGivenASleepDateStateIsInSleep(t *testing.T) {
	cfg := DefaultConfiguration()
	settings, _ := CreateAppSettings(cfg)
	now := time.Date(2021, 4, 10, 0, 0, 0, 0, &settings.Localization)

	expected := Sleep
	actual := calculatePeriod(cfg, now)

	if expected != actual.Name || actual.DaysLeft != 0 {
		t.Fail()
	}
}
