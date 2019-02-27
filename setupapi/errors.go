package setupapi

import (
	"errors"
	"syscall"
)

var (
	// ErrEmptyBuffer is returned when a nil or zero-sized buffer is provided
	// to a system call.
	ErrEmptyBuffer = errors.New("nil or empty buffer provided")

	// ErrInvalidRegistry is returned when an unexpected registry value type
	// is encountered.
	//ErrInvalidRegistry = errors.New("invalid registry type")

	// ErrInvalidData is returned when a property isn't present or isn't valid.
	ErrInvalidData = syscall.Errno(13)
)
