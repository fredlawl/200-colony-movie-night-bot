package main

import "time"

type AppState int32

const (
	SUGGESTING AppState = iota
	VOTING
	MOVIE
)

type AppConfig struct {
	suggestionPeriod string
	votePeriod       string
	movieNightPeriod string
	appState         AppState
}

func CreateAppConfig() *AppConfig {
	cfg := AppConfig{
		suggestionPeriod: "R/2021-04-05T00:00/P3D",
		votePeriod:       "P1D",
		movieNightPeriod: "P1D",
		appState:         SUGGESTING,
	}

	cfg.appState = calculateState(cfg)

	return &cfg
}

func calculateState(cfg AppConfig) AppState {
	now := time.Now()
	suggestionDate, err := time.Parse(time.RFC3339, cfg.suggestionPeriod)
	if err != nil {
		//
	}

	if now.After(suggestionDate) {
		return VOTING
	}

	return SUGGESTING
}
