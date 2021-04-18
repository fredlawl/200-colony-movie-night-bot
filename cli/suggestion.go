package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

var SUGGESTION_BUCKET_NAME string = "suggestions"
var SUGGESTION_BUCKET_LOOKUP_NAME string = "lookup"

type SuggestionId string

func (SuggestionId SuggestionId) String() string {
	return string(SuggestionId)
}

type Suggestion struct {
	id     SuggestionId
	weekId WeekId
	author string
	movie  Movie
	order  uint64
}

func NewSuggestion(weekId WeekId, author string, movie Movie) (*Suggestion, error) {
	if len(movie.Encode()) == 0 {
		return nil, errors.New("movie could not be encoded")
	}

	suggestionId := SuggestionId(uuid.New().String())
	return &Suggestion{
		id:     suggestionId,
		weekId: weekId,
		author: author,
		movie:  movie,
		order:  1,
	}, nil
}

func (suggestion *Suggestion) SaveSuggestion(db *bolt.DB) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	weekBucket, err := tx.CreateBucketIfNotExists([]byte(suggestion.weekId.String()))
	if err != nil {
		return err
	}

	suggestionBucket, err := weekBucket.CreateBucketIfNotExists([]byte(SUGGESTION_BUCKET_NAME))
	if err != nil {
		return err
	}

	lookupBucket, err := weekBucket.CreateBucketIfNotExists([]byte(SUGGESTION_BUCKET_LOOKUP_NAME))
	if err != nil {
		return err
	}

	// TODO: Check if movie currently exists prior to save, this will work
	//   based off the suggestion order or movie encoding

	orderId, err := suggestionBucket.NextSequence()
	if err != nil {
		return err
	}
	suggestion.order = orderId

	if buf, err := json.Marshal(suggestion); err != nil {
		return err
	} else if err := suggestionBucket.Put([]byte(suggestion.id.String()), buf); err != nil {
		return err
	}

	// Following two insertions creates two lookups per suggestion
	//  1. Order
	//  2. Movie encoding
	// This allows people to either type the movie name, hash, or order to make a vote
	orderLookupKey := fmt.Sprintf("%s:%s", "order", strconv.FormatUint(orderId, 10))
	if err := lookupBucket.Put([]byte(orderLookupKey), []byte(suggestion.id.String())); err != nil {
		return err
	}

	movieHashLookupKey := fmt.Sprintf("%s:%s", "hash", suggestion.movie.Encode())
	if err := lookupBucket.Put([]byte(movieHashLookupKey), []byte(suggestion.id.String())); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func suggestMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	db, err := bolt.Open("cli.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if c.NArg() < 1 {
		_, writeErr := c.App.Writer.Write([]byte("Movie name not provided as argument.\n"))
		return writeErr
	}

	if settings.curPeriod.name != SUGGESTING && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to add the movie to suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	suggestion, err := NewSuggestion(settings.weekId, c.String("user"), MovieFromString(c.Args().First()))
	if err != nil {
		c.App.Writer.Write([]byte(err.Error() + "\n"))
		return err
	}

	saveErr := suggestion.SaveSuggestion(db)
	if saveErr != nil {
		c.App.Writer.Write([]byte("Unable to save this movie.\n"))
		return saveErr
	}

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
