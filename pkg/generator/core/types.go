package core

import (
	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// Generator defines the interface that all generators must implement
type Generator interface {
	// Generate generates content from template
	Generate() (string, error)

	// GetName returns the generator name
	GetName() string

	// Validate validates the generator configuration
	Validate() error

	// GetContext returns the generator context
	GetContext() *context.GeneratorContext
}

// GeneratorCreator is a function type that creates a generator
type GeneratorCreator func(ctx *context.GeneratorContext, options ...interface{}) (Generator, error)
