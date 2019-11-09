package filters

import (
	"github.com/jabulba/disgord"
	"github.com/jabulba/disgord/std"
	"nodewarmanager/config"
)

// PrefixFilter is a disgord filter with the config.Bot.Prefix already configured
var PrefixFilter *std.MsgFilter

// HelpCommand is a disgord filter with the help prefix already configured
var HelpCommand *std.MsgFilter

// Load or reload the filters
func Load(client disgord.Session) {
	PrefixFilter, _ = std.NewMsgFilter(client)
	PrefixFilter.SetPrefix(config.Bot.Prefix)

	HelpCommand, _ = std.NewMsgFilter(client)
	HelpCommand.SetPrefix("help")
}
