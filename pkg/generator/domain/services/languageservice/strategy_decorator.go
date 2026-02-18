package languageservice

// StrategyDecorator is the base decorator for LanguageStrategy
// It implements the Decorator Pattern to add functionality to strategies
type StrategyDecorator struct {
	wrapped LanguageStrategy
}

// NewStrategyDecorator creates a new strategy decorator
func NewStrategyDecorator(strategy LanguageStrategy) *StrategyDecorator {
	return &StrategyDecorator{
		wrapped: strategy,
	}
}

// GetName delegates to the wrapped strategy
func (d *StrategyDecorator) GetName() string {
	return d.wrapped.GetName()
}

// GetDependencyFiles delegates to the wrapped strategy
func (d *StrategyDecorator) GetDependencyFiles() []string {
	return d.wrapped.GetDependencyFiles()
}

// GetDepsInstallCommand delegates to the wrapped strategy
func (d *StrategyDecorator) GetDepsInstallCommand() string {
	return d.wrapped.GetDepsInstallCommand()
}

// GetPackageManager delegates to the wrapped strategy
func (d *StrategyDecorator) GetPackageManager() string {
	return d.wrapped.GetPackageManager()
}

// GetDefaultBuildCommand delegates to the wrapped strategy
func (d *StrategyDecorator) GetDefaultBuildCommand() string {
	return d.wrapped.GetDefaultBuildCommand()
}

// GetDependencyFilesWithDetection delegates to the wrapped strategy
func (d *StrategyDecorator) GetDependencyFilesWithDetection(projectDir string) []string {
	return d.wrapped.GetDependencyFilesWithDetection(projectDir)
}

// Unwrap returns the wrapped strategy
func (d *StrategyDecorator) Unwrap() LanguageStrategy {
	return d.wrapped
}
