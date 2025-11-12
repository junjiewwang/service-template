package core

import "errors"

// Common errors
var (
	// ErrNilContext indicates that the generator context is nil
	ErrNilContext = errors.New("generator context is nil")

	// ErrNilConfig indicates that the service config is nil
	ErrNilConfig = errors.New("service config is nil")

	// ErrGeneratorNotFound indicates that the requested generator type is not registered
	ErrGeneratorNotFound = errors.New("generator not found")

	// ErrInvalidOptions indicates that the provided options are invalid
	ErrInvalidOptions = errors.New("invalid generator options")

	// ErrTemplateRender indicates a template rendering error
	ErrTemplateRender = errors.New("template rendering failed")

	// ErrValidation indicates a validation error
	ErrValidation = errors.New("validation failed")
)
