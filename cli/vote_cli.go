package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

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

		if sr.err.(sqlite3.Error).Code == sqlite3.ErrConstraint {
			c.App.Writer.Write([]byte(fmt.Sprintf("Suggestion %d does not exist.\n", sr.vote.SuggestionOrderID)))
		} else {
			c.App.Writer.Write([]byte(fmt.Sprintf("Vote for suggestion %d resulted in an error.\n", sr.vote.SuggestionOrderID)))
		}

		hasErr = true
	}

	if hasErr {
		c.App.Writer.Write([]byte("Unable to cast votes. Something went wrong with the transaction.\n"))
		return errors.New("END of vote bulk save errors")
	}

	return nil
}
