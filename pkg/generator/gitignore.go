package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	gitignoreStartMarker = "# >>> svcgen generated - DO NOT EDIT >>>"
	gitignoreEndMarker   = "# <<< svcgen generated <<<"
)

// gitignoreEntries 收集所有需要被 .gitignore 忽略的生成文件条目
func (g *Generator) gitignoreEntries() []string {
	entries := make(map[string]struct{})

	// 1. .tad/ 目录（包含 devops.yaml、Dockerfile、所有构建/部署脚本）
	entries[".tad/"] = struct{}{}

	// 2. 如果 CI script_dir 自定义且不在 .tad/ 下，额外添加
	scriptDir := g.ctx.Paths.CI.ScriptDir
	if !strings.HasPrefix(scriptDir, ".tad/") && !strings.HasPrefix(scriptDir, ".tad\\") {
		// 自定义路径，需要单独忽略
		// 确保以 / 结尾表示目录
		dir := strings.TrimRight(scriptDir, "/\\") + "/"
		entries[dir] = struct{}{}
	}

	// 3. compose.yaml
	entries["compose.yaml"] = struct{}{}

	// 4. Makefile
	entries["Makefile"] = struct{}{}

	// 5. .env.make（由 Makefile 从 devops.yaml 解析生成）
	entries[".env.make"] = struct{}{}

	// 排序以保证输出稳定
	sorted := make([]string, 0, len(entries))
	for e := range entries {
		sorted = append(sorted, e)
	}
	sort.Strings(sorted)

	return sorted
}

// buildGitignoreBlock 构建 marker block 内容
func buildGitignoreBlock(entries []string) string {
	var sb strings.Builder
	sb.WriteString(gitignoreStartMarker)
	sb.WriteByte('\n')
	sb.WriteString("# Generated files (only service.yaml needs to be tracked)")
	sb.WriteByte('\n')
	for _, entry := range entries {
		sb.WriteString(entry)
		sb.WriteByte('\n')
	}
	sb.WriteString(gitignoreEndMarker)
	return sb.String()
}

// updateGitignore 更新目标目录下的 .gitignore 文件
// - 如果不存在：创建新文件
// - 如果已存在：用 marker block 增量更新，保留用户自定义内容
func (g *Generator) updateGitignore() error {
	entries := g.gitignoreEntries()
	newBlock := buildGitignoreBlock(entries)

	gitignorePath := filepath.Join(g.outputDir, ".gitignore")

	existing, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建新文件
			if writeErr := os.WriteFile(gitignorePath, []byte(newBlock+"\n"), 0644); writeErr != nil {
				return fmt.Errorf("failed to create .gitignore: %w", writeErr)
			}
			fmt.Println("✓ Created .gitignore")
			return nil
		}
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	// 文件已存在，检查是否有 marker block
	updated := replaceOrAppendBlock(string(existing), newBlock)

	if updated == string(existing) {
		// 内容未变化
		fmt.Println("✓ .gitignore is up to date")
		return nil
	}

	if writeErr := os.WriteFile(gitignorePath, []byte(updated), 0644); writeErr != nil {
		return fmt.Errorf("failed to update .gitignore: %w", writeErr)
	}
	fmt.Println("✓ Updated .gitignore")
	return nil
}

// replaceOrAppendBlock 替换已有 marker block 或追加新 block
func replaceOrAppendBlock(existing, newBlock string) string {
	pattern := fmt.Sprintf(`(?s)%s\n.*?%s`,
		regexp.QuoteMeta(gitignoreStartMarker),
		regexp.QuoteMeta(gitignoreEndMarker))

	re := regexp.MustCompile(pattern)

	if re.MatchString(existing) {
		// 替换已有 block
		return re.ReplaceAllString(existing, newBlock)
	}

	// 追加新 block
	result := strings.TrimRight(existing, "\n")
	if result != "" {
		result += "\n\n"
	}
	result += newBlock + "\n"
	return result
}
