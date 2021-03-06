package runner

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fredlawl/200-colony-movie-night-bot/general"
	"github.com/fredlawl/200-colony-movie-night-bot/suggestion"
	"github.com/fredlawl/200-colony-movie-night-bot/vote"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func Run(args []string) {
	appID := uuid.New().String()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	errorLogName := fmt.Sprintf("logs/%s.error.log", today)

	os.Mkdir("logs", 0700)

	errorLogFile, errorLogFileErr := os.OpenFile(errorLogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if errorLogFileErr != nil {
		log.Fatalf("[error] %s error creating logfile %+v", appID, errorLogFileErr)
		return
	}
	defer errorLogFile.Close()

	log.SetOutput(errorLogFile)
	log.SetPrefix(appID + " - ")

	// Load app settings
	cfg := general.DefaultConfiguration()
	settings, settingsErr := general.CreateAppSettings(cfg)
	if settingsErr != nil {
		log.Fatalf("[error] %s error establishing settings %+v", appID, settingsErr)
		return
	}

	settings.AppID = appID

	dbSession, dbSessionErr := sql.Open("sqlite3", settings.Config.DbFilePath)
	if dbSessionErr != nil {
		log.Fatalf("[error] %s error establishing settings %+v", appID, dbSessionErr)
		return
	}
	defer dbSession.Close()

	// SQLITE3 does not have foreign_keys turned on by default
	dbSession.Exec("PRAGMA foreign_keys = ON")

	// Load CLI
	app := &cli.App{
		Metadata: map[string]interface{}{
			"settings":  settings,
			"dbSession": dbSession,
		},
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
			suggestion.Command(),
			vote.Command(),
		},
	}

	cliErr := app.Run(args)
	if cliErr != nil {
		log.Printf("[error] %s",
			strings.Join(os.Args, " "))
		log.Fatalf("[error] %+v", cliErr)
	}
}
