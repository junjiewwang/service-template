package models

import "github.com/junjiewwang/service-template/pkg/generator/domain/errors"

// PathsConfig is a value object that encapsulates all path-related configuration
// This eliminates hardcoded paths throughout the codebase
type PathsConfig struct {
	// Plugin paths
	PluginInstallDir string

	// Service paths
	ServiceDeployDir string
	ServiceBinDir    string
	ServiceConfigDir string
	ServiceLogDir    string
	ServiceDataDir   string

	// Build paths
	BuildOutputDir string
	BuildCacheDir  string

	// Script paths
	ScriptDir string

	// CI/CD paths
	CIDir string
}

// NewPathsConfig creates a new PathsConfig with default values
func NewPathsConfig() *PathsConfig {
	return &PathsConfig{
		PluginInstallDir: "/plugins",
		ServiceDeployDir: "/data/services",
		ServiceBinDir:    "bin",
		ServiceConfigDir: "conf",
		ServiceLogDir:    "logs",
		ServiceDataDir:   "data",
		BuildOutputDir:   "build",
		BuildCacheDir:    ".cache",
		ScriptDir:        "scripts",
		CIDir:            ".ci",
	}
}

// WithPluginInstallDir sets the plugin installation directory
func (p *PathsConfig) WithPluginInstallDir(dir string) *PathsConfig {
	p.PluginInstallDir = dir
	return p
}

// WithServiceDeployDir sets the service deployment directory
func (p *PathsConfig) WithServiceDeployDir(dir string) *PathsConfig {
	p.ServiceDeployDir = dir
	return p
}

// WithServiceBinDir sets the service binary directory
func (p *PathsConfig) WithServiceBinDir(dir string) *PathsConfig {
	p.ServiceBinDir = dir
	return p
}

// WithServiceConfigDir sets the service configuration directory
func (p *PathsConfig) WithServiceConfigDir(dir string) *PathsConfig {
	p.ServiceConfigDir = dir
	return p
}

// WithServiceLogDir sets the service log directory
func (p *PathsConfig) WithServiceLogDir(dir string) *PathsConfig {
	p.ServiceLogDir = dir
	return p
}

// WithServiceDataDir sets the service data directory
func (p *PathsConfig) WithServiceDataDir(dir string) *PathsConfig {
	p.ServiceDataDir = dir
	return p
}

// GetServiceRoot returns the full service root path
func (p *PathsConfig) GetServiceRoot(serviceName string) string {
	return p.ServiceDeployDir + "/" + serviceName
}

// GetServiceBinPath returns the full service binary path
func (p *PathsConfig) GetServiceBinPath(serviceName string) string {
	return p.GetServiceRoot(serviceName) + "/" + p.ServiceBinDir
}

// GetServiceConfigPath returns the full service config path
func (p *PathsConfig) GetServiceConfigPath(serviceName string) string {
	return p.GetServiceRoot(serviceName) + "/" + p.ServiceConfigDir
}

// GetServiceLogPath returns the full service log path
func (p *PathsConfig) GetServiceLogPath(serviceName string) string {
	return p.GetServiceRoot(serviceName) + "/" + p.ServiceLogDir
}

// GetServiceDataPath returns the full service data path
func (p *PathsConfig) GetServiceDataPath(serviceName string) string {
	return p.GetServiceRoot(serviceName) + "/" + p.ServiceDataDir
}

// Clone creates a deep copy of PathsConfig
func (p *PathsConfig) Clone() *PathsConfig {
	return &PathsConfig{
		PluginInstallDir: p.PluginInstallDir,
		ServiceDeployDir: p.ServiceDeployDir,
		ServiceBinDir:    p.ServiceBinDir,
		ServiceConfigDir: p.ServiceConfigDir,
		ServiceLogDir:    p.ServiceLogDir,
		ServiceDataDir:   p.ServiceDataDir,
		BuildOutputDir:   p.BuildOutputDir,
		BuildCacheDir:    p.BuildCacheDir,
		ScriptDir:        p.ScriptDir,
		CIDir:            p.CIDir,
	}
}

// Validate validates the paths configuration
func (p *PathsConfig) Validate() error {
	if p.PluginInstallDir == "" {
		return errors.ErrInvalidPluginInstallDir
	}
	if p.ServiceDeployDir == "" {
		return errors.ErrInvalidServiceDeployDir
	}
	return nil
}
