package badgerdb

import (
	"github.com/dgraph-io/badger"
	"log"
	"nodewarmanager/config"
)

// BadgerDB implementation for IDB
type BadgerDB struct {
	db *badger.DB
}

// Connect to the database
func (bdb *BadgerDB) Connect() error {
	options := badger.DefaultOptions(config.DB.BadgerDB.Path)
	options.Truncate = true

	var err error
	bdb.db, err = badger.Open(options)

	return err
}

// Disconnect from the database
func (bdb *BadgerDB) Disconnect() {
	err := bdb.db.Close()
	if err != nil {
		log.Print("Error disconnecting from BadgerDB: ", err)
	} else {
		log.Print("Disconnected from BadgerDB...")
	}
}
