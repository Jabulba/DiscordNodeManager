package main

import (
	"log"
	"nodewarmanager/config"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
}
