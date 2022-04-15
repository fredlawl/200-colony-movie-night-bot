package info

import (
	"fmt"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	description := `Shows application state.`

	return &cli.Command{
		Name:        "info",
		Aliases:     []string{"v"},
		Usage:       "shows application state",
		Description: description,
		Action:      infoAction,
	}
}

func infoAction(c *cli.Context) error {
	settings := c.App.Metadata["settings"].(*general.AppSettings)
	c.App.Writer.Write([]byte(fmt.Sprintf("Week: %s\n", settings.WeekID)))
	return nil
}
