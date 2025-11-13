package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Loader handles loading configuration from YAML files
type Loader struct {
	configPath string
}

// NewLoader creates a new configuration loader
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load reads and parses the service.yaml configuration file
func (l *Loader) Load() (*ServiceConfig, error) {
	data, err := os.ReadFile(l.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply default values
	applyDefaults(&config)

	return &config, nil
}

// LoadFromBytes loads configuration from byte slice
func LoadFromBytes(data []byte) (*ServiceConfig, error) {
	var config ServiceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply default values
	applyDefaults(&config)

	return &config, nil
}

// Save writes the configuration to a YAML file
func (l *Loader) Save(config *ServiceConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(l.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// applyDefaults applies default values to the configuration
func applyDefaults(config *ServiceConfig) {
	// Default deploy directory
	if config.Service.DeployDir == "" {
		config.Service.DeployDir = "/usr/local/services"
	}

	// Default auto_detect for dependency files
	// Note: YAML unmarshaling sets bool to false by default
	// We need to check if it was explicitly set or not
	// Since we can't distinguish between "not set" and "false" with bool,
	// we treat missing field as "true" by checking if Files is also empty
	if !config.Build.DependencyFiles.AutoDetect && len(config.Build.DependencyFiles.Files) == 0 {
		config.Build.DependencyFiles.AutoDetect = true
	}
}
