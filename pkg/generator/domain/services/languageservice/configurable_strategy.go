package languageservice

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// ConfigurableStrategy decorates a strategy with custom configuration support
// It allows overriding the default install command from config
type ConfigurableStrategy struct {
	*StrategyDecorator
	config *config.LanguageConfig
}

// NewConfigurableStrategy creates a new configurable strategy decorator
func NewConfigurableStrategy(strategy LanguageStrategy, config *config.LanguageConfig) *ConfigurableStrategy {
	return &ConfigurableStrategy{
		StrategyDecorator: NewStrategyDecorator(strategy),
		config:            config,
	}
}

// GetDepsInstallCommand returns the install command
// If custom command is configured, use it; otherwise use the default from wrapped strategy
func (s *ConfigurableStrategy) GetDepsInstallCommand() string {
	// Check if custom install command is configured
	if s.config != nil {
		customCmd := s.config.GetString("deps_install_command", "")
		if customCmd != "" {
			return customCmd
		}
	}

	// Use default command from wrapped strategy
	return s.wrapped.GetDepsInstallCommand()
}
