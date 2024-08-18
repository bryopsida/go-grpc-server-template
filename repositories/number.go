package repositories

import (
	"encoding/json"

	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/dgraph-io/badger"
)

type badgerNumberRepository struct {
	db *badger.DB
}

// NewBadgerNumberRepository creates a new badgerNumberRepository instance
func NewBadgerNumberRepository(db *badger.DB) interfaces.INumberRepository {
	return &badgerNumberRepository{db: db}
}

// Save saves a number
// - number: the number to save
// Returns an error if the save operation fails
func (r *badgerNumberRepository) Save(number interfaces.Number) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(number)
		if err != nil {
			return err
		}
		return txn.Set([]byte(number.ID), data)
	})
}

// FindByID finds a number by its ID
// - id: the ID of the number to find
// Returns the number if found, otherwise returns an error
func (r *badgerNumberRepository) FindByID(id string) (*interfaces.Number, error) {
	var number interfaces.Number
	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &number)
		})
	})
	if err != nil {
		return nil, err
	}
	return &number, nil
}

// DeleteByID deletes a number by its ID
// - id: the ID of the number to delete
// Returns an error if the delete operation fails
func (r *badgerNumberRepository) DeleteByID(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
