package examples

import (
	"context"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/filewriter"
	"github.com/junjiewwang/service-template/pkg/generator/filewriter/strategies"
)

// Example 1: Simple overwrite (default behavior)
func ExampleSimpleOverwrite() error {
	ctx := context.Background()
	writer := filewriter.New()

	content := "# My Makefile\n.PHONY: build\nbuild:\n\tgo build\n"
	err := writer.WriteString(ctx, "/tmp/Makefile", content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("File written successfully")
	return nil
}

// Example 2: Skip if file exists
func ExampleSkipExisting() error {
	ctx := context.Background()

	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.SkipStrategyID))

	content := "# Configuration file\nkey: value\n"
	err := writer.WriteString(ctx, "/tmp/config.yaml", content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("File written (or skipped if exists)")
	return nil
}

// Example 3: Incremental update with markers
func ExampleIncrementalUpdate() error {
	ctx := context.Background()

	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))

	// First write
	content1 := ".PHONY: test\ntest:\n\tgo test ./...\n"
	err := writer.WriteString(ctx, "/tmp/Makefile", content1)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Second write with updated content
	content2 := ".PHONY: test\ntest:\n\tgo test -v ./...\n"
	err = writer.WriteString(ctx, "/tmp/Makefile", content2)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	fmt.Println("File updated incrementally")
	return nil
}

// Example 4: Using in a Generator
type MakefileGenerator struct {
	outputPath string
	content    string
}

func (g *MakefileGenerator) Generate(ctx context.Context) error {
	// Use incremental strategy for Makefile
	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))

	return writer.WriteString(ctx, g.outputPath, g.content)
}

// Example 5: Conditional strategy selection
func ExampleConditionalStrategy(overwrite bool) error {
	ctx := context.Background()

	var strategyID string
	if overwrite {
		strategyID = strategies.OverwriteStrategyID
	} else {
		strategyID = strategies.IncrementalStrategyID
	}

	writer := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategyID))

	content := "# Generated content\n"
	err := writer.WriteString(ctx, "/tmp/output.txt", content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("File written with selected strategy")
	return nil
}

// Example 6: Multiple files with different strategies
func ExampleMultipleFiles() error {
	ctx := context.Background()

	// Overwrite strategy for generated files
	overwriteWriter := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.OverwriteStrategyID))

	// Skip strategy for user configuration files
	skipWriter := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.SkipStrategyID))

	// Incremental strategy for build files
	incrementalWriter := filewriter.New().
		WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))

	// Write different files with different strategies
	if err := overwriteWriter.WriteString(ctx, "/tmp/generated.go", "package main\n"); err != nil {
		return err
	}

	if err := skipWriter.WriteString(ctx, "/tmp/config.yaml", "key: value\n"); err != nil {
		return err
	}

	if err := incrementalWriter.WriteString(ctx, "/tmp/Makefile", ".PHONY: build\n"); err != nil {
		return err
	}

	fmt.Println("Multiple files written with different strategies")
	return nil
}
