package cmdchannel

import (
	"github.com/bwmarrin/discordgo"
	"nodewarmanager/config"
	"nodewarmanager/idb"
	"nodewarmanager/utils"
	"sort"
	"strconv"
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
var Prefix = "channel"

const errMsg = "Aww, snap! Something terrible has happened... Your request has failed :("

func MessageCreate(s *discordgo.Session, evt *discordgo.MessageCreate) {
	var chatReply string
	if len(evt.Content) == 0 {
		// No channel was specified, show the help message and the summary of channels
		chatReply = displaySummary(s, evt.GuildID)
	} else {
		// Toggle specified channel
		chatReply = toggleChannelMonitoring(evt.Content, s, evt.GuildID)
	}

	_, _ = s.ChannelMessageSend(evt.ChannelID, chatReply)
}

func toggleChannelMonitoring(userMsg string, s *discordgo.Session, guildID string) string {
	var chatMsg string
	chanNum, err := strconv.Atoi(userMsg)
	if err != nil {
		// Specified channel number is not a number...
		chatMsg = "Your input '" + userMsg + "' is not a number... You need to specify the number of the channel here!"
	} else {
		// We got a number! Lets add or remove it from the list
		channels, err := s.GuildChannels(guildID)
		if err != nil {
			return errMsg
		}

		for _, c := range channels {
			// Skip non voice channels
			if c.Type != discordgo.ChannelTypeGuildVoice {
				continue
			}

			// Check for the desired channel
			if c.Position == chanNum {
				// Toggle channel
				monitored, err := idb.DB.ToggleMonitoredChannel(guildID, c.ID)
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

func displaySummary(s *discordgo.Session, guildID string) string {
	chanList, err := idb.DB.GetMonitoredGuildChannelIDs(guildID)
	if err != nil {
		return errMsg
	}

	var chatMsg string
	if len(chanList) == 0 {
		chatMsg = helpText + "BOOO! No channel is being monitored :("
	} else {
		chatMsg = helpText
		channelNames, err := utils.GetChannelNames(chanList, s)
		if err != nil {
			return errMsg
		}

		chatMsg += channelNames
	}

	// Now we will list all available channels and their numbers so the user can add more channels to the monitored list!
	channels, err := s.GuildChannels(guildID)
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
		if c.Type != discordgo.ChannelTypeGuildVoice {
			continue
		}

		// Ignore all monitored channels
		for _, cID := range chanList {
			if c.ID == cID {
				continue ChannelsFor
			}
		}

		chatMsg += strconv.Itoa(c.Position) + ") " + c.Name + "\n"
	}

	return chatMsg
}
