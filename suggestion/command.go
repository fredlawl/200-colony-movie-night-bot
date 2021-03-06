package suggestion

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
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
	settings := c.App.Metadata["settings"].(*general.AppSettings)
	dbSession := c.App.Metadata["dbSession"].(*sql.DB)

	suggestionRepository := NewRepository(dbSession)

	if c.NArg() < 1 {
		_, writeErr := c.App.Writer.Write([]byte("Movie name not provided as argument.\n"))
		return writeErr
	}

	if settings.CurPeriod.Name != general.Suggesting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to add the movie to suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	suggestion, err := NewSuggestion(settings.WeekID, c.String("user"), general.MovieFromString(c.Args().First()))
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
	settings := c.App.Metadata["settings"].(*general.AppSettings)
	dbSession := c.App.Metadata["dbSession"].(*sql.DB)

	suggestionRepository := NewRepository(dbSession)

	var outputBuffer strings.Builder

	outputBuffer.WriteString(fmt.Sprintf("%-4s%-.32s\n", "ID", "Movie"))

	suggestionRepository.AllSuggestions(settings.WeekID, func(k []byte, s *Suggestion) error {
		outputBuffer.WriteString(fmt.Sprintf("%-4d%-.32s\n",
			s.Order,
			s.Movie.String()))
		return nil
	})

	c.App.Writer.Write([]byte(outputBuffer.String()))

	return nil
}

func removeMovieAction(c *cli.Context) error {
	settings := c.App.Metadata["settings"].(*general.AppSettings)
	dbSession := c.App.Metadata["dbSession"].(*sql.DB)

	suggestionRepository := NewRepository(dbSession)

	orderID, err := strconv.ParseUint(c.Args().First(), 10, 64)
	if err != nil {
		_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a number.\n", c.Args().First())))
		return writeErr
	}

	if settings.CurPeriod.Name != general.Suggesting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to remove the movie from suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	// Need to first get a suggestion
	foundSuggestion := suggestionRepository.GetSuggestionByOrder(OrderedID(orderID))
	if foundSuggestion == nil {
		_, _ = c.App.Writer.Write([]byte("Unable to find a matching suggestion.\n"))
		return nil
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
