package services

import (
	"os"
	"testing"
)

func TestLanguageService_GetStrategy(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{"go language", "go", false},
		{"python language", "python", false},
		{"nodejs language", "nodejs", false},
		{"java language", "java", false},
		{"rust language", "rust", false},
		{"unsupported language", "ruby", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := service.GetStrategy(tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && strategy == nil {
				t.Error("GetStrategy() returned nil strategy")
			}
		})
	}
}

func TestLanguageService_GetDependencyFiles(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name        string
		language    string
		autoDetect  bool
		customFiles []string
		expected    []string
	}{
		{
			name:       "go auto detect",
			language:   "go",
			autoDetect: true,
			expected:   []string{"go.mod", "go.sum"},
		},
		{
			name:       "python auto detect",
			language:   "python",
			autoDetect: true,
			expected:   []string{"requirements.txt"},
		},
		{
			name:        "custom files",
			language:    "go",
			autoDetect:  false,
			customFiles: []string{"custom.txt"},
			expected:    []string{"custom.txt"},
		},
		{
			name:       "unsupported language",
			language:   "unknown",
			autoDetect: true,
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetDependencyFiles(tt.language, tt.autoDetect, tt.customFiles)
			if len(got) != len(tt.expected) {
				t.Errorf("GetDependencyFiles() length = %d, want %d", len(got), len(tt.expected))
				return
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("GetDependencyFiles()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestLanguageService_GetDepsInstallCommand(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name     string
		language string
		expected string
	}{
		{"go", "go", "go mod download"},
		{"python", "python", "pip install -r requirements.txt"},
		{"nodejs", "nodejs", "npm install"},
		{"java", "java", "mvn dependency:go-offline || gradle dependencies --refresh-dependencies"},
		{"rust", "rust", "cargo fetch"},
		{"unknown", "unknown", "echo 'No dependency installation needed'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetDepsInstallCommand(tt.language)
			if got != tt.expected {
				t.Errorf("GetDepsInstallCommand() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLanguageService_GetPackageManager(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name     string
		language string
		expected string
	}{
		{"go", "go", "go"},
		{"python", "python", "pip"},
		{"nodejs", "nodejs", "npm"},
		{"java", "java", "mvn"},
		{"rust", "rust", "cargo"},
		{"unknown", "unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetPackageManager(tt.language)
			if got != tt.expected {
				t.Errorf("GetPackageManager() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLanguageService_IsSupported(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name     string
		language string
		expected bool
	}{
		{"go supported", "go", true},
		{"python supported", "python", true},
		{"ruby not supported", "ruby", false},
		{"php not supported", "php", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.IsSupported(tt.language)
			if got != tt.expected {
				t.Errorf("IsSupported() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLanguageService_ListSupportedLanguages(t *testing.T) {
	service := NewLanguageService()

	languages := service.ListSupportedLanguages()

	if len(languages) != 5 {
		t.Errorf("ListSupportedLanguages() length = %d, want 5", len(languages))
	}

	// Check that all expected languages are present
	expectedLanguages := map[string]bool{
		"go":     true,
		"python": true,
		"nodejs": true,
		"java":   true,
		"rust":   true,
	}

	for _, lang := range languages {
		if !expectedLanguages[lang] {
			t.Errorf("Unexpected language in list: %s", lang)
		}
	}
}

func TestLanguageService_Register(t *testing.T) {
	service := NewLanguageService()

	// Create a custom strategy
	customStrategy := &GoStrategy{} // Reuse GoStrategy for testing

	// Register with a different name
	service.Register(customStrategy)

	// Verify it's registered
	strategy, err := service.GetStrategy("go")
	if err != nil {
		t.Errorf("GetStrategy() error = %v", err)
	}
	if strategy == nil {
		t.Error("GetStrategy() returned nil")
	}
}

// Test individual strategies
func TestGoStrategy(t *testing.T) {
	strategy := NewGoStrategy()

	if strategy.GetName() != "go" {
		t.Errorf("GetName() = %v, want go", strategy.GetName())
	}

	files := strategy.GetDependencyFiles()
	if len(files) != 2 || files[0] != "go.mod" || files[1] != "go.sum" {
		t.Errorf("GetDependencyFiles() = %v, want [go.mod go.sum]", files)
	}

	if strategy.GetDepsInstallCommand() != "go mod download" {
		t.Errorf("GetDepsInstallCommand() = %v, want go mod download", strategy.GetDepsInstallCommand())
	}

	if strategy.GetPackageManager() != "go" {
		t.Errorf("GetPackageManager() = %v, want go", strategy.GetPackageManager())
	}
}

func TestPythonStrategy(t *testing.T) {
	strategy := NewPythonStrategy()

	if strategy.GetName() != "python" {
		t.Errorf("GetName() = %v, want python", strategy.GetName())
	}

	files := strategy.GetDependencyFiles()
	if len(files) != 1 || files[0] != "requirements.txt" {
		t.Errorf("GetDependencyFiles() = %v, want [requirements.txt]", files)
	}
}

func TestNodeJSStrategy(t *testing.T) {
	strategy := NewNodeJSStrategy()

	if strategy.GetName() != "nodejs" {
		t.Errorf("GetName() = %v, want nodejs", strategy.GetName())
	}

	files := strategy.GetDependencyFiles()
	if len(files) != 2 || files[0] != "package.json" {
		t.Errorf("GetDependencyFiles() = %v, want [package.json package-lock.json]", files)
	}
}

func TestJavaStrategy(t *testing.T) {
	strategy := NewJavaStrategy()

	if strategy.GetName() != "java" {
		t.Errorf("GetName() = %v, want java", strategy.GetName())
	}

	files := strategy.GetDependencyFiles()
	expected := []string{"pom.xml", "build.gradle", "settings.gradle"}
	if len(files) != len(expected) {
		t.Errorf("GetDependencyFiles() length = %v, want %v", len(files), len(expected))
	}
	for i, file := range expected {
		if i >= len(files) || files[i] != file {
			t.Errorf("GetDependencyFiles() = %v, want %v", files, expected)
			break
		}
	}
}

func TestRustStrategy(t *testing.T) {
	strategy := NewRustStrategy()

	if strategy.GetName() != "rust" {
		t.Errorf("GetName() = %v, want rust", strategy.GetName())
	}

	files := strategy.GetDependencyFiles()
	if len(files) != 2 || files[0] != "Cargo.toml" {
		t.Errorf("GetDependencyFiles() = %v, want [Cargo.toml Cargo.lock]", files)
	}
}

func TestLanguageService_GetDependencyFilesWithDetection(t *testing.T) {
	service := NewLanguageService()

	tests := []struct {
		name        string
		language    string
		projectDir  string
		autoDetect  bool
		customFiles []string
		setupFiles  []string // Files to create for testing
		expected    []string
	}{
		{
			name:       "java maven project",
			language:   "java",
			projectDir: t.TempDir(),
			autoDetect: true,
			setupFiles: []string{"pom.xml"},
			expected:   []string{"pom.xml"},
		},
		{
			name:       "java gradle project",
			language:   "java",
			projectDir: t.TempDir(),
			autoDetect: true,
			setupFiles: []string{"build.gradle", "settings.gradle"},
			expected:   []string{"build.gradle", "settings.gradle"},
		},
		{
			name:       "java mixed project (both maven and gradle)",
			language:   "java",
			projectDir: t.TempDir(),
			autoDetect: true,
			setupFiles: []string{"pom.xml", "build.gradle", "settings.gradle"},
			expected:   []string{"pom.xml", "build.gradle", "settings.gradle"},
		},
		{
			name:        "custom files without detection",
			language:    "java",
			projectDir:  t.TempDir(),
			autoDetect:  false,
			customFiles: []string{"custom-deps.txt"},
			expected:    []string{"custom-deps.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test files
			for _, file := range tt.setupFiles {
				filePath := tt.projectDir + "/" + file
				f, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", file, err)
				}
				f.Close()
			}

			got := service.GetDependencyFilesWithDetection(tt.language, tt.projectDir, tt.autoDetect, tt.customFiles)
			if len(got) != len(tt.expected) {
				t.Errorf("GetDependencyFilesWithDetection() length = %d, want %d", len(got), len(tt.expected))
				t.Errorf("Got: %v, Want: %v", got, tt.expected)
				return
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("GetDependencyFilesWithDetection()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestJavaStrategy_GetDependencyFilesWithDetection(t *testing.T) {
	strategy := NewJavaStrategy()

	tests := []struct {
		name       string
		setupFiles []string
		expected   []string
	}{
		{
			name:       "maven only",
			setupFiles: []string{"pom.xml"},
			expected:   []string{"pom.xml"},
		},
		{
			name:       "gradle only",
			setupFiles: []string{"build.gradle", "settings.gradle"},
			expected:   []string{"build.gradle", "settings.gradle"},
		},
		{
			name:       "gradle build.gradle only",
			setupFiles: []string{"build.gradle"},
			expected:   []string{"build.gradle"},
		},
		{
			name:       "both maven and gradle",
			setupFiles: []string{"pom.xml", "build.gradle", "settings.gradle"},
			expected:   []string{"pom.xml", "build.gradle", "settings.gradle"},
		},
		{
			name:       "no dependency files (fallback)",
			setupFiles: []string{},
			expected:   []string{"pom.xml", "build.gradle", "settings.gradle"}, // Returns all as fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Setup test files
			for _, file := range tt.setupFiles {
				filePath := tmpDir + "/" + file
				f, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", file, err)
				}
				f.Close()
			}

			got := strategy.GetDependencyFilesWithDetection(tmpDir)
			if len(got) != len(tt.expected) {
				t.Errorf("GetDependencyFilesWithDetection() length = %d, want %d", len(got), len(tt.expected))
				t.Errorf("Got: %v, Want: %v", got, tt.expected)
				return
			}
			for i, v := range got {
				if v != tt.expected[i] {
					t.Errorf("GetDependencyFilesWithDetection()[%d] = %v, want %v", i, v, tt.expected[i])
				}
			}
		})
	}
}
