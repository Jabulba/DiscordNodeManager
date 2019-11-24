package idb

import (
	"fmt"
	"nodewarmanager/config"
	"nodewarmanager/idb/badgerdb"
)

// IDatabase is the interface used to connect and manage the database defined in t he configuration file
type IDatabase interface {
	Connect() error
	Disconnect()
}

// Init will initialize the database defined in the configuration
func Init() (IDatabase, error) {
	if len(config.DB.BadgerDB.Path) != 0 {
		// Using Badger DB implementation
		return &badgerdb.BadgerDB{}, nil
	}

	return nil, fmt.Errorf("no database has been configured")
}
