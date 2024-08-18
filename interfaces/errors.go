package interfaces

import "errors"

const (
	ErrMsgNotFound   = "not found"
	ErrMsgSaveFailed = "save failed"
)

var (
	ErrNotFound   = errors.New(ErrMsgNotFound)
	ErrSaveFailed = errors.New(ErrMsgSaveFailed)
)
