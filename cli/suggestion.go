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

const SUGGESTION_BUCKET_NAME string = "suggestions"
const SUGGESTION_BUCKET_LOOKUP_NAME string = "lookup"

type SuggestionId string

func (SuggestionId SuggestionId) String() string {
	return string(SuggestionId)
}

type Suggestion struct {
	Id     SuggestionId
	WeekId WeekId
	Author string
	Movie  Movie
	Order  uint64
}

func NewSuggestion(weekId WeekId, author string, movie Movie) (*Suggestion, error) {
	if len(movie.Encode()) == 0 {
		return nil, errors.New("movie could not be encoded")
	}

	suggestionId := SuggestionId(uuid.New().String())
	return &Suggestion{
		Id:     suggestionId,
		WeekId: weekId,
		Author: author,
		Movie:  movie,
		Order:  1,
	}, nil
}

func (suggestion *Suggestion) SaveSuggestion(db *bolt.DB) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	weekBucket, err := tx.CreateBucketIfNotExists([]byte(suggestion.WeekId.String()))
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
	suggestion.Order = orderId

	if buf, err := json.Marshal(suggestion); err != nil {
		return err
	} else if err := suggestionBucket.Put([]byte(suggestion.Id.String()), buf); err != nil {
		return err
	}

	// Following two insertions creates two lookups per suggestion
	//  1. Order
	//  2. Movie encoding
	// This allows people to either type the movie name, hash, or order to make a vote
	orderLookupKey := fmt.Sprintf("%s:%s", "order", strconv.FormatUint(orderId, 10))
	if err := lookupBucket.Put([]byte(orderLookupKey), []byte(suggestion.Id.String())); err != nil {
		return err
	}

	movieHashLookupKey := fmt.Sprintf("%s:%s", "hash", suggestion.Movie.Encode())
	if err := lookupBucket.Put([]byte(movieHashLookupKey), []byte(suggestion.Id.String())); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func openDb(week WeekId) (*bolt.DB, error) {
	db, err := bolt.Open("cli.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	dbErr := db.Update(func(tx *bolt.Tx) error {
		weekBucket, err := tx.CreateBucketIfNotExists([]byte(week.String()))
		if err != nil {
			return err
		}

		_, serr := weekBucket.CreateBucketIfNotExists([]byte(SUGGESTION_BUCKET_NAME))
		if serr != nil {
			return serr
		}

		_, lerr := weekBucket.CreateBucketIfNotExists([]byte(SUGGESTION_BUCKET_LOOKUP_NAME))
		if lerr != nil {
			return lerr
		}

		return nil
	})

	if dbErr != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func suggestMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	db, err := openDb(settings.weekId)
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
