package config

import (
	"errors"
	"flag"
	"github.com/go-akka/configuration"
	"io/ioutil"
	"log"
	"os"
)

var (
	configPath = flag.String("config", "./config.hocon", "Configuration File Location. Ex: ./config.hocon")
	conf       *configuration.Config
	// Version as defined by the configuration file
	Version string
	// Bot specific configurations
	Bot bot
	// Database specific configurations
	DB db
)

type bot struct {
	// Token used to authenticate with Discord API
	Token string
	// Prefix recognized by the bot in text channels
	Prefix string
	// Debug true/false to enable/disable debug logging
	Debug bool
}

type db struct {
	// BadgerDB config
	BadgerDB bdb
}

type bdb struct {
	// Path used to save all Database files
	Path string
}

// Load or reload the configuration file
func Load() error {
	// Read the configuration file
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, save a default config file and exit
			err = ioutil.WriteFile(*configPath, []byte(defaultConfig), 0660)
			if err != nil {
				log.Printf("Failed to create configuration file: %s", err.Error())
				return err
			}

			log.Printf("A new config file has been created at '%s' for you to customize before launching the bot again.", *configPath)
			return errors.New("no config file found")
		} else {
			log.Printf("Failed to load configuration file: %s", err.Error())
			return err
		}
	} else {
		conf = configuration.ParseString(string(file))
	}

	// Validate config file structure
	if !conf.HasPath("bot") {
		return errors.New("config file is missing the bot specific values")
	}

	if !conf.HasPath("bot.token") {
		return errors.New("config file is missing the bot.token value")
	}

	if !conf.HasPath("bot.prefix") {
		return errors.New("config file is missing the bot.prefix value")
	}

	if !conf.HasPath("database.badger.path") {
		return errors.New("config file is missing the database.badger.path value")
	}

	if !conf.HasPath("version") {
		return errors.New("config file is missing the version value")
	}

	// Read config file values
	Version = conf.GetString("version")

	Bot = bot{
		Token:  conf.GetString("bot.token"),
		Prefix: conf.GetString("bot.prefix"),
		Debug:  conf.GetBoolean("bot.debug"),
	}

	DB = db{
		BadgerDB: bdb{
			Path: conf.GetString("database.badger.path"),
		},
	}

	// Validate empty token
	if len(Bot.Token) == 0 {
		log.Printf("You need to edit the config file '%s' and add your bot token before running the bot.", *configPath)
		return errors.New("config file has an empty value for bot.token")
	}

	return nil
}

var defaultConfig = `####################################
# Nodewar Manager Configuration    #
####################################
version: "0.0.0-Adalwolf"

# Bot specific configurations
bot {
  # Your Discord bot token
  token: "",

  # The prefix used to identify a command in chat
  prefix: "?",
}

# Database specific configurations
database {
  # Configurations for BadgerDB implementation.
  badger {
    # The folder to store files related to the database
    path: "./.badgerdb/",
  }
}`
