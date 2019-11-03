package config

import (
	"testing"
)

const assetionError = "Expected %s but got %s."

func TestConfigLoadingWithTestConfigFile(t *testing.T) {
	// Test using config file config_test.hocon
	*configPath = "../config_test.hocon"

	// Load the config
	err := Load()
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}

	// Check version value
	expectedVersion := "test-version"
	if Version != expectedVersion {
		t.Errorf(assetionError, expectedVersion, Version)
	}

	// Check bot.token value
	expectedBotToken := "test-token"
	if Bot.Token != expectedBotToken {
		t.Errorf(assetionError, expectedBotToken, Bot.Token)
	}

	// Check bot.prefix value
	expectedBotPrefix := "test-prefix"
	if Bot.Prefix != expectedBotPrefix {
		t.Errorf(assetionError, expectedVersion, Bot.Prefix)
	}

	// Check database.path value
	expectedDatabasePath := "test-database-path"
	if DB.Path != expectedDatabasePath {
		t.Errorf(assetionError, expectedDatabasePath, DB.Path)
	}
}

func TestConfigLoadingWithoutConfigFile(t *testing.T) {
	// Test using config file config_test.hocon
	*configPath = "../config_test.idontexist"

	// Load the config
	err := Load()
	if err != nil {
		t.Errorf("Unable to Read Configuration: " + err.Error())
	}

	// Check version value
	expectedVersion := "0.0.0-Adalwolf"
	if Version != expectedVersion {
		t.Errorf(assetionError, expectedVersion, Version)
	}

	// Check bot.token value
	expectedBotToken := "your-bot-token-here"
	if Bot.Token != expectedBotToken {
		t.Errorf(assetionError, expectedBotToken, Bot.Token)
	}

	// Check bot.prefix value
	expectedBotPrefix := "?"
	if Bot.Prefix != expectedBotPrefix {
		t.Errorf(assetionError, expectedBotPrefix, Bot.Prefix)
	}

	// Check database.path value
	expectedDatabasePath := "./databases/"
	if DB.Path != expectedDatabasePath {
		t.Errorf(assetionError, expectedDatabasePath, DB.Path)
	}
}
