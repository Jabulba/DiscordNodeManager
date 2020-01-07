package main

import (
	"log"
	"nodewarmanager/bot"
	"nodewarmanager/config"
	"nodewarmanager/idb"
)

func main() {
	var err error
	if err = config.Load(); err != nil {
		log.Fatal(err)
	}

	if err = idb.Init(); err != nil {
		log.Fatal(err)
	}

	if err = idb.DB.Connect(); err != nil {
		log.Fatal(err)
	}

	defer idb.DB.Disconnect()

	bot.Connect()
}
