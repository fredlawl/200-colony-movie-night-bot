package vote

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/fredlawl/200-colony-movie-night-bot/suggestion"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
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
	settings := c.App.Metadata["settings"].(*general.AppSettings)
	dbSession := c.App.Metadata["dbSession"].(*sql.DB)
	author := c.String("user")
	week := settings.WeekID

	if settings.CurPeriod.Name != general.Voting && !c.Bool("bypass") {
		_, writeErr := c.App.Writer.Write([]byte("Sorry, unable to cast votes. The vote period has already ended.\n"))
		return writeErr
	}

	voteRepository := NewRepository(dbSession)

	numSuggestions := voteRepository.SuggestionCnt(week)
	if numSuggestions == 0 {
		_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("There are no suggestions this week! Add some :D\n")))
		return writeErr
	}

	uniqueVotes := make(map[suggestion.OrderedID]struct{})
	var emptyMember struct{}
	var votes []Vote
	for i := 0; i < c.NArg(); i++ {
		suggestionOrderIDArg := c.Args().Get(i)
		suggestionOrderID, parseErr := strconv.ParseUint(suggestionOrderIDArg, 10, 64)
		if parseErr != nil {
			_, writeErr := c.App.Writer.Write([]byte(fmt.Sprintf("\"%s\" is not a valid movie ID.\n", suggestionOrderIDArg)))
			return writeErr
		}

		id := suggestion.OrderedID(suggestionOrderID)
		_, exists := uniqueVotes[id]
		if exists {
			continue
		}

		votes = append(votes, Vote{
			VoteID:              ID(uuid.New().String()),
			SuggestionOrderedID: id,
			Author:              author,
			Preference:          uint(i + 1),
			WeekID:              week,
		})

		uniqueVotes[id] = emptyMember
	}

	saveResults, err := voteRepository.BulkSaveVotes(author, week, votes)
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
			c.App.Writer.Write([]byte(fmt.Sprintf("Suggestion %d does not exist.\n", sr.vote.SuggestionOrderedID)))
		} else {
			c.App.Writer.Write([]byte(fmt.Sprintf("Vote for suggestion %d resulted in an error.\n", sr.vote.SuggestionOrderedID)))
		}

		hasErr = true
	}

	if hasErr {
		c.App.Writer.Write([]byte("Unable to cast votes. Something went wrong with the transaction.\n"))
		return errors.New("END of vote bulk save errors")
	}

	return nil
}
