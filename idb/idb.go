package idb

import (
	"fmt"
	"nodewarmanager/config"
	"nodewarmanager/idb/badgerdb"
)

// IDatabase is the interface used to connect and manage the database defined in t he configuration file
type IDatabase interface {
	// Connect to the database
	Connect() error
	// Disconnect from the database
	Disconnect()
	// GetMonitoredGuildChannelIDs will return a sorted slice with all monitored channels ids
	GetMonitoredGuildChannelIDs(guildID string) ([]string, error)
	// ToggleMonitoredChannel will either add or remove the channel from the monitored channels list
	ToggleMonitoredChannel(guildID string, channelID string) (bool, error)
}

// Init will initialize the database defined in the configuration
func Init() (IDatabase, error) {
	if len(config.DB.BadgerDB.Path) != 0 {
		// Using Badger DB implementation
		return &badgerdb.BadgerDB{}, nil
	}

	return nil, fmt.Errorf("no database has been configured")
}
