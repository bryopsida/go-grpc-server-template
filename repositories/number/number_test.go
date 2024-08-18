package number

import (
	"encoding/json"
	"testing"

	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBadgerNumberRepository(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()

	// Open a Badger database
	opts := badger.DefaultOptions(tempDir)
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository
	repo := NewBadgerNumberRepository(db)
	assert.NotNil(t, repo)
}

func TestBadgerNumberRepository_Save(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()

	// Open a Badger database
	opts := badger.DefaultOptions(tempDir)
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository
	repo := NewBadgerNumberRepository(db)

	// Create a sample number
	number := interfaces.Number{ID: "1", Number: 42}

	// Save the number
	err = repo.Save(number)
	assert.NoError(t, err)

	// Verify the number was saved
	var savedNumber interfaces.Number
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(number.ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &savedNumber)
		})
	})
	assert.NoError(t, err)
	assert.Equal(t, number, savedNumber)
}

func TestBadgerNumberRepository_FindByID(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()

	// Open a Badger database
	opts := badger.DefaultOptions(tempDir)
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository
	repo := NewBadgerNumberRepository(db)

	// Create a sample number
	number := interfaces.Number{ID: "1", Number: 42}

	// Save the number
	err = repo.Save(number)
	require.NoError(t, err)

	// Find the number by ID
	foundNumber, err := repo.FindByID(number.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundNumber)
	assert.Equal(t, number, *foundNumber)
}

func TestBadgerNumberRepository_DeleteByID(t *testing.T) {
	// Create a temporary directory for the database
	tempDir := t.TempDir()

	// Open a Badger database
	opts := badger.DefaultOptions(tempDir)
	db, err := badger.Open(opts)
	require.NoError(t, err)
	defer db.Close()

	// Create a new repository
	repo := NewBadgerNumberRepository(db)

	// Create a sample number
	number := interfaces.Number{ID: "1", Number: 42}

	// Save the number
	err = repo.Save(number)
	require.NoError(t, err)

	// Verify the number was saved
	var savedNumber interfaces.Number
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(number.ID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &savedNumber)
		})
	})
	require.NoError(t, err)
	assert.Equal(t, number, savedNumber)

	// Delete the number by ID
	err = repo.DeleteByID(number.ID)
	assert.NoError(t, err)

	// Verify the number was deleted
	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(number.ID))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		return err
	})
	assert.NoError(t, err)
}
