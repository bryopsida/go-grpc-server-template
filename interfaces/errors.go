package interfaces

import "errors"

const (
	// ErrMsgNotFound is the error message for when a resource is not found
	ErrMsgNotFound = "not found"
	// ErrMsgSaveFailed is the error message for when a save operation fails
	ErrMsgSaveFailed = "save failed"
)

var (
	// ErrNotFound is an error for when a resource is not found
	ErrNotFound = errors.New(ErrMsgNotFound)
	// ErrSaveFailed is an error for when a save operation fails
	ErrSaveFailed = errors.New(ErrMsgSaveFailed)
)
