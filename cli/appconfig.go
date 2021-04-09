package main

import (
	"time"
)

type AppState int32

const (
	SUGGESTING AppState = iota
	VOTING
	MOVIENIGHT
	SLEEP
)

type AppConfig struct {
	suggestionDayStart     time.Weekday
	suggestionPeriodInDays int
	votePeriodInDays       int
	movieNightPeriodInDays int
	appState               AppState
	timeLocation           time.Location
	isoYear                int
	isoWeek                int
	daysLeftInCurrentState int
	now                    time.Time
}

func CreateAppConfig() *AppConfig {
	loc, _ := time.LoadLocation("America/Chicago")
	cfg := AppConfig{
		suggestionDayStart:     time.Monday,
		suggestionPeriodInDays: 3,
		votePeriodInDays:       1,
		movieNightPeriodInDays: 1,
		appState:               SLEEP,
		timeLocation:           *loc,
	}

	now := time.Now().In(loc)
	cfg.now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	cfg.appState, cfg.daysLeftInCurrentState = calculateState(cfg, cfg.now)
	cfg.isoYear, cfg.isoWeek = now.ISOWeek()

	return &cfg
}

func calculateState(cfg AppConfig, now time.Time) (curState AppState, daysLeft int) {
	suggestionPeriodStart := now.AddDate(0, 0, -int(now.Weekday()))
	suggestionPeriodEnd := suggestionPeriodStart.AddDate(0, 0, cfg.suggestionPeriodInDays)
	votingPeriodEnd := suggestionPeriodEnd.AddDate(0, 0, cfg.votePeriodInDays)
	movieNightPeriodEnd := votingPeriodEnd.AddDate(0, 0, cfg.movieNightPeriodInDays)

	if now.After(suggestionPeriodStart) && now.Before(suggestionPeriodEnd.AddDate(0, 0, 1)) {
		return SUGGESTING, int(suggestionPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(suggestionPeriodEnd) && now.Before(votingPeriodEnd.AddDate(0, 0, 1)) {
		return VOTING, int(votingPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(votingPeriodEnd) && now.Before(movieNightPeriodEnd.AddDate(0, 0, 1)) {
		return MOVIENIGHT, int(movieNightPeriodEnd.Sub(now).Hours()) / 24
	}

	return SLEEP, 0
}
