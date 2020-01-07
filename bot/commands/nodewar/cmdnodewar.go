package cmdnodewar

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"nodewarmanager/idb"
	"nodewarmanager/utils"
	"strconv"
	"sync"
	"time"
)

var Prefix = "nw"
var monitoredWars = make(map[string]chan bool)
var nwMux sync.Mutex

func MessageCreate(s *discordgo.Session, evt *discordgo.MessageCreate) {
	// Lock the map, allowing only one event to modify it at a time
	nwMux.Lock()

	// Unlock the map after processing the event
	defer nwMux.Unlock()

	// Check if there is a nodewar being monitored
	if c, ok := monitoredWars[evt.GuildID]; ok {
		// Signals the channel that terminates the nodewar monitoring
		c <- true
	} else {
		// Start monitoring a new nodewar
		go monitorWar(s, evt)
	}
}

// Start monitoring a nodewar
func monitorWar(s *discordgo.Session, evt *discordgo.MessageCreate) {
	nwMux.Lock()
	complete := make(chan bool)
	monitoredWars[evt.GuildID] = complete
	nwMux.Unlock()

	defer func() {
		nwMux.Lock()
		delete(monitoredWars, evt.GuildID)
		nwMux.Unlock()
	}()

	ticker := time.NewTicker(30 * time.Second)
	deadlineTimer := time.NewTimer(4 * time.Hour)

	guild, err := s.Guild(evt.GuildID)
	if err != nil {
		_, _ = s.ChannelMessageSend(evt.ChannelID, "I was unable to retrieve the guild info for the server so I can't monitor a war at this time :(")
		return
	}

	chanList, err := idb.DB.GetMonitoredGuildChannelIDs(guild.ID)
	if err != nil {
		_, _ = s.ChannelMessageSend(evt.ChannelID, "YIKES! I was unable to retrieve the list o monitored channels for the server so I can't monitor a war at this time :(\nTry using the 'channel' command!")
		return
	}

	if len(chanList) == 0 {
		_, _ = s.ChannelMessageSend(evt.ChannelID, "BOOO! No channel is being monitored :(")
		return
	}

	names, err := utils.GetChannelNames(chanList, s)
	if err != nil {
		names = "\nGGGAAAAAHHHHHH!!!!! I somehow failed to get the channel names... Use the channel command if you want to confirm the monitored channels."
	}
	_, _ = s.ChannelMessageSend(evt.ChannelID, "You war is now being monitored on the following channels: "+names)

	warDate := time.Now().Format("20060102")
	tick := 0

	for {
		select {
		case <-ticker.C:
			tick++

			guild, err = s.Guild(evt.GuildID)
			if err != nil {
				_, _ = s.ChannelMessageSend(evt.ChannelID, "I was unable to retrieve the guild info for the server so I can't monitor a war at this time :(")
				return
			}

			var participants []string
			for _, cid := range chanList {
				for _, vs := range guild.VoiceStates {
					// Skip deaf users, skip users not in monitored channels
					if vs.SelfDeaf || cid != vs.ChannelID {
						continue
					}

					participants = append(participants, vs.UserID)
				}
			}

			//_, _ = evt.Message.Reply(context.Background(), s, "tick "+strconv.Itoa(tick)+": "+strings.Join(participants, ", "))
			err = idb.DB.SaveWarStatus(guild.ID, warDate, participants, tick)
			if err != nil {
				log.Print(err)
			}
		case <-deadlineTimer.C:
			complete <- true
		case <-complete:
			sendWarSummary(guild, warDate, s, evt)
			ticker.Stop()
			deadlineTimer.Stop()
			return
		}
	}
}

func sendWarSummary(guild *discordgo.Guild, date string, s *discordgo.Session, evt *discordgo.MessageCreate) {
	participants, err := idb.DB.GetWarStatus(guild.ID, date)
	if err != nil {
		_, _ = s.ChannelMessageSend(evt.ChannelID, "Unable to retrieve war summary :(")
		return
	}
	_, _ = s.ChannelMessageSend(evt.ChannelID, "Please wait... Fetching participants nicknames with discord, ETA: "+strconv.Itoa(((len(participants)-1)/10)*10)+" seconds")

	summaryMessage := "War summary:\nId, Name, 1/2 minutes listening"
	for userId, t := range participants {
		userName := getUserName(guild.ID, userId, s)
		summaryMessage += "\n" + userId + ", " + userName + "," + strconv.Itoa(t)
	}
	_, _ = s.ChannelMessageSend(evt.ChannelID, summaryMessage)
}

func getUserName(guildId string, userId string, s *discordgo.Session) string {
	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		log.Printf("Failed to get user data. Using id instead of name.")
		return userId
	} else {
		if member.Nick != "" {
			return member.Nick
		} else if member.User.Username != "" {
			return member.User.Username
		} else {
			log.Printf("Failed to get user name or nick. Using id instead of name.")
			return userId
		}
	}
}
