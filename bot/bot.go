package bot

import (
	"github.com/jabulba/disgord"
	"nodewarmanager/bot/commands/help"
	"nodewarmanager/bot/filters"
	"nodewarmanager/config"
)

// Connect the bot with Discord and register all commands
func Connect() {
	client := disgord.New(disgord.Config{
		BotToken: config.Bot.Token,
		Logger:   disgord.DefaultLogger(config.Bot.Debug), // debug=false
	})

	defer client.StayConnectedUntilInterrupted()
	filters.Load(client)
	help.Register(client)
}
