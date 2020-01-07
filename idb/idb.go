package idb

import (
	"fmt"
	"nodewarmanager/config"
	"nodewarmanager/idb/badgerdb"
)

var (
	DB IDatabase
)

// IDatabase is the interface used to connect and manage the database defined in the configuration file
type IDatabase interface {
	// Connect to the database
	Connect() error
	// Disconnect from the database
	Disconnect()
	// GetMonitoredGuildChannelIDs will return a sorted slice with all monitored channels ids
	GetMonitoredGuildChannelIDs(guildID string) ([]string, error)
	// ToggleMonitoredChannel will either add or remove the channel from the monitored channels list
	ToggleMonitoredChannel(guildID string, channelID string) (bool, error)
	// SaveWarStatus will add to the database with the 'date' prefix all the participants for later use
	SaveWarStatus(guildID string, date string, participants []string, tick int) error
	// GetWarStatus will retrieve the participation status of all members for the given day
	GetWarStatus(guildID string, date string) (map[string]int, error)
}

// Init will initialize the database defined in the configuration
func Init() error {
	if len(config.DB.BadgerDB.Path) != 0 {
		// Using Badger DB implementation
		DB = &badgerdb.BadgerDB{}
		return nil
	}

	return fmt.Errorf("no database has been configured")
}
