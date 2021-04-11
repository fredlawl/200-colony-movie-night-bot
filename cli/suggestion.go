package main

import (
	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v2"
)

type Suggestion struct {
	userName string
	movie    string
}

func suggestMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	if settings.curPeriod.name != SUGGESTING {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to add the movie to suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	db, dbErr := bolt.Open("cli.db", 0600, nil)
	if dbErr != nil {
		return dbErr
	}

	defer db.Close()
	return nil
}

func SuggestionCliCommand() *cli.Command {
	description := `List this weeks suggestions:
    mov suggestions list

Add suggestion:
    mov suggestions suggest "[movie name]"
`

	return &cli.Command{
		Name:        "suggestions",
		Aliases:     []string{"s"},
		Usage:       "manages movie suggestions",
		Description: description,
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Lists suggested movies",
			},
			{
				Name:    "suggest",
				Aliases: []string{"s"},
				Usage:   "Suggest a movie",
				Action:  suggestMovieAction,
			},
		},
	}
}
