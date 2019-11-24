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
)

func TestBadgerDB_Connect_Disconnect(t *testing.T) {
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
	defer os.RemoveAll(config.DB.BadgerDB.Path)

	// Create a new BadgerDB and connect to it
	bdb := BadgerDB{}
	if err = bdb.Connect(); err != nil {
		log.Fatal(err)
	}

	// Test connection to DB by adding a pair and reading it
	testKey := []byte("test key")
	testVal := []byte("test val")
	if err = bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(testKey, testVal)
	}); err != nil {
		t.Fatal("Failed to insert test pair in db: ", err)
	}

	if err = bdb.db.View(func(txn *badger.Txn) error {
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
