package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// LanguageStrategy defines the interface for language-specific logic
type LanguageStrategy interface {
	// GetName returns the language name
	GetName() string

	// GetDependencyFiles returns the list of dependency files
	GetDependencyFiles() []string

	// GetDepsInstallCommand returns the dependency installation command
	GetDepsInstallCommand() string

	// GetPackageManager returns the package manager name
	GetPackageManager() string

	// GetDependencyFilesWithDetection returns dependency files that actually exist in the project
	// projectDir: the root directory of the project to scan
	GetDependencyFilesWithDetection(projectDir string) []string
}

// LanguageService manages language-specific logic
type LanguageService struct {
	strategies map[string]LanguageStrategy
}

// NewLanguageService creates a new language service with built-in strategies
func NewLanguageService() *LanguageService {
	service := &LanguageService{
		strategies: make(map[string]LanguageStrategy),
	}

	// Register built-in language strategies
	service.Register(NewGoStrategy())
	service.Register(NewPythonStrategy())
	service.Register(NewNodeJSStrategy())
	service.Register(NewJavaStrategy())
	service.Register(NewRustStrategy())

	return service
}

// Register registers a language strategy
func (s *LanguageService) Register(strategy LanguageStrategy) {
	s.strategies[strategy.GetName()] = strategy
}

// GetStrategy returns the language strategy for the given language
func (s *LanguageService) GetStrategy(language string) (LanguageStrategy, error) {
	strategy, exists := s.strategies[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return strategy, nil
}

// GetDependencyFiles returns dependency files for the given language
func (s *LanguageService) GetDependencyFiles(language string, autoDetect bool, customFiles []string) []string {
	if !autoDetect {
		return customFiles
	}

	strategy, err := s.GetStrategy(language)
	if err != nil {
		return []string{}
	}

	return strategy.GetDependencyFiles()
}

// GetDependencyFilesWithDetection returns dependency files that actually exist in the project
func (s *LanguageService) GetDependencyFilesWithDetection(language string, projectDir string, autoDetect bool, customFiles []string) []string {
	if !autoDetect {
		return customFiles
	}

	strategy, err := s.GetStrategy(language)
	if err != nil {
		return []string{}
	}

	return strategy.GetDependencyFilesWithDetection(projectDir)
}

// GetDepsInstallCommand returns the dependency installation command
func (s *LanguageService) GetDepsInstallCommand(language string) string {
	strategy, err := s.GetStrategy(language)
	if err != nil {
		return "echo 'No dependency installation needed'"
	}

	return strategy.GetDepsInstallCommand()
}

// GetPackageManager returns the package manager for the given language
func (s *LanguageService) GetPackageManager(language string) string {
	strategy, err := s.GetStrategy(language)
	if err != nil {
		return "unknown"
	}

	return strategy.GetPackageManager()
}

// IsSupported checks if the language is supported
func (s *LanguageService) IsSupported(language string) bool {
	_, exists := s.strategies[language]
	return exists
}

// ListSupportedLanguages returns a list of all supported languages
func (s *LanguageService) ListSupportedLanguages() []string {
	languages := make([]string, 0, len(s.strategies))
	for lang := range s.strategies {
		languages = append(languages, lang)
	}
	return languages
}

// --- Go Language Strategy ---

// GoStrategy implements LanguageStrategy for Go
type GoStrategy struct{}

// NewGoStrategy creates a new Go strategy
func NewGoStrategy() *GoStrategy {
	return &GoStrategy{}
}

func (s *GoStrategy) GetName() string {
	return "go"
}

func (s *GoStrategy) GetDependencyFiles() []string {
	return []string{"go.mod", "go.sum"}
}

func (s *GoStrategy) GetDependencyFilesWithDetection(projectDir string) []string {
	return filterExistingFiles(projectDir, s.GetDependencyFiles())
}

func (s *GoStrategy) GetDepsInstallCommand() string {
	return "go mod download"
}

func (s *GoStrategy) GetPackageManager() string {
	return "go"
}

// --- Python Language Strategy ---

// PythonStrategy implements LanguageStrategy for Python
type PythonStrategy struct{}

// NewPythonStrategy creates a new Python strategy
func NewPythonStrategy() *PythonStrategy {
	return &PythonStrategy{}
}

func (s *PythonStrategy) GetName() string {
	return "python"
}

func (s *PythonStrategy) GetDependencyFiles() []string {
	return []string{"requirements.txt"}
}

func (s *PythonStrategy) GetDependencyFilesWithDetection(projectDir string) []string {
	return filterExistingFiles(projectDir, s.GetDependencyFiles())
}

func (s *PythonStrategy) GetDepsInstallCommand() string {
	return "pip install -r requirements.txt"
}

func (s *PythonStrategy) GetPackageManager() string {
	return "pip"
}

// --- NodeJS Language Strategy ---

// NodeJSStrategy implements LanguageStrategy for NodeJS
type NodeJSStrategy struct{}

// NewNodeJSStrategy creates a new NodeJS strategy
func NewNodeJSStrategy() *NodeJSStrategy {
	return &NodeJSStrategy{}
}

func (s *NodeJSStrategy) GetName() string {
	return "nodejs"
}

func (s *NodeJSStrategy) GetDependencyFiles() []string {
	return []string{"package.json", "package-lock.json"}
}

func (s *NodeJSStrategy) GetDependencyFilesWithDetection(projectDir string) []string {
	return filterExistingFiles(projectDir, s.GetDependencyFiles())
}

func (s *NodeJSStrategy) GetDepsInstallCommand() string {
	return "npm install"
}

func (s *NodeJSStrategy) GetPackageManager() string {
	return "npm"
}

// --- Java Language Strategy ---

// JavaStrategy implements LanguageStrategy for Java
// Supports both Maven and Gradle build tools
type JavaStrategy struct{}

// NewJavaStrategy creates a new Java strategy
func NewJavaStrategy() *JavaStrategy {
	return &JavaStrategy{}
}

func (s *JavaStrategy) GetName() string {
	return "java"
}

func (s *JavaStrategy) GetDependencyFiles() []string {
	// Support both Maven and Gradle
	// Maven: pom.xml
	// Gradle: build.gradle, settings.gradle
	return []string{"pom.xml", "build.gradle", "settings.gradle"}
}

func (s *JavaStrategy) GetDependencyFilesWithDetection(projectDir string) []string {
	var detectedFiles []string

	// Check for Maven (pom.xml)
	pomPath := filepath.Join(projectDir, "pom.xml")
	if fileExists(pomPath) {
		detectedFiles = append(detectedFiles, "pom.xml")
	}

	// Check for Gradle (build.gradle and settings.gradle)
	buildGradlePath := filepath.Join(projectDir, "build.gradle")
	settingsGradlePath := filepath.Join(projectDir, "settings.gradle")

	if fileExists(buildGradlePath) {
		detectedFiles = append(detectedFiles, "build.gradle")
	}
	if fileExists(settingsGradlePath) {
		detectedFiles = append(detectedFiles, "settings.gradle")
	}

	return detectedFiles
}

func (s *JavaStrategy) GetDepsInstallCommand() string {
	// Return Maven command by default
	// Note: In practice, the build script should detect which build tool is present
	// and use the appropriate command (mvn or gradle)
	return "mvn dependency:go-offline || gradle dependencies --refresh-dependencies"
}

func (s *JavaStrategy) GetPackageManager() string {
	return "mvn"
}

// --- Rust Language Strategy ---

// RustStrategy implements LanguageStrategy for Rust
type RustStrategy struct{}

// NewRustStrategy creates a new Rust strategy
func NewRustStrategy() *RustStrategy {
	return &RustStrategy{}
}

func (s *RustStrategy) GetName() string {
	return "rust"
}

func (s *RustStrategy) GetDependencyFiles() []string {
	return []string{"Cargo.toml", "Cargo.lock"}
}

func (s *RustStrategy) GetDependencyFilesWithDetection(projectDir string) []string {
	return filterExistingFiles(projectDir, s.GetDependencyFiles())
}

func (s *RustStrategy) GetDepsInstallCommand() string {
	return "cargo fetch"
}

func (s *RustStrategy) GetPackageManager() string {
	return "cargo"
}

// --- Helper Functions ---

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// filterExistingFiles filters the list of files to only include those that exist
func filterExistingFiles(projectDir string, files []string) []string {
	var existingFiles []string
	for _, file := range files {
		filePath := filepath.Join(projectDir, file)
		if fileExists(filePath) {
			existingFiles = append(existingFiles, file)
		}
	}
	// If no files exist, return the original list as fallback
	if len(existingFiles) == 0 {
		return files
	}
	return existingFiles
}
