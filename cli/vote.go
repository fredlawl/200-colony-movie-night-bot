package main

import "github.com/urfave/cli/v2"

type VoteID string

type Vote struct {
	VoteID       VoteID
	SuggestionID SuggestionID
	Author       string
	Preference   uint
}

func castVotesAction(c *cli.Context) error {
	cfg := DefaultConfiguration()
	_, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		return settingsErr
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
