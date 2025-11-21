package strategies

import (
	"context"
	"os"
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/generator/filewriter"
)

// OverwriteStrategyID is the unique identifier for the overwrite strategy
const OverwriteStrategyID = "overwrite"

// OverwriteStrategy always overwrites existing files
type OverwriteStrategy struct{}

func init() {
	filewriter.DefaultStrategyRegistry.MustRegister(&OverwriteStrategy{})
}

// ID returns the strategy identifier
func (s *OverwriteStrategy) ID() string {
	return OverwriteStrategyID
}

// Description returns the strategy description
func (s *OverwriteStrategy) Description() string {
	return "Always overwrite existing files"
}

// Write writes content to the file, overwriting if it exists
func (s *OverwriteStrategy) Write(ctx context.Context, path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}
