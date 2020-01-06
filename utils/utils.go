package utils

import (
	"github.com/bwmarrin/discordgo"
	"sort"
	"strconv"
)

// GetChannelNames will return the list of channels separated by a new line with the following format: '#) NAME'
func GetChannelNames(chanList []string, s *discordgo.Session) (string, error) {
	strChanList := ""
	var ordChanList []*discordgo.Channel

	for _, cid := range chanList {
		c, err := s.Channel(cid)
		if err != nil {
			return "", err
		}

		ordChanList = append(ordChanList, c)
	}

	sort.SliceStable(ordChanList, func(i, j int) bool {
		return ordChanList[i].Position < ordChanList[j].Position
	})

	for _, c := range ordChanList {
		strChanList += "\n" + strconv.Itoa(c.Position) + ") " + c.Name
	}

	return strChanList, nil
}
