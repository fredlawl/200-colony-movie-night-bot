package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// TODO: Use https://github.com/boltdb/bolt for embeded DB
// TODO: Use github.com/urfave/cli for cli configuration

func main() {
	log.Println("before create")
	cfg := CreateAppConfig()
	log.Println("after create")

	db, err := bolt.Open("cli.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("test")
	log.Println(cfg.appState)

	defer db.Close()
}
