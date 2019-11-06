package main

import (
	"log"
	"nodewarmanager/bot"
	"nodewarmanager/config"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	defer bot.Connect()
}
