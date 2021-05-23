package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func SuggestionCliCommand() *cli.Command {
	description := `List this weeks suggestions:
    mov suggestions list

Add suggestion:
    mov suggestions add "[movie name]"

Remove suggestion:
	mov suggestions remove [id]

	Only users may remove their own suggestions.
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
				Action:  listMoviesAction,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Suggest a movie",
				Action:  suggestMovieAction,
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "Remove suggestion",
				Action:  removeMovieAction,
			},
		},
	}
}

func suggestMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	dbSession, err := sql.Open("sqlite3", settings.config.dbFilePath)
	if err != nil {
		return err
	}
	defer dbSession.Close()

	suggestionRepository := NewSuggestionRepository(dbSession)

	if c.NArg() < 1 {
		_, writeErr := c.App.Writer.Write([]byte("Movie name not provided as argument.\n"))
		return writeErr
	}

	if settings.curPeriod.name != Suggesting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to add the movie to suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	suggestion, err := NewSuggestion(settings.weekID, c.String("user"), MovieFromString(c.Args().First()))
	if err != nil {
		c.App.Writer.Write([]byte(err.Error() + "\n"))
		return err
	}

	saveErr := suggestionRepository.Save(*suggestion)
	if saveErr != nil {
		c.App.Writer.Write([]byte(fmt.Sprintf("Movie \"%s\" was already suggested.\n", suggestion.Movie.String())))
		return saveErr
	}

	return nil
}

func listMoviesAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	dbSession, err := sql.Open("sqlite3", settings.config.dbFilePath)
	if err != nil {
		return err
	}
	defer dbSession.Close()

	suggestionRepository := NewSuggestionRepository(dbSession)

	var outputBuffer strings.Builder

	outputBuffer.WriteString(fmt.Sprintf("%-4s%-.32s\n", "ID", "Movie"))

	suggestionRepository.AllSuggestions(settings.weekID, func(k []byte, s *Suggestion) error {
		outputBuffer.WriteString(fmt.Sprintf("%-4d%-.32s\n",
			s.Order,
			s.Movie.String()))
		return nil
	})

	c.App.Writer.Write([]byte(outputBuffer.String()))

	return nil
}

func removeMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	dbSession, err := sql.Open("sqlite3", settings.config.dbFilePath)
	if err != nil {
		return err
	}
	defer dbSession.Close()

	suggestionRepository := NewSuggestionRepository(dbSession)

	orderID, err := strconv.ParseUint(c.Args().First(), 10, 64)
	if err != nil {
		_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a number.\n", c.Args().First())))
		return writeErr
	}

	if settings.curPeriod.name != Suggesting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to remove the movie from suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	// Need to first get a suggestion
	foundSuggestion := suggestionRepository.GetSuggestionByOrder(SuggestionOrderID(orderID))
	if foundSuggestion == nil {
		_, _ = c.App.Writer.Write([]byte("Unable to find a matching suggestion.\n"))
		return err
	}

	// Compare suggestion authors to validate this user can remove suggestion
	if strings.Compare(foundSuggestion.Author, c.String("user")) != 0 {
		_, writeErr := c.App.Writer.Write([]byte("You did not suggest this movie, and can't remove it.\n"))
		return writeErr
	}

	// Remove suggestion
	if removeErr := suggestionRepository.Remove(*foundSuggestion); removeErr != nil {
		_, _ = c.App.Writer.Write([]byte("Unable to remove suggestion from DB.\n"))
		return removeErr
	}

	return nil
}
