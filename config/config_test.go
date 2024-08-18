package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewViperConfig(t *testing.T) {
	// Create a new Viper config
	config := NewViperConfig()

	// Assert that the config is not nil
	assert.NotNil(t, config)
}

func TestViperConfig_GetDatabasePath(t *testing.T) {
	// Create a new Viper config
	config := NewViperConfig()

	// Get the database path
	dbPath := config.GetDatabasePath()

	// Assert that the database path is the default value
	expectedPath := path.Join("data", "db")
	assert.Equal(t, expectedPath, dbPath)
}
