package generator

import (
	"os"
	"path/filepath"
	"testing"

	configtestutil "github.com/junjiewwang/service-template/pkg/config/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildGitignoreBlock(t *testing.T) {
	entries := []string{".env.make", ".tad/", "Makefile", "compose.yaml"}

	block := buildGitignoreBlock(entries)

	assert.Contains(t, block, gitignoreStartMarker)
	assert.Contains(t, block, gitignoreEndMarker)
	assert.Contains(t, block, ".tad/")
	assert.Contains(t, block, "compose.yaml")
	assert.Contains(t, block, "Makefile")
	assert.Contains(t, block, ".env.make")
	assert.Contains(t, block, "# Generated files (only service.yaml needs to be tracked)")
}

func TestReplaceOrAppendBlock_NoExisting(t *testing.T) {
	newBlock := buildGitignoreBlock([]string{".tad/", "compose.yaml"})

	result := replaceOrAppendBlock("", newBlock)

	assert.Contains(t, result, gitignoreStartMarker)
	assert.Contains(t, result, gitignoreEndMarker)
	assert.Contains(t, result, ".tad/")
	assert.Contains(t, result, "compose.yaml")
}

func TestReplaceOrAppendBlock_AppendToExisting(t *testing.T) {
	existing := `# My custom ignores
*.log
vendor/
`
	newBlock := buildGitignoreBlock([]string{".tad/", "compose.yaml"})

	result := replaceOrAppendBlock(existing, newBlock)

	// 应保留用户自定义内容
	assert.Contains(t, result, "*.log")
	assert.Contains(t, result, "vendor/")
	// 应包含生成的 block
	assert.Contains(t, result, gitignoreStartMarker)
	assert.Contains(t, result, ".tad/")
	assert.Contains(t, result, "compose.yaml")
}

func TestReplaceOrAppendBlock_ReplaceExistingBlock(t *testing.T) {
	existing := `# My custom ignores
*.log

# >>> svcgen generated - DO NOT EDIT >>>
# Generated files (only service.yaml needs to be tracked)
.tad/
compose.yaml
# <<< svcgen generated <<<
`
	// 新 block 新增了 Makefile
	newBlock := buildGitignoreBlock([]string{".tad/", "Makefile", "compose.yaml"})

	result := replaceOrAppendBlock(existing, newBlock)

	// 应保留用户自定义内容
	assert.Contains(t, result, "*.log")
	// 应包含更新后的 block
	assert.Contains(t, result, "Makefile")
	assert.Contains(t, result, ".tad/")
	assert.Contains(t, result, "compose.yaml")

	// 确保只有一个 marker block（不重复追加）
	assert.Equal(t, 1, countOccurrences(result, gitignoreStartMarker))
	assert.Equal(t, 1, countOccurrences(result, gitignoreEndMarker))
}

func TestReplaceOrAppendBlock_PreservesUserContent(t *testing.T) {
	existing := `# IDE files
.vscode/
.idea/

# >>> svcgen generated - DO NOT EDIT >>>
# Generated files (only service.yaml needs to be tracked)
.tad/
# <<< svcgen generated <<<

# Build artifacts
dist/
build/
`
	newBlock := buildGitignoreBlock([]string{".tad/", "Makefile", "compose.yaml"})

	result := replaceOrAppendBlock(existing, newBlock)

	// 保留 block 前后的用户内容
	assert.Contains(t, result, ".vscode/")
	assert.Contains(t, result, ".idea/")
	assert.Contains(t, result, "dist/")
	assert.Contains(t, result, "build/")
	// block 已更新
	assert.Contains(t, result, "Makefile")
}

func TestGenerator_GitignoreEntries_DefaultPaths(t *testing.T) {
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, "/tmp/test-output")
	entries := gen.gitignoreEntries()

	assert.Contains(t, entries, ".tad/")
	assert.Contains(t, entries, "compose.yaml")
	assert.Contains(t, entries, "Makefile")
	assert.Contains(t, entries, ".env.make")
}

func TestGenerator_GitignoreEntries_CustomScriptDir(t *testing.T) {
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		WithCIScriptDir("bk-ci/tcs").
		BuildWithDefaults()

	gen := NewGenerator(cfg, "/tmp/test-output")
	entries := gen.gitignoreEntries()

	// 自定义路径不在 .tad/ 下，应同时包含 .tad/ 和自定义路径
	assert.Contains(t, entries, ".tad/")
	assert.Contains(t, entries, "bk-ci/tcs/")
	assert.Contains(t, entries, "compose.yaml")
	assert.Contains(t, entries, "Makefile")
	assert.Contains(t, entries, ".env.make")
}

func TestGenerator_GitignoreEntries_ScriptDirUnderTad(t *testing.T) {
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		WithCIScriptDir(".tad/custom/scripts").
		BuildWithDefaults()

	gen := NewGenerator(cfg, "/tmp/test-output")
	entries := gen.gitignoreEntries()

	// script_dir 在 .tad/ 下，不应重复添加
	assert.Contains(t, entries, ".tad/")
	// 不应该有额外的 .tad/custom/scripts/ 条目（因为已被 .tad/ 覆盖）
	assert.NotContains(t, entries, ".tad/custom/scripts/")
}

func TestGenerator_UpdateGitignore_CreateNew(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	err := gen.updateGitignore()
	require.NoError(t, err)

	// 验证文件已创建
	content, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, gitignoreStartMarker)
	assert.Contains(t, contentStr, gitignoreEndMarker)
	assert.Contains(t, contentStr, ".tad/")
	assert.Contains(t, contentStr, "compose.yaml")
	assert.Contains(t, contentStr, "Makefile")
	assert.Contains(t, contentStr, ".env.make")
}

func TestGenerator_UpdateGitignore_AppendToExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建已有 .gitignore
	existingContent := "# My rules\n*.log\nvendor/\n"
	err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(existingContent), 0644)
	require.NoError(t, err)

	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	err = gen.updateGitignore()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	require.NoError(t, err)

	contentStr := string(content)
	// 保留用户内容
	assert.Contains(t, contentStr, "*.log")
	assert.Contains(t, contentStr, "vendor/")
	// 包含生成内容
	assert.Contains(t, contentStr, gitignoreStartMarker)
	assert.Contains(t, contentStr, ".tad/")
}

func TestGenerator_UpdateGitignore_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建已有 .gitignore（包含旧的 svcgen block）
	existingContent := `# My rules
*.log

# >>> svcgen generated - DO NOT EDIT >>>
# Generated files (only service.yaml needs to be tracked)
.tad/
compose.yaml
# <<< svcgen generated <<<
`
	err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(existingContent), 0644)
	require.NoError(t, err)

	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	err = gen.updateGitignore()
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	require.NoError(t, err)

	contentStr := string(content)
	// 保留用户内容
	assert.Contains(t, contentStr, "*.log")
	// block 已更新，包含新增的条目
	assert.Contains(t, contentStr, "Makefile")
	assert.Contains(t, contentStr, ".env.make")
	// 确保只有一个 block
	assert.Equal(t, 1, countOccurrences(contentStr, gitignoreStartMarker))
}

func TestGenerator_UpdateGitignore_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	// 第一次生成
	err := gen.updateGitignore()
	require.NoError(t, err)

	content1, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	require.NoError(t, err)

	// 第二次生成（应该幂等）
	err = gen.updateGitignore()
	require.NoError(t, err)

	content2, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	require.NoError(t, err)

	assert.Equal(t, string(content1), string(content2), "Running updateGitignore twice should produce the same result")
}

func TestGenerator_ManageGitignore_DefaultFalse(t *testing.T) {
	tmpDir := t.TempDir()

	// 默认不设置 ManageGitignore（零值为 false）
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	// 验证 ManageGitignore 默认为 false
	assert.False(t, cfg.Metadata.ManageGitignore, "ManageGitignore should default to false")

	// Generate() 不应创建 .gitignore
	err := gen.Generate()
	require.NoError(t, err)

	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	_, err = os.Stat(gitignorePath)
	assert.True(t, os.IsNotExist(err), ".gitignore should NOT be created when manage_gitignore is false")
}

func TestGenerator_ManageGitignore_EnabledTrue(t *testing.T) {
	tmpDir := t.TempDir()

	// 显式设置 ManageGitignore = true
	cfg := configtestutil.NewConfigBuilder().
		WithService("test-service", "Test Service").
		WithPort("http", 8080, "TCP", true).
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/test-service").
		WithDeployDir("/opt/services").
		WithManageGitignore(true).
		BuildWithDefaults()

	gen := NewGenerator(cfg, tmpDir)

	// Generate() 应创建 .gitignore
	err := gen.Generate()
	require.NoError(t, err)

	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	require.NoError(t, err, ".gitignore should be created when manage_gitignore is true")

	contentStr := string(content)
	assert.Contains(t, contentStr, gitignoreStartMarker)
	assert.Contains(t, contentStr, gitignoreEndMarker)
	assert.Contains(t, contentStr, ".tad/")
	assert.Contains(t, contentStr, "compose.yaml")
	assert.Contains(t, contentStr, "Makefile")
}

// countOccurrences 统计字符串出现次数
func countOccurrences(s, substr string) int {
	count := 0
	offset := 0
	for {
		idx := indexOf(s[offset:], substr)
		if idx == -1 {
			break
		}
		count++
		offset += idx + len(substr)
	}
	return count
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
