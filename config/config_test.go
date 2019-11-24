package config

import (
	"flag"
	"nodewarmanager/tests"
	"os"
	"testing"
)

func TestConfigLoadingWithTestConfigFile(t *testing.T) {
	// Test using config file config_test.hocon
	err := flag.Set("config", "../tests/configs/config_test.hocon")
	if err != nil {
		t.Fatal("Unable to change configuration file path: " + err.Error())
	}

	// Load the config
	err = Load()
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}

	// Check version value
	expectedVersion := "test-version"
	if Version != expectedVersion {
		t.Errorf(tests.AssertionError, expectedVersion, Version)
	}

	// Check bot.token value
	expectedBotToken := "test-token"
	if Bot.Token != expectedBotToken {
		t.Errorf(tests.AssertionError, expectedBotToken, Bot.Token)
	}

	// Check bot.prefix value
	expectedBotPrefix := "test-prefix"
	if Bot.Prefix != expectedBotPrefix {
		t.Errorf(tests.AssertionError, expectedBotPrefix, Bot.Prefix)
	}

	// Check bot.debug value
	if Bot.Debug {
		t.Errorf(tests.AssertionError, false, Bot.Debug)
	}

	// Check database.path value
	expectedDatabasePath := "test-database-path"
	if DB.BadgerDB.Path != expectedDatabasePath {
		t.Errorf(tests.AssertionError, expectedDatabasePath, DB.BadgerDB.Path)
	}
}

func TestConfigLoadingWithoutConfigFile(t *testing.T) {
	// Test using config file config_test.fail.hocon
	err := flag.Set("config", "../tests/configs/config_test.fail.hocon")
	if err != nil {
		t.Fatal("Unable to change configuration file path: " + err.Error())
	}

	// Cleanup after test
	defer os.Remove(*configPath)

	// Load the config
	err = Load()
	if err == nil {
		t.Error("Expected an error but received none.")
	}

	info, err := os.Stat(*configPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected %s to exist but it doesnt.", *configPath)
	}

	if info.IsDir() {
		t.Errorf("Expected %s to be a file but its not.", *configPath)
	}
}
