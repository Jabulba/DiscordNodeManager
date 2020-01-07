package badgerdb

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"nodewarmanager/config"
	"nodewarmanager/tests"
	"os"
	"testing"
	"time"
)

func TestBadgerDB_Connect_Disconnect(t *testing.T) {
	prepareTest(t)
	defer os.RemoveAll(config.DB.BadgerDB.Path)

	// Create a new BadgerDB and connect to it
	bdb := BadgerDB{}
	if err := bdb.Connect(); err != nil {
		log.Fatal(err)
	}

	// Test connection to DB by adding a pair and reading it
	testKey := []byte("test key")
	testVal := []byte("test val")
	if err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(testKey, testVal)
	}); err != nil {
		t.Fatal("Failed to insert test pair in db: ", err)
	}

	if err := bdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(testKey)
		if err != nil {
			return fmt.Errorf("failed to retrieve test pair from db: %w", err)
		}

		if err = item.Value(func(val []byte) error {
			if bytes.Compare(val, testVal) != 0 {
				t.Errorf(tests.AssertionError, testVal, item)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("failed to retrieve test value from db: %w", err)
		}

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// Disconnect from DB
	bdb.Disconnect()
}

func TestBadgerDB_GetMonitoredGuildChannelIDs_ToggleMonitoredChannel(t *testing.T) {
	prepareTest(t)
	defer os.RemoveAll(config.DB.BadgerDB.Path)

	bdb := BadgerDB{}
	if err := bdb.Connect(); err != nil {
		log.Fatal(err)
	}

	guildId := "12345"

	// Assert no channels return from monitored channels list
	channels, err := bdb.GetMonitoredGuildChannelIDs(guildId)
	if err != nil {
		log.Fatal(err)
	}

	l := len(channels)
	if l != 0 {
		t.Errorf(tests.AssertionError, 0, l)
	}

	// Monitor a channel
	channel1ID := "987654"

	b, err := bdb.ToggleMonitoredChannel(guildId, channel1ID)
	if err != nil {
		log.Fatal(err)
	}
	if !b {
		t.Errorf(tests.AssertionError, true, b)
	}

	// Make sure 1 channel is returned
	channels, err = bdb.GetMonitoredGuildChannelIDs(guildId)
	if err != nil {
		log.Fatal(err)
	}

	l = len(channels)
	if l != 1 {
		t.Errorf(tests.AssertionError, 1, l)
	}

	// Stop monitoring a channel
	b, err = bdb.ToggleMonitoredChannel(guildId, channel1ID)
	if err != nil {
		log.Fatal(err)
	}
	if b {
		t.Errorf(tests.AssertionError, false, b)
	}

	// Assert no channels return from monitored channels list
	channels, err = bdb.GetMonitoredGuildChannelIDs(guildId)
	if err != nil {
		log.Fatal(err)
	}

	l = len(channels)
	if l != 0 {
		t.Errorf(tests.AssertionError, 0, l)
	}

	// Disconnect from DB
	bdb.Disconnect()
}

func TestBadgerDB_SaveWarStatus_GetWarStatus(t *testing.T) {
	prepareTest(t)
	defer os.RemoveAll(config.DB.BadgerDB.Path)

	bdb := BadgerDB{}
	if err := bdb.Connect(); err != nil {
		log.Fatal(err)
	}

	guildID := "12345"
	warDate := time.Date(2020, 01, 01, 21, 00, 00, 0, time.Local).Format("20060102")
	participants := []string{"001", "002"}
	err := bdb.SaveWarStatus(guildID, warDate, participants, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = bdb.SaveWarStatus(guildID, warDate, participants, 1)
	if err != nil {
		log.Fatal(err)
	}

	err = bdb.SaveWarStatus(guildID, warDate, []string{"001"}, 2)
	if err != nil {
		log.Fatal(err)
	}

	status, err := bdb.GetWarStatus(guildID, warDate)
	if err != nil {
		log.Fatal(err)
	}

	u001 := status["001"]
	if u001 != 3 {
		t.Errorf(tests.AssertionError, 3, u001)
	}

	u002 := status["002"]
	if u002 != 2 {
		t.Errorf(tests.AssertionError, 2, u001)
	}

	// Disconnect from DB
	bdb.Disconnect()
}

func prepareTest(t *testing.T) {
	// Test using config file badgerdb.hocon
	err := flag.Set("config", "../../tests/configs/badgerdb.hocon")
	if err != nil {
		t.Fatal("Unable to change configuration file path: " + err.Error())
	}
	if err = config.Load(); err != nil {
		t.Fatal("Failed to load configuration: ", err)
	}
	// Check database.path value
	expectedDatabasePath := "../../tests/.badgerdb_test/"
	if config.DB.BadgerDB.Path != expectedDatabasePath {
		t.Fatalf(tests.AssertionError, expectedDatabasePath, config.DB.BadgerDB.Path)
	}

	if _, err := os.Stat(config.DB.BadgerDB.Path); err == nil {
		t.Fatal("Directory used for database tests already exists! Dir: " + config.DB.BadgerDB.Path)
	}
}
