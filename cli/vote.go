package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type VoteID string

type Vote struct {
	VoteID            VoteID
	SuggestionOrderID SuggestionOrderID
	WeekID            WeekID
	Author            string
	Preference        uint
}

type VoteRepository struct {
	session *sql.DB
}

type BulkVoteResult struct {
	err  error
	vote Vote
}

func NewVoteRepository(session *sql.DB) *VoteRepository {
	return &VoteRepository{
		session: session,
	}
}

func (context *VoteRepository) BulkSaveVotes(votes []Vote) ([]BulkVoteResult, error) {
	emptyBulkResult := []BulkVoteResult{}

	context.session.Exec("PRAGMA foreign_keys = ON;")

	tx, err := context.session.Begin()
	if err != nil {
		return emptyBulkResult, err
	}

	stmt, err := context.session.Prepare(`
		INSERT INTO votes (suggestionID, weekID, author, preference)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return emptyBulkResult, tx.Rollback()
	}

	var hasErrors = false
	var bulkResults = make([]BulkVoteResult, len(votes))
	for i, v := range votes {
		bulkResults[i].vote = v
		_, bulkResults[i].err = tx.Stmt(stmt).Exec(v.SuggestionOrderID,
			v.WeekID.String(), v.Author, v.Preference)
	}

	var txErr error
	if hasErrors {
		txErr = tx.Rollback()
	} else {
		txErr = tx.Commit()
	}

	context.session.Exec("PRAGMA foreign_keys = OFF;")

	return bulkResults, txErr
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
			_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a valid movie ID.\n", suggestionOrderIDArg)))
			return writeErr
		}

		votes[i] = Vote{
			VoteID:            VoteID(uuid.New().String()),
			SuggestionOrderID: SuggestionOrderID(suggestionOrderID),
			Author:            c.String("user"),
			Preference:        uint(i + 1),
			WeekID:            settings.weekID,
		}
	}

	dbSession, err := sql.Open("sqlite3", settings.config.dbFilePath)
	if err != nil {
		return err
	}
	defer dbSession.Close()

	voteRepository := NewVoteRepository(dbSession)

	saveResults, err := voteRepository.BulkSaveVotes(votes)
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
    mov votes cast [Suggestion ID 1], [Suggestion ID 2], ... [Suggestion ID N]

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
