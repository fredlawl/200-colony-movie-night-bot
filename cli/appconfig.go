package main

import (
	"time"
)

type PeriodName int32

const (
	Suggesting PeriodName = iota
	Voting
	MovieNight
	Sleep
)

// Configuration is based on a single 7 day week, with the seggestion
// date starting on Monday.
type AppConfig struct {
	localization           string
	suggestionPeriodInDays int
	votePeriodInDays       int
	movieNightPeriodInDays int
}

type Period struct {
	name     PeriodName
	daysLeft int
}

type AppSettings struct {
	config       AppConfig
	curPeriod    Period
	localization time.Location
	weekID       WeekID
	curDay       time.Time // This is the current day with no time.
}

func DefaultConfiguration() AppConfig {
	return AppConfig{
		localization:           "America/Chicago",
		suggestionPeriodInDays: 3,
		votePeriodInDays:       1,
		movieNightPeriodInDays: 1,
	}
}

func CreateAppSettings(cfg AppConfig) (*AppSettings, error) {
	loc, locErr := time.LoadLocation(cfg.localization)

	if locErr != nil {
		return nil, locErr
	}

	settings := AppSettings{
		config:       cfg,
		localization: *loc,
	}

	settings.setTime(time.Now())

	return &settings, nil
}

// Reconfigure settings to a new time. This is especially useful for testing
// purposes.
func (settings *AppSettings) setTime(now time.Time) {
	now = now.In(&settings.localization)
	settings.curDay = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0,
		0, &settings.localization)
	settings.weekID = WeekIDFromTime(now)
	settings.curPeriod = calculatePeriod(settings.config, settings.curDay)
}

func calculatePeriod(cfg AppConfig, now time.Time) Period {
	curPeriod := Period{
		name:     Sleep,
		daysLeft: 0,
	}

	suggestionPeriodStart := now.AddDate(0, 0, -int(now.Weekday()))
	suggestionPeriodEnd := suggestionPeriodStart.AddDate(0, 0, cfg.suggestionPeriodInDays)
	votingPeriodEnd := suggestionPeriodEnd.AddDate(0, 0, cfg.votePeriodInDays)
	movieNightPeriodEnd := votingPeriodEnd.AddDate(0, 0, cfg.movieNightPeriodInDays)

	if now.After(suggestionPeriodStart) && now.Before(suggestionPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = Suggesting
		curPeriod.daysLeft = int(suggestionPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(suggestionPeriodEnd) && now.Before(votingPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = Voting
		curPeriod.daysLeft = int(votingPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(votingPeriodEnd) && now.Before(movieNightPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = MovieNight
		curPeriod.daysLeft = int(movieNightPeriodEnd.Sub(now).Hours()) / 24
	}

	return curPeriod
}
