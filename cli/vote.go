package main

import (
	"fmt"
	"strconv"

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

	// TODO: Bulk create votes

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
