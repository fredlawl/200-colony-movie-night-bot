package general

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
	Localization           string
	SuggestionPeriodInDays int
	VotePeriodInDays       int
	MovieNightPeriodInDays int
	DbFilePath             string
}

type Period struct {
	Name     PeriodName
	DaysLeft int
}

type AppSettings struct {
	Config       AppConfig
	CurPeriod    Period
	Localization time.Location
	WeekID       WeekID
	CurDay       time.Time // This is the current day with no time.
}

func DefaultConfiguration() AppConfig {
	return AppConfig{
		Localization:           "America/Chicago",
		SuggestionPeriodInDays: 3,
		VotePeriodInDays:       1,
		MovieNightPeriodInDays: 1,
		DbFilePath:             "./sqlite-cli.db",
	}
}

func CreateAppSettings(cfg AppConfig) (*AppSettings, error) {
	loc, locErr := time.LoadLocation(cfg.Localization)

	if locErr != nil {
		return nil, locErr
	}

	settings := AppSettings{
		Config:       cfg,
		Localization: *loc,
	}

	settings.setTime(time.Now())

	return &settings, nil
}

// Reconfigure settings to a new time. This is especially useful for testing
// purposes.
func (settings *AppSettings) setTime(now time.Time) {
	now = now.In(&settings.Localization)
	settings.CurDay = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0,
		0, &settings.Localization)
	settings.WeekID = WeekIDFromTime(now)
	settings.CurPeriod = calculatePeriod(settings.Config, settings.CurDay)
}

func calculatePeriod(cfg AppConfig, now time.Time) Period {
	curPeriod := Period{
		Name:     Sleep,
		DaysLeft: 0,
	}

	suggestionPeriodStart := now.AddDate(0, 0, -int(now.Weekday()))
	suggestionPeriodEnd := suggestionPeriodStart.AddDate(0, 0, cfg.SuggestionPeriodInDays)
	votingPeriodEnd := suggestionPeriodEnd.AddDate(0, 0, cfg.VotePeriodInDays)
	movieNightPeriodEnd := votingPeriodEnd.AddDate(0, 0, cfg.MovieNightPeriodInDays)

	if now.After(suggestionPeriodStart) && now.Before(suggestionPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.Name = Suggesting
		curPeriod.DaysLeft = int(suggestionPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(suggestionPeriodEnd) && now.Before(votingPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.Name = Voting
		curPeriod.DaysLeft = int(votingPeriodEnd.Sub(now).Hours()) / 24
	}

	if now.After(votingPeriodEnd) && now.Before(movieNightPeriodEnd.AddDate(0, 0, 1)) {
		curPeriod.Name = MovieNight
		curPeriod.DaysLeft = int(movieNightPeriodEnd.Sub(now).Hours()) / 24
	}

	return curPeriod
}
