package main

import (
	"log"
	"nodewarmanager/bot"
	"nodewarmanager/config"
	"nodewarmanager/idb"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := idb.Init()
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Connect(); err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	bot.Connect(db)
}
