package generator

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/config"
)

// HealthcheckStrategy defines the interface for health check strategies
type HealthcheckStrategy interface {
	// GetType returns the strategy type
	GetType() string
	// GenerateScript generates the health check script content
	GenerateScript(vars map[string]interface{}) (string, error)
	// Validate validates the strategy configuration
	Validate() error
}

// HealthcheckStrategyFactory creates health check strategies
type HealthcheckStrategyFactory struct {
	config *config.ServiceConfig
}

// NewHealthcheckStrategyFactory creates a new factory
func NewHealthcheckStrategyFactory(cfg *config.ServiceConfig) *HealthcheckStrategyFactory {
	return &HealthcheckStrategyFactory{
		config: cfg,
	}
}

// CreateStrategy creates a health check strategy based on configuration
func (f *HealthcheckStrategyFactory) CreateStrategy() (HealthcheckStrategy, error) {
	if !f.config.Runtime.Healthcheck.Enabled {
		return NewDefaultHealthcheckStrategy(f.config), nil
	}

	switch f.config.Runtime.Healthcheck.Type {
	case "default", "":
		return NewDefaultHealthcheckStrategy(f.config), nil
	case "custom":
		return NewCustomHealthcheckStrategy(f.config), nil
	default:
		return nil, fmt.Errorf("unsupported healthcheck type: %s (valid: default, custom)", f.config.Runtime.Healthcheck.Type)
	}
}

// DefaultHealthcheckStrategy implements default process-based health check
type DefaultHealthcheckStrategy struct {
	config *config.ServiceConfig
}

// NewDefaultHealthcheckStrategy creates a new default strategy
func NewDefaultHealthcheckStrategy(cfg *config.ServiceConfig) *DefaultHealthcheckStrategy {
	return &DefaultHealthcheckStrategy{
		config: cfg,
	}
}

// GetType returns the strategy type
func (s *DefaultHealthcheckStrategy) GetType() string {
	return "default"
}

// GenerateScript generates default health check script
func (s *DefaultHealthcheckStrategy) GenerateScript(vars map[string]interface{}) (string, error) {
	script := `#!/bin/sh

# Export service paths as environment variables
export SERVICE_ROOT="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}"
export SERVICE_BIN_DIR="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}/bin"
export SERVICE_NAME="{{ .SERVICE_NAME }}"

# Default healthcheck: check if service process is running
ps=$(ls -l /proc/*/exe 2>/dev/null | grep "${SERVICE_NAME}" | grep -v grep)

# abnormal
[[ "$ps" == "" ]] && exit 1

# normal
exit 0
`
	return script, nil
}

// Validate validates the default strategy configuration
func (s *DefaultHealthcheckStrategy) Validate() error {
	// Default strategy has no additional validation requirements
	return nil
}

// CustomHealthcheckStrategy implements custom user-defined health check
type CustomHealthcheckStrategy struct {
	config *config.ServiceConfig
}

// NewCustomHealthcheckStrategy creates a new custom strategy
func NewCustomHealthcheckStrategy(cfg *config.ServiceConfig) *CustomHealthcheckStrategy {
	return &CustomHealthcheckStrategy{
		config: cfg,
	}
}

// GetType returns the strategy type
func (s *CustomHealthcheckStrategy) GetType() string {
	return "custom"
}

// GenerateScript generates custom health check script
func (s *CustomHealthcheckStrategy) GenerateScript(vars map[string]interface{}) (string, error) {
	if s.config.Runtime.Healthcheck.CustomScript == "" {
		return "", fmt.Errorf("custom_script is required for custom healthcheck type")
	}

	script := `#!/bin/sh

# Export service paths as environment variables
export SERVICE_ROOT="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}"
export SERVICE_BIN_DIR="{{ .DEPLOY_DIR }}/{{ .SERVICE_NAME }}/bin"
export SERVICE_NAME="{{ .SERVICE_NAME }}"

# Custom healthcheck script
{{ .CUSTOM_SCRIPT }}
`
	return script, nil
}

// Validate validates the custom strategy configuration
func (s *CustomHealthcheckStrategy) Validate() error {
	if s.config.Runtime.Healthcheck.CustomScript == "" {
		return fmt.Errorf("runtime.healthcheck.custom_script is required when type is 'custom'")
	}
	return nil
}
