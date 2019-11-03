package main

import (
	"nodewarmanager/config"
)

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}
}
