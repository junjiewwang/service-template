package strategies

import (
	"context"
	"os"
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/generator/filewriter"
)

// SkipStrategyID is the unique identifier for the skip strategy
const SkipStrategyID = "skip"

// SkipStrategy skips writing if the file already exists
type SkipStrategy struct{}

func init() {
	filewriter.DefaultStrategyRegistry.MustRegister(&SkipStrategy{})
}

// ID returns the strategy identifier
func (s *SkipStrategy) ID() string {
	return SkipStrategyID
}

// Description returns the strategy description
func (s *SkipStrategy) Description() string {
	return "Skip writing if file already exists"
}

// Write writes content to the file only if it doesn't exist
func (s *SkipStrategy) Write(ctx context.Context, path string, content []byte) error {
	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		// File exists, skip writing
		return nil
	}

	// File doesn't exist, create it
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}
