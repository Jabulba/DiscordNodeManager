package config

import (
	"flag"
	"github.com/go-akka/configuration"
	"io/ioutil"
	"log"
	"os"
)

var (
	configPath = flag.String("config", "./config.hocon", "Configuration File Location. Ex: ./config.hocon")
	conf       *configuration.Config
	Version    string
	Bot        bot
	DB         db
)

type bot struct {
	Token  string
	Prefix string
}

type db struct {
	Path string
}

func Load() error {
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			conf = configuration.ParseString(defaultConfig)
		} else {
			log.Printf("Failed to load configuration file: %s", err.Error())
			return err
		}
	} else {
		conf = configuration.ParseString(string(file))
	}

	Version = conf.GetString("version")

	Bot = bot{
		Token:  conf.GetString("bot.token"),
		Prefix: conf.GetString("bot.prefix"),
	}

	DB = db{
		Path: conf.GetString("database.path"),
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
  # The folder to store files related to the database
  path: "./databases/",
}`
