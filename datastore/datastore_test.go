package datastore

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabase(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()
	dbPath := path.Join(tempDir, "testdb")

	// Create a mock config
	mockConfig := new(MockConfig)
	mockConfig.On("GetDatabasePath").Return(dbPath)

	// Call GetDatabase
	db, err := GetDatabase(mockConfig)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Close the database
	err = db.Close()
	assert.NoError(t, err)

	// Verify the mock expectations
	mockConfig.AssertExpectations(t)

	// Clean up
	os.RemoveAll(tempDir)
}
