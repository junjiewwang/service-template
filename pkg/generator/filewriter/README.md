# FileWriter Framework

A flexible and extensible file writing framework that supports multiple writing strategies, including incremental updates with marker blocks.

## Features

- **Multiple Writing Strategies**: Overwrite, Skip, and Incremental
- **Marker-based Incremental Updates**: Automatically merge new content with existing files
- **Idempotent Operations**: Multiple writes with the same content produce the same result
- **Extensible Architecture**: Easy to add new strategies, mergers, and conflict resolvers
- **Self-describing Components**: All components register themselves automatically
- **Type-safe**: No hardcoded strings, all IDs are exposed as constants

## Architecture

```
FileWriter (Facade)
    ↓
WriteStrategy (Overwrite / Skip / Incremental)
    ↓
ContentMerger (Marker-based)
    ↓
ConflictResolver (KeepExisting / UseNew)
```

## Quick Start

### 1. Simple Overwrite (Default)

```go
import (
    "context"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter"
)

func main() {
    ctx := context.Background()
    writer := filewriter.New()
    
    err := writer.WriteString(ctx, "/path/to/file.txt", "content")
    if err != nil {
        // handle error
    }
}
```

### 2. Skip if File Exists

```go
import (
    "context"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter/strategies"
)

func main() {
    ctx := context.Background()
    
    writer := filewriter.New().
        WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.SkipStrategyID))
    
    err := writer.WriteString(ctx, "/path/to/file.txt", "content")
    if err != nil {
        // handle error
    }
}
```

### 3. Incremental Update with Markers

```go
import (
    "context"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter/strategies"
)

func main() {
    ctx := context.Background()
    
    writer := filewriter.New().
        WithStrategy(filewriter.DefaultStrategyRegistry.MustGet(strategies.IncrementalStrategyID))
    
    // First write: appends content with markers
    err := writer.WriteString(ctx, "/path/to/Makefile", "generated content")
    if err != nil {
        // handle error
    }
    
    // Second write with same content: no change (idempotent)
    err = writer.WriteString(ctx, "/path/to/Makefile", "generated content")
    
    // Third write with different content: updates marker block
    err = writer.WriteString(ctx, "/path/to/Makefile", "updated content")
}
```

## How Incremental Update Works

When using the incremental strategy:

1. **First Write (File doesn't exist)**: Creates file with content directly
2. **First Write (File exists, no markers)**: Appends content with markers:
   ```
   # User content
   user_var = 1
   
   # ===== GENERATED_START =====
   generated_var = 2
   # ===== GENERATED_END =====
   ```

3. **Subsequent Writes (Same content)**: No change (idempotent)
4. **Subsequent Writes (Different content)**: Updates only the marker block:
   ```
   # User content
   user_var = 1
   
   # ===== GENERATED_START =====
   generated_var = 3  # Updated!
   # ===== GENERATED_END =====
   ```

## Advanced Usage

### Custom Markers

```go
import (
    "github.com/junjiewwang/service-template/pkg/generator/filewriter/mergers"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter/strategies"
)

func main() {
    // Create custom merger with custom markers
    customMerger := mergers.NewMarkerMerger().
        WithMarkers("<!-- BEGIN GENERATED -->", "<!-- END GENERATED -->")
    
    // Register the custom merger
    mergers.DefaultMergerRegistry.MustRegister(customMerger)
    
    // Use incremental strategy with custom merger
    strategy := strategies.NewIncrementalStrategy().
        WithMerger(customMerger.ID())
    
    writer := filewriter.New().WithStrategy(strategy)
    
    // Use the writer...
}
```

## Extending the Framework

### Adding a New Write Strategy

```go
package strategies

import (
    "context"
    "github.com/junjiewwang/service-template/pkg/generator/filewriter"
)

// 1. Define the strategy ID constant
const BackupStrategyID = "backup"

// 2. Implement the WriteStrategy interface
type BackupStrategy struct{}

func init() {
    // 3. Auto-register the strategy
    filewriter.DefaultStrategyRegistry.MustRegister(&BackupStrategy{})
}

func (s *BackupStrategy) ID() string {
    return BackupStrategyID
}

func (s *BackupStrategy) Description() string {
    return "Backup existing file before overwriting"
}

func (s *BackupStrategy) Write(ctx context.Context, path string, content []byte) error {
    // Implement backup logic
    return nil
}
```

### Adding a New Content Merger

```go
package mergers

import (
    "context"
)

// 1. Define the merger ID constant
const YAMLMergerID = "yaml"

// 2. Implement the ContentMerger interface
type YAMLMerger struct{}

func init() {
    // 3. Auto-register the merger
    DefaultMergerRegistry.MustRegister(&YAMLMerger{})
}

func (m *YAMLMerger) ID() string {
    return YAMLMergerID
}

func (m *YAMLMerger) Description() string {
    return "Merge YAML files by structure"
}

func (m *YAMLMerger) Merge(ctx context.Context, input *MergeInput) ([]byte, error) {
    // Implement YAML merge logic
    return nil, nil
}
```

## Components

### Write Strategies

| Strategy | ID | Description |
|----------|----|----|
| **OverwriteStrategy** | `overwrite` | Always overwrite existing files |
| **SkipStrategy** | `skip` | Skip writing if file already exists |
| **IncrementalStrategy** | `incremental` | Merge new content with existing content using marker blocks |

### Content Mergers

| Merger | ID | Description |
|--------|----|----|
| **MarkerMerger** | `marker` | Merge content using marker blocks (start/end markers) |

### Conflict Resolvers

| Resolver | ID | Description |
|----------|----|----|
| **KeepExistingResolver** | `keep_existing` | Keep existing content when conflict occurs |
| **UseNewResolver** | `use_new` | Use new content when conflict occurs |

## Design Principles

1. **Simple by Default**: 90% of use cases require just one line of code
2. **Progressive Enhancement**: Complex scenarios can be handled with additional configuration
3. **Zero Intrusion**: Doesn't change existing Generator interfaces
4. **Self-describing**: Components expose their own IDs as constants
5. **Testable**: Each component is independently testable

## Testing

Run tests:

```bash
go test -v ./pkg/generator/filewriter/...
```

## License

This framework is part of the service-template project.
