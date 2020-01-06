package badgerdb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/pb"
	"log"
	"nodewarmanager/config"
	"sort"
	"strconv"
)

// BadgerDB implementation for IDB
type BadgerDB struct {
	db *badger.DB
}

func (bdb *BadgerDB) GetWarStatus(guildID string, date string) (map[string]int, error) {
	warStats := make(map[string]int)
	err := bdb.db.View(func(txn *badger.Txn) error {
		// Create a new stream to process KVs
		stream := bdb.db.NewStream()
		// Set the prefix to return only the specified guild channels
		stream.Prefix = []byte("G-" + guildID + "-War-" + date + "-")
		// Logger prefix to keep things clean!
		stream.LogPrefix = "Badger.Stream[" + string(stream.Prefix) + "]"

		// Stream will append all KVs to the chanList slice
		stream.Send = func(list *pb.KVList) error {
			for _, kv := range list.Kv {
				var ws []string
				err := json.Unmarshal(kv.Value, &ws)
				if err != nil {
					return err
				}

				for _, userId := range ws {
					if c, ok := warStats[userId]; ok {
						warStats[userId] = c + 1
					} else {
						warStats[userId] = 1
					}
				}
			}

			return nil
		}

		// Process the stream
		if err := stream.Orchestrate(context.Background()); err != nil {
			return err
		}

		// All good!
		return nil
	})

	return warStats, err
}

func (bdb *BadgerDB) SaveWarStatus(guildID string, date string, participants []string, tick int) error {
	return bdb.db.Update(func(txn *badger.Txn) error {
		// Create the key with prefix
		key := []byte("G-" + guildID + "-War-" + date + "-" + strconv.Itoa(tick))
		pJson, err := json.Marshal(participants)
		if err != nil {
			return err
		}

		return txn.Set(key, pJson)
	})
}

func (bdb *BadgerDB) ToggleMonitoredChannel(guildID string, channelID string) (bool, error) {
	// Default channel state is: NOT monitored
	monitored := false
	err := bdb.db.Update(func(txn *badger.Txn) error {
		// Create the key with prefix
		key := []byte("G-" + guildID + "-MonChan-" + channelID)
		item, err := txn.Get(key)
		if err != nil && !(errors.Is(err, badger.ErrKeyNotFound)) {
			// The error is not that the kv doesn't exist :(
			return err
		}

		if item != nil {
			// We have the kv stored so deleting it makes the channel no longer monitored
			return item.Value(func(val []byte) error {
				return txn.Delete(key)
			})
		} else {
			// The channel isn't being monitored, create a new entry to make it so!
			monitored = true
			// Save the value and return the outcome (err or nil)
			return txn.Set(key, []byte(channelID))
		}
	})

	return monitored, err
}

func (bdb *BadgerDB) GetMonitoredGuildChannelIDs(guildID string) ([]string, error) {
	var chanList []string
	err := bdb.db.View(func(txn *badger.Txn) error {
		// Create a new stream to process KVs
		stream := bdb.db.NewStream()
		// Set the prefix to return only the specified guild channels
		stream.Prefix = []byte("G-" + guildID + "-MonChan-")
		// Logger prefix to keep things clean!
		stream.LogPrefix = "Badger.Stream[" + string(stream.Prefix) + "]"

		// Stream will append all KVs to the chanList slice
		stream.Send = func(list *pb.KVList) error {
			for _, kv := range list.Kv {
				chanList = append(chanList, string(kv.Value))
			}

			return nil
		}

		// Process the stream
		if err := stream.Orchestrate(context.Background()); err != nil {
			return err
		}

		// All good!
		return nil
	})

	// Sort the slice before returning it allows for binary search and other neat stuff
	sort.Strings(chanList)

	return chanList, err
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
