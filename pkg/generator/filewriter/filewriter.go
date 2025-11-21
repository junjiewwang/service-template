package filewriter

import (
	"context"
)

// FileWriter provides a unified interface for file writing operations
type FileWriter struct {
	strategy WriteStrategy
}

// New creates a new FileWriter with the default strategy (overwrite)
func New() *FileWriter {
	// Get overwrite strategy as default
	// We use a lazy approach: try to get "overwrite", fallback to registry default
	strategy, exists := DefaultStrategyRegistry.Get("overwrite")
	if !exists {
		strategy = DefaultStrategyRegistry.GetDefault()
	}

	return &FileWriter{
		strategy: strategy,
	}
}

// WithStrategy sets the write strategy
func (w *FileWriter) WithStrategy(strategy WriteStrategy) *FileWriter {
	return &FileWriter{
		strategy: strategy,
	}
}

// Write writes content to the specified file path
func (w *FileWriter) Write(ctx context.Context, path string, content []byte) error {
	return w.strategy.Write(ctx, path, content)
}

// WriteString writes string content to the specified file path
func (w *FileWriter) WriteString(ctx context.Context, path string, content string) error {
	return w.Write(ctx, path, []byte(content))
}
