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

type SuggestionPersistanceContext struct {
	db     *bolt.DB
	weekId WeekId
}

func NewSuggestionPersistance(week WeekId) (*SuggestionPersistanceContext, error) {
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

	return &SuggestionPersistanceContext{
		db:     db,
		weekId: week,
	}, nil
}

func (context *SuggestionPersistanceContext) Save(s Suggestion) error {
	tx, err := context.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	weekBucket := tx.Bucket([]byte(s.WeekId.String()))
	suggestionBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_NAME))
	lookupBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_LOOKUP_NAME))

	// TODO: Check if movie currently exists prior to save, this will work
	//   based off the suggestion order or movie encoding

	orderId, err := suggestionBucket.NextSequence()
	if err != nil {
		return err
	}
	s.Order = orderId

	if buf, err := json.Marshal(s); err != nil {
		return err
	} else if err := suggestionBucket.Put([]byte(s.Id.String()), buf); err != nil {
		return err
	}

	// Following two insertions creates two lookups per suggestion
	//  1. Order
	//  2. Movie encoding
	// This allows people to either type the movie name, hash, or order to make a vote
	orderLookupKey := fmt.Sprintf("%s:%s", "order", strconv.FormatUint(orderId, 10))
	if err := lookupBucket.Put([]byte(orderLookupKey), []byte(s.Id.String())); err != nil {
		return err
	}

	movieHashLookupKey := fmt.Sprintf("%s:%s", "hash", s.Movie.Encode())
	if err := lookupBucket.Put([]byte(movieHashLookupKey), []byte(s.Id.String())); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (context *SuggestionPersistanceContext) AllSuggestions(callback func(cursor *bolt.Cursor) error) error {
	return context.db.View(func(tx *bolt.Tx) error {
		weekBucket := tx.Bucket([]byte(context.weekId.String()))
		suggestionsBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_NAME))
		suggestionsCursor := suggestionsBucket.Cursor()
		return callback(suggestionsCursor)
	})
}

func (context *SuggestionPersistanceContext) Close() {
	context.db.Close()
}

func suggestMovieAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	db, err := NewSuggestionPersistance(settings.weekId)
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

	saveErr := db.Save(*suggestion)
	if saveErr != nil {
		c.App.Writer.Write([]byte("Unable to save this movie.\n"))
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

	db, err := NewSuggestionPersistance(settings.weekId)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.AllSuggestions(func(cursor *bolt.Cursor) error {
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var suggestion Suggestion
			unmarshalErr := json.Unmarshal(v, &suggestion)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
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
				Action:  listMoviesAction,
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
