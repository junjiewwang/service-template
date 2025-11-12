package context

import "errors"

// Common errors for context package
var (
	// ErrNilConfig indicates that the service config is nil
	ErrNilConfig = errors.New("service config is nil")

	// ErrNilVariables indicates that the variables are nil
	ErrNilVariables = errors.New("variables are nil")

	// ErrNilPaths indicates that the paths are nil
	ErrNilPaths = errors.New("paths are nil")
)
