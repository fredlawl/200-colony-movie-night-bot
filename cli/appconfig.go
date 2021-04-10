package main

import (
	"time"
)

type PeriodName int32

const (
	SUGGESTING PeriodName = iota
	VOTING
	MOVIENIGHT
	SLEEP
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
	curPeriod    Period
	localization time.Location
	isoYear      int
	isoWeek      int
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
		localization: *loc,
	}

	now := time.Now().In(&settings.localization)
	settings.curDay = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0,
		0, &settings.localization)
	settings.isoYear, settings.isoWeek = now.ISOWeek()
	settings.curPeriod = calculatePeriod(cfg, settings.curDay)

	return &settings, nil
}

func calculatePeriod(cfg AppConfig, now time.Time) Period {
	curPeriod := Period{
		name:     SLEEP,
		daysLeft: 0,
	}

	suggestionPeriodStart := now.AddDate(0, 0, -int(now.Weekday()))
	suggestionPeriodEnd := suggestionPeriodStart.AddDate(0, 0, cfg.suggestionPeriodInDays)
	votingPeriodEnd := suggestionPeriodEnd.AddDate(0, 0, cfg.votePeriodInDays)
	movieNightPeriodEnd := votingPeriodEnd.AddDate(0, 0, cfg.movieNightPeriodInDays)

	if now.After(suggestionPeriodStart) && now.Before(suggestionPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = SUGGESTING
		curPeriod.daysLeft = int(suggestionPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(suggestionPeriodEnd) && now.Before(votingPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = VOTING
		curPeriod.daysLeft = int(votingPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(votingPeriodEnd) && now.Before(movieNightPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.name = MOVIENIGHT
		curPeriod.daysLeft = int(movieNightPeriodEnd.Sub(now).Hours()) / 24
	}

	return curPeriod
}
