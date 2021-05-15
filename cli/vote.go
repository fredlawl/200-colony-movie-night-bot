package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type VoteID string

type Vote struct {
	VoteID            VoteID
	SuggestionOrderID SuggestionOrderID
	Author            string
	Preference        uint
}

const VoteBucketName string = "votes"

type VotePersistanceContext struct {
	db     *bolt.DB
	weekID WeekID
}

type BulkVoteResult struct {
	err  error
	vote Vote
}

func NewVotePersistance(week WeekID) (*VotePersistanceContext, error) {
	db, err := bolt.Open("cli.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	dbErr := db.Update(func(tx *bolt.Tx) error {
		weekBucket, err := tx.CreateBucketIfNotExists([]byte(week.String()))
		if err != nil {
			return err
		}

		_, serr := weekBucket.CreateBucketIfNotExists([]byte(VoteBucketName))
		if serr != nil {
			return serr
		}

		return nil
	})

	if dbErr != nil {
		db.Close()
		return nil, err
	}

	return &VotePersistanceContext{
		db:     db,
		weekID: week,
	}, nil
}

func (context *VotePersistanceContext) BulkSaveVotes(votes []Vote) ([]BulkVoteResult, error) {
	emptyBulkResult := []BulkVoteResult{}

	tx, err := context.db.Begin(true)
	if err != nil {
		return emptyBulkResult, err
	}
	defer tx.Rollback()

	weekBucket := tx.Bucket([]byte(context.weekID.String()))
	voteBucket := weekBucket.Bucket([]byte(VoteBucketName))

	var hasErrors = false
	var bulkResults = make([]BulkVoteResult, len(votes))
	for i, v := range votes {
		bulkResults[i].vote = v

		if buf, err := json.Marshal(v); err != nil {
			bulkResults[i].err = err
			hasErrors = true
		} else if err := voteBucket.Put([]byte(v.VoteID), buf); err != nil {
			bulkResults[i].err = err
			hasErrors = true
		}
	}

	if hasErrors {
		return bulkResults, tx.Rollback()
	}

	return bulkResults, tx.Commit()
}

// Close Closes the persistence context
func (context *VotePersistanceContext) Close() {
	context.db.Close()
}

func castVotesAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
	}

	if c.NArg() < 1 {
		_, writeErr := c.App.Writer.Write([]byte("You must make at least one vote.\n"))
		return writeErr
	}

	if settings.curPeriod.name != Voting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to cast votes. The vote period has already ended.\n"))
		return writeErr
	}

	votes := make([]Vote, c.NArg())
	for i := 0; i < len(votes); i++ {
		suggestionOrderIDArg := c.Args().Get(i)
		suggestionOrderID, parseErr := strconv.ParseUint(suggestionOrderIDArg, 10, 64)
		if parseErr != nil {
			_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a valid movie id.\n", suggestionOrderIDArg)))
			return writeErr
		}

		votes[i] = Vote{
			VoteID:            VoteID(uuid.New().String()),
			SuggestionOrderID: SuggestionOrderID(suggestionOrderID),
			Author:            c.String("user"),
			Preference:        uint(i + 1),
		}
	}

	db, err := NewVotePersistance(settings.weekID)
	if err != nil {
		return err
	}
	defer db.Close()

	saveResults, err := db.BulkSaveVotes(votes)
	if err != nil {
		c.App.Writer.Write([]byte("Unable to save votes. Something went wrong with the transaction.\n"))
		return err
	}

	var hasErr = false
	for _, sr := range saveResults {
		if sr.err == nil {
			continue
		}

		log.Printf("[error] %v", sr.err)
		c.App.Writer.Write([]byte(fmt.Sprintf("Vote for suggestion %d resulted in an error.\n", sr.vote.SuggestionOrderID)))
		hasErr = true
	}

	if hasErr {
		c.App.Writer.Write([]byte("Unable to save votes. Something went wrong with the transaction.\n"))
		return errors.New("END of vote bulk save errors")
	}

	return nil
}

func VoteCliCommand() *cli.Command {
	description := `Vote for a movie in order of preference:
    mov votes cast [SuggestionID 1], [SuggestionID 2], ... [SuggestionID N]

	To recast votes, this command must be written again. All previous votes will be nullified and replaced with this new order.
`

	return &cli.Command{
		Name:        "votes",
		Aliases:     []string{"v"},
		Usage:       "manages movie votes",
		Description: description,
		Subcommands: []*cli.Command{
			{
				Name:    "cast",
				Aliases: []string{"c"},
				Usage:   "Casts votes for for movies",
				Action:  castVotesAction,
			},
		},
	}
}
