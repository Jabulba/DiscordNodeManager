package cmdchannel

import (
	"github.com/jabulba/disgord"
	"github.com/jabulba/disgord/std"
	"nodewarmanager/bot/chatfilters"
	"nodewarmanager/config"
	"nodewarmanager/idb"
	"sort"
	"strconv"
	"strings"
)

var helpText = `Use this command to define your Node War channels. All channels registered with this command will be monitored during a Node War!

Usage:
	` + config.Bot.Prefix + ` channel [#]
		-	Add or remove a channel from the monitored channels list.
			If no channel is specified, a summary of channels will be displayed.

When monitoring a node war, the bot will keep an eye on who enters and leaves the channels added to the monitored channels list. When you end the node war, all attendees will be saved and a summary displayed.
The payout period is composed by a start date and an end date. All node wars that happen during the payout period are combined in the payout summary.

Monitored Channels:
`

const errMsg = "Aww, snap! Something terrible has happened... Your request has failed :("

// Register the help command with disgord
func Register(client disgord.Session, db idb.IDatabase) {
	client.On(disgord.EvtMessageCreate,
		chatfilters.PrefixFilter.NotByBot,
		chatfilters.PrefixFilter.HasPrefix,
		std.CopyMsgEvt,
		chatfilters.PrefixFilter.StripPrefix,
		chatfilters.ChannelCommand.HasPrefix,
		chatfilters.ChannelCommand.StripPrefix,
		func(s disgord.Session, evt *disgord.MessageCreate) {
			var chatMsg string
			userMsg := strings.TrimSpace(evt.Message.Content)
			if len(userMsg) == 0 {
				// No channel was specified, show the help message and the summary of channels
				chatMsg = displaySummary(s, evt, db)
			} else {
				// Toggle specified channel
				chatMsg = toggleChannelMonitoring(userMsg, s, evt, db)
			}

			_, _ = evt.Message.Reply(s, chatMsg)
		})
}

func toggleChannelMonitoring(userMsg string, s disgord.Session, evt *disgord.MessageCreate, db idb.IDatabase) string {
	var chatMsg string
	chanNum, err := strconv.Atoi(userMsg)
	if err != nil {
		// Specified channel number is not a number...
		chatMsg = "Your input '" + userMsg + "' is not a number... You need to specify the number of the channel here!"
	} else {
		// We got a number! Lets add or remove it from the list
		channels, err := s.GetGuildChannels(evt.Message.GuildID)
		if err != nil {
			return errMsg
		}

		for _, c := range channels {
			// Skip non voice channels
			if c.Type != disgord.ChannelTypeGuildVoice {
				continue
			}

			// Check for the desired channel
			if c.Position == chanNum {
				// Toggle channel
				monitored, err := db.ToggleMonitoredChannel(evt.Message.GuildID.String(), c.ID.String())
				if err != nil {
					return errMsg
				}

				return "The channel '" + c.Name + "' monitoring status has been set to: '" + strconv.FormatBool(monitored) + "'"
			}
		}

		// No voice channel with the desired number was found...
		return "OH NOES! I was unable to find the voice channel '" + userMsg + "' 😱"
	}

	return chatMsg
}

func displaySummary(s disgord.Session, evt *disgord.MessageCreate, db idb.IDatabase) string {
	chanList, err := db.GetMonitoredGuildChannelIDs(evt.Message.GuildID.String())
	if err != nil {
		return errMsg
	}

	var chatMsg string
	var ordChanList []*disgord.Channel
	if len(chanList) == 0 {
		chatMsg = helpText + "BOOO! No channel is being monitored :("
	} else {
		chatMsg = helpText
		for _, cid := range chanList {
			snowflake, err := disgord.GetSnowflake(cid)
			if err != nil {
				return errMsg
			}

			c, err := s.GetChannel(snowflake)
			if err != nil {
				return errMsg
			}

			ordChanList = append(ordChanList, c)
		}

		sort.SliceStable(ordChanList, func(i, j int) bool {
			return ordChanList[i].Position < ordChanList[j].Position
		})

		for _, c := range ordChanList {
			chatMsg += "\n" + strconv.Itoa(c.Position) + ") " + c.Name
		}
	}

	// Now we will list all available channels and their numbers so the user can add more channels to the monitored list!
	channels, err := s.GetGuildChannels(evt.Message.GuildID)
	if err != nil {
		// Something went wrong but we already have a lot of information here!
		return chatMsg
	}

	chatMsg += "\n\nAvailable Channels:\n"

	sort.SliceStable(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})

ChannelsFor:
	for _, c := range channels {
		if c.Type != disgord.ChannelTypeGuildVoice {
			continue
		}

		// Ignore all monitored channels
		for _, cID := range chanList {
			if c.ID.String() == cID {
				continue ChannelsFor
			}
		}

		chatMsg += strconv.Itoa(c.Position) + ") " + c.Name + "\n"
	}

	return chatMsg
}
