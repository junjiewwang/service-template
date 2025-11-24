package strategies

import (
	"context"
	"os"
	"path/filepath"

	"github.com/junjiewwang/service-template/pkg/generator/filewriter"
	"github.com/junjiewwang/service-template/pkg/generator/filewriter/mergers"
)

// IncrementalStrategyID is the unique identifier for the incremental strategy
const IncrementalStrategyID = "incremental"

// IncrementalStrategy performs incremental updates using content mergers
type IncrementalStrategy struct {
	mergerID string
}

func init() {
	filewriter.DefaultStrategyRegistry.MustRegister(NewIncrementalStrategy())
}

// NewIncrementalStrategy creates a new incremental strategy with default merger
func NewIncrementalStrategy() *IncrementalStrategy {
	return &IncrementalStrategy{
		mergerID: mergers.MarkerMergerID,
	}
}

// WithMerger sets the merger ID to use
func (s *IncrementalStrategy) WithMerger(mergerID string) *IncrementalStrategy {
	return &IncrementalStrategy{
		mergerID: mergerID,
	}
}

// ID returns the strategy identifier
func (s *IncrementalStrategy) ID() string {
	return IncrementalStrategyID
}

// Description returns the strategy description
func (s *IncrementalStrategy) Description() string {
	return "Merge new content with existing content using marker blocks"
}

// Write writes content to the file using incremental merge
func (s *IncrementalStrategy) Write(ctx context.Context, path string, content []byte) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Get merger
	merger := mergers.DefaultMergerRegistry.MustGet(s.mergerID)

	// Check if file exists
	existingContent, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use merger to wrap content with markers
			// This ensures first-time generation also has markers
			mergedContent, err := merger.Merge(ctx, &mergers.MergeInput{
				ExistingContent: []byte{}, // Empty existing content
				NewContent:      content,
				FilePath:        path,
			})
			if err != nil {
				return err
			}
			return os.WriteFile(path, mergedContent, 0644)
		}
		return err
	}

	// File exists, use merger to merge content
	mergedContent, err := merger.Merge(ctx, &mergers.MergeInput{
		ExistingContent: existingContent,
		NewContent:      content,
		FilePath:        path,
	})
	if err != nil {
		return err
	}

	// Write merged content
	return os.WriteFile(path, mergedContent, 0644)
}
