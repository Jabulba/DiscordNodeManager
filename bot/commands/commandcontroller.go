package commandcontroller

import (
	"github.com/bwmarrin/discordgo"
	cmdchannel "nodewarmanager/bot/commands/channel"
	cmdhelp "nodewarmanager/bot/commands/help"
	cmdnodewar "nodewarmanager/bot/commands/nodewar"
	"nodewarmanager/config"
	"strings"
)

func MessageCreate(s *discordgo.Session, evt *discordgo.MessageCreate) {
	// Ignore bots
	if evt.Author.Bot {
		return
	}

	// Ignore messages that don't start with the prefix
	userMsg := strings.TrimSpace(evt.Content)
	if !strings.HasPrefix(userMsg, config.Bot.Prefix) {
		return
	}

	// Remove Prefix
	userMsg = strings.TrimPrefix(userMsg, config.Bot.Prefix)
	userMsg = strings.TrimSpace(userMsg)

	if strings.HasPrefix(userMsg, cmdchannel.Prefix) {
		// channel command
		userMsg = strings.TrimPrefix(userMsg, cmdchannel.Prefix)
		userMsg = strings.TrimSpace(userMsg)
		evt.Content = userMsg
		cmdchannel.MessageCreate(s, evt)
	} else if strings.HasPrefix(userMsg, cmdhelp.Prefix) {
		// help command
		cmdhelp.MessageCreate(s, evt)
	} else if strings.HasPrefix(userMsg, cmdnodewar.Prefix) {
		// nw command
		cmdnodewar.MessageCreate(s, evt)
	}
}
