package context

import "errors"

// Common errors for context package
var (
	// ErrNilConfig indicates that the service config is nil
	ErrNilConfig = errors.New("service config is nil")

	// ErrNilPaths indicates that the paths are nil
	ErrNilPaths = errors.New("paths are nil")

	// ErrNilVariablePool indicates that the variable pool is nil
	ErrNilVariablePool = errors.New("variable pool is nil")
)
