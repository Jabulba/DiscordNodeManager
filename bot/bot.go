package bot

import (
	"github.com/bwmarrin/discordgo"
	commandcontroller "nodewarmanager/bot/commands"
	"nodewarmanager/config"
	"os"
	"os/signal"
	"syscall"
)

// Connect the bot with Discord and register all commands
func Connect() {
	dg, err := discordgo.New("Bot " + config.Bot.Token)
	if err != nil {
		panic(err)
	}

	// Register the command controller MessageCreate events callback.
	dg.AddHandler(commandcontroller.MessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	// Keep bot running until stopped
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session, ignoring errors
	_ = dg.Close()
}
