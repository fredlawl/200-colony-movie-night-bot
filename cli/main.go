package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func main() {
	appID := uuid.New().String()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	errorLogName := fmt.Sprintf("logs/%s.error.log", today)

	os.Mkdir("logs", 0700)

	errorLogFile, errorLogFileErr := os.OpenFile(errorLogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if errorLogFileErr != nil {
		log.Fatalf("[error] %s error creating logfile %v", appID, errorLogFileErr)
	}
	defer errorLogFile.Close()

	log.SetOutput(errorLogFile)
	log.SetPrefix(appID + " - ")

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
			// To help with testing, allow app to bypass Period restrictions
			&cli.BoolFlag{
				Name:     "bypass",
				Aliases:  []string{"bp"},
				Usage:    "disable app state check",
				Hidden:   true,
				Required: false,
				Value:    false,
			},
		},
		Commands: []*cli.Command{
			SuggestionCliCommand(),
			VoteCliCommand(),
		},
	}

	cliErr := app.Run(os.Args)
	if cliErr != nil {
		log.Printf("[error] %s",
			strings.Join(os.Args, " "))
		log.Fatalf("[error] %v", cliErr)
	}
}
