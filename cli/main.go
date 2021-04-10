package main

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli/v2"
)

// TODO: Use https://github.com/boltdb/bolt for embeded DB
// TODO: Use github.com/urfave/cli for cli configuration

func main() {
	cfg := DefaultConfiguration()
	settings, settingsErr := CreateAppSettings(cfg)

	if settingsErr != nil {
		log.Fatal(settingsErr)
	}

	_ = settings

	db, dbErr := bolt.Open("cli.db", 0600, nil)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	app := &cli.App{
		Name:     "mov",
		HelpName: "mov",
		Usage:    "an application to manage movie night movie suggestions and votes.",
		Flags: []cli.Flag{
			// Because this application depends on multiple users, we need
			// to supply the user calling these commands. Further, the option
			// is hidden so that it does not appear in documenation so a user
			// doesn't have to supply it, but the bot will behind the scenes.
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "user interfacing with the application",
				Hidden:   true,
				Required: true,
			},
		},
		Commands: []*cli.Command{
			SuggestionCliCommand(),
		},
	}

	cliErr := app.Run(os.Args)
	if cliErr != nil {
		log.Fatal(cliErr)
	}

	defer db.Close()
}
