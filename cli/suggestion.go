package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

type SuggestionID string
type SuggestionOrderID uint64

func (SuggestionId SuggestionID) String() string {
	return string(SuggestionId)
}

type Suggestion struct {
	ID     SuggestionID
	WeekID WeekID
	Author string
	Movie  Movie
	Order  SuggestionOrderID
}

func NewSuggestion(weekID WeekID, author string, movie Movie) (*Suggestion, error) {
	if len(movie.Encode()) == 0 {
		return nil, errors.New("movie could not be encoded")
	}

	suggestionID := SuggestionID(uuid.New().String())
	return &Suggestion{
		ID:     suggestionID,
		WeekID: weekID,
		Author: author,
		Movie:  movie,
		Order:  1,
	}, nil
}

type SuggestionRepository struct {
	session *sql.DB
}

func NewSuggestionRepository(session *sql.DB) *SuggestionRepository {
	return &SuggestionRepository{
		session: session,
	}
}

func (context *SuggestionRepository) Save(s Suggestion) error {
	stmt, err := context.session.Prepare(
		`INSERT INTO suggestions (
			uuid,
			weekID,
			author,
			movie,
			movieHash
		) VALUES (
			?,
			?,
			?,
			?,
			?
		)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.ID.String(), s.WeekID.String(), s.Author,
		s.Movie.String(), s.Movie.Encode())

	return err
}

func (context *SuggestionRepository) AllSuggestions(weekID WeekID, callback func(key []byte, suggestion *Suggestion) error) {
	stmt, err := context.session.Prepare("SELECT id, uuid, author, movie FROM suggestions WHERE weekID = ? ORDER BY id ASC")
	if err != nil {
		return
	}

	rows, err := stmt.Query(weekID.String())
	if err != nil {
		return
	}

	var id int
	var suggestionID string
	var author string
	var movie string

	for rows.Next() {
		err = rows.Scan(&id, &suggestionID, &author, &movie)
		if err != nil {
			return
		}

		err = callback(
			[]byte(suggestionID),
			&Suggestion{
				ID:     SuggestionID(suggestionID),
				WeekID: weekID,
				Author: author,
				Movie:  MovieFromString(movie),
				Order:  SuggestionOrderID(id),
			})

		if err != nil {
			return
		}
	}
}

// GetSuggestionByOrder Given the order id, return the suggestion at that position
func (context *SuggestionRepository) GetSuggestionByOrder(orderID SuggestionOrderID) *Suggestion {
	stmt, err := context.session.Prepare("SELECT id, uuid, weekID, author, movie FROM suggestions WHERE id = ?")
	if err != nil {
		return nil
	}

	row := stmt.QueryRow(orderID)

	var id int
	var suggestionID string
	var weekID string
	var author string
	var movie string

	err = row.Scan(&id, &suggestionID, &weekID, &author, &movie)
	if err != nil {
		return nil
	}

	parsedWeekID, _ := WeekIDFromString(weekID)

	return &Suggestion{
		ID:     SuggestionID(suggestionID),
		WeekID: *parsedWeekID,
		Author: author,
		Movie:  MovieFromString(movie),
		Order:  SuggestionOrderID(id),
	}
}

func (context *SuggestionRepository) Remove(s Suggestion) error {
	context.session.Exec("PRAGMA foreign_keys = ON;")
	stmt, err := context.session.Prepare("DELETE FROM suggestions WHERE uuid = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.ID.String())

	context.session.Exec("PRAGMA foreign_keys = OFF;")
	return err
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
