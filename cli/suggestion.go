package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	orderLookupKey := context.OrderLookupKey(orderId)
	if err := lookupBucket.Put([]byte(orderLookupKey), []byte(s.Id.String())); err != nil {
		return err
	}

	movieHashLookupKey := context.MovieHashLookupKey(s.Movie)
	if err := lookupBucket.Put([]byte(movieHashLookupKey), []byte(s.Id.String())); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (context *SuggestionPersistanceContext) AllSuggestions(callback func(key []byte, suggestion *Suggestion) error) error {
	return context.db.View(func(tx *bolt.Tx) error {
		weekBucket := tx.Bucket([]byte(context.weekId.String()))
		suggestionsBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_NAME))
		cursor := suggestionsBucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var suggestion Suggestion
			unmarshalErr := json.Unmarshal(v, &suggestion)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			err := callback(k, &suggestion)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// OrderLookupKey Given a orderId number, returns order:[number]
func (context *SuggestionPersistanceContext) OrderLookupKey(orderID uint64) string {
	return fmt.Sprintf("%s:%s", "order", strconv.FormatUint(orderID, 10))
}

// MovieHashLookupKey Given a movie, returns move:[moviehash]
func (context *SuggestionPersistanceContext) MovieHashLookupKey(movie Movie) string {
	return fmt.Sprintf("%s:%s", "hash", movie.Encode())
}

// GetSuggestionByOrder Given the order id, return the suggestion at that position
func (context *SuggestionPersistanceContext) GetSuggestionByOrder(orderID uint64) (*Suggestion, error) {
	var suggestion Suggestion

	err := context.db.View(func(tx *bolt.Tx) error {
		weekBucket := tx.Bucket([]byte(context.weekId.String()))
		lookupBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_LOOKUP_NAME))
		suggestionsBucket := weekBucket.Bucket([]byte(SUGGESTION_BUCKET_NAME))

		key := lookupBucket.Get([]byte(context.OrderLookupKey(orderID)))
		value := suggestionsBucket.Get([]byte(key))

		unmarshalErr := json.Unmarshal(value, &suggestion)
		if unmarshalErr != nil {
			return unmarshalErr
		}

		return nil
	})

	return &suggestion, err
}

// Close Closes the persistence context
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

	var outputBuffer strings.Builder

	outputBuffer.WriteString(fmt.Sprintf("%-4s%-.32s\n", "ID", "Movie"))

	listerr := db.AllSuggestions(func(k []byte, s *Suggestion) error {
		outputBuffer.WriteString(fmt.Sprintf("%-4d%-.32s\n",
			s.Order,
			s.Movie.String()))
		return nil
	})

	if listerr != nil {
		return listerr
	}

	c.App.Writer.Write([]byte(outputBuffer.String()))

	return nil
}

func removeMovieAction(c *cli.Context) error {
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

	orderID, err := strconv.ParseUint(c.Args().First(), 10, 64)
	if err != nil {
		_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a number.\n", c.Args().First())))
		return writeErr
	}

	if settings.curPeriod.name != SUGGESTING && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to remove the movie from suggestions. The suggestion period has already ended.\n"))
		return writeErr
	}

	// Need to first get a suggestion
	foundSuggestion, err := db.GetSuggestionByOrder(orderID)
	if err != nil {
		_, _ = c.App.Writer.Write([]byte("Unable to find a matching suggestion.\n"))
		return err
	}

	// Compare suggestion authors to validate this user can remove suggestion
	if strings.Compare(foundSuggestion.Author, c.String("user")) != 0 {
		_, writeErr := c.App.Writer.Write([]byte("You did not suggest this movie, and can't remove it.\n"))
		return writeErr
	}

	// Remove suggestion

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
				Action:  listMoviesAction,
			},
			{
				Name:    "suggest",
				Aliases: []string{"s"},
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
