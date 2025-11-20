// Package testutil provides utilities for testing service configurations.
//
// This package implements three design patterns to simplify test configuration creation:
//
// 1. Builder Pattern (ConfigBuilder)
//   - Provides fluent API for building configurations
//   - Example:
//     cfg := testutil.NewConfigBuilder().
//     WithService("my-service", "My Service").
//     WithLanguage("go").
//     WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
//     WithBuilderImage("@builders.go_1.21").
//     Build()
//
// 2. Preset Pattern (Presets)
//   - Provides pre-configured templates for common scenarios
//   - Example:
//     cfg := testutil.GoServiceConfig()  // Returns a complete Go service config
//
// 3. Options Pattern (ConfigOption)
//   - Provides functional options for flexible configuration modification
//   - Example:
//     cfg := testutil.NewConfigWithOptions(
//     testutil.MinimalConfig(),
//     testutil.WithServiceNameOpt("custom-service"),
//     testutil.WithPortOpt("http", 8080, "TCP", true),
//     )
//
// Usage Examples:
//
// Example 1: Using Builder Pattern
//
//	func TestMyFeature(t *testing.T) {
//	    cfg := testutil.NewConfigBuilder().
//	        WithService("test-service", "Test Service").
//	        WithLanguage("go").
//	        WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
//	        WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
//	        WithBuilderImage("@builders.go_1.21").
//	        WithRuntimeImage("@runtimes.alpine_3.18").
//	        WithBuildCommand("go build -o bin/app").
//	        Build()
//
//	    // Use cfg in your test
//	}
//
// Example 2: Using Preset Pattern
//
//	func TestGoService(t *testing.T) {
//	    cfg := testutil.GoServiceConfig()
//	    // cfg is ready to use with sensible defaults
//	}
//
// Example 3: Using Options Pattern
//
//	func TestCustomConfig(t *testing.T) {
//	    cfg := testutil.NewConfigWithOptions(
//	        testutil.MinimalConfig(),
//	        testutil.WithServiceNameOpt("my-service"),
//	        testutil.WithBuildCommandOpt("make build"),
//	    )
//	}
//
// Example 4: Combining Patterns
//
//	func TestCombined(t *testing.T) {
//	    // Start with a preset
//	    cfg := testutil.GoServiceConfig()
//
//	    // Modify using options
//	    cfg = testutil.ApplyOptions(cfg,
//	        testutil.WithServiceNameOpt("custom-go-service"),
//	        testutil.WithPortOpt("grpc", 9000, "TCP", true),
//	    )
//	}
//
// Benefits:
//
// - High Cohesion: Configuration creation logic is centralized
// - Low Coupling: Tests are decoupled from configuration structure
// - Easy Maintenance: Configuration changes only affect this package
// - Better Readability: Test intent is clear and concise
// - Extensibility: Easy to add new presets and options
package testutil
