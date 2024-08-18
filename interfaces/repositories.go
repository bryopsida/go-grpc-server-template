package interfaces

// Number is a struct to represent a number
type Number struct {
	// ID is the unique identifier of the number
	ID string
	// Number is the value of the number
	Number uint64
}

// INumberRepository is an interface for number repositories
type INumberRepository interface {
	// Save saves a number
	// - number: the number to save
	// Returns an error if the save operation fails
	Save(number Number) error
	// FindByID finds a number by its ID
	// - id: the ID of the number to find
	// Returns the number if found, otherwise returns an error
	FindByID(id string) (*Number, error)
	// DeleteByID deletes a number by its ID
	// - id: the ID of the number to delete
	// Returns an error if the delete operation fails
	DeleteByID(id string) error
}
