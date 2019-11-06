package help

import (
	"github.com/jabulba/disgord"
	"nodewarmanager/bot/filters"
	"nodewarmanager/config"
)

// Register the help command with disgord
func Register(client disgord.Session) {
	helpText := `Some catchy phrase goes here! But none has been thought of yet...

Commands:
	` + config.Bot.Prefix + ` channel [#]
		-	Add or remove a channel from the monitored channels list.
			If no channel is specified, a summary of channels will be displayed.
	` + config.Bot.Prefix + ` start
		-	Star monitoring the node war.
	` + config.Bot.Prefix + ` stop
		-	Stop monitoring the node war and print a summary of attendance.
	` + config.Bot.Prefix + ` payout
		-	Finish a payout period and start a new one.
	` + config.Bot.Prefix + ` histogram [#]
		-	Displays the history of attendance for the last # payout periods.
			If you don't specify #, only the last payout period will be shown'

When monitoring a node war, the bot will keep an eye on who enters and leaves the channels added to the monitored channels list. When you end the node war, all attendees will be saved and a summary displayed.
The payout period is composed by a start date and an end date. All node wars that happen during the payout period are combined in the payout summary.`

	client.On(disgord.EvtMessageCreate,
		filters.PrefixFilter.HasPrefix,
		filters.PrefixFilter.StripPrefix,
		filters.HelpCommand.HasPrefix,
		func(s disgord.Session, evt *disgord.MessageCreate) {
			_, _ = evt.Message.Reply(s, helpText)
		})
}
