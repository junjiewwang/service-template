package models

import (
	"testing"
)

func TestNewPathsConfig(t *testing.T) {
	config := NewPathsConfig()

	if config.PluginInstallDir != "/plugins" {
		t.Errorf("PluginInstallDir = %v, want /plugins", config.PluginInstallDir)
	}

	if config.ServiceDeployDir != "/data/services" {
		t.Errorf("ServiceDeployDir = %v, want /data/services", config.ServiceDeployDir)
	}

	if config.ServiceBinDir != "bin" {
		t.Errorf("ServiceBinDir = %v, want bin", config.ServiceBinDir)
	}
}

func TestPathsConfig_WithMethods(t *testing.T) {
	config := NewPathsConfig().
		WithPluginInstallDir("/custom/plugins").
		WithServiceDeployDir("/custom/services").
		WithServiceBinDir("binaries")

	if config.PluginInstallDir != "/custom/plugins" {
		t.Errorf("PluginInstallDir = %v, want /custom/plugins", config.PluginInstallDir)
	}

	if config.ServiceDeployDir != "/custom/services" {
		t.Errorf("ServiceDeployDir = %v, want /custom/services", config.ServiceDeployDir)
	}

	if config.ServiceBinDir != "binaries" {
		t.Errorf("ServiceBinDir = %v, want binaries", config.ServiceBinDir)
	}
}

func TestPathsConfig_GetServiceRoot(t *testing.T) {
	config := NewPathsConfig()

	root := config.GetServiceRoot("myservice")
	expected := "/data/services/myservice"

	if root != expected {
		t.Errorf("GetServiceRoot() = %v, want %v", root, expected)
	}
}

func TestPathsConfig_GetServiceBinPath(t *testing.T) {
	config := NewPathsConfig()

	binPath := config.GetServiceBinPath("myservice")
	expected := "/data/services/myservice/bin"

	if binPath != expected {
		t.Errorf("GetServiceBinPath() = %v, want %v", binPath, expected)
	}
}

func TestPathsConfig_GetServiceConfigPath(t *testing.T) {
	config := NewPathsConfig()

	configPath := config.GetServiceConfigPath("myservice")
	expected := "/data/services/myservice/conf"

	if configPath != expected {
		t.Errorf("GetServiceConfigPath() = %v, want %v", configPath, expected)
	}
}

func TestPathsConfig_GetServiceLogPath(t *testing.T) {
	config := NewPathsConfig()

	logPath := config.GetServiceLogPath("myservice")
	expected := "/data/services/myservice/logs"

	if logPath != expected {
		t.Errorf("GetServiceLogPath() = %v, want %v", logPath, expected)
	}
}

func TestPathsConfig_GetServiceDataPath(t *testing.T) {
	config := NewPathsConfig()

	dataPath := config.GetServiceDataPath("myservice")
	expected := "/data/services/myservice/data"

	if dataPath != expected {
		t.Errorf("GetServiceDataPath() = %v, want %v", dataPath, expected)
	}
}

func TestPathsConfig_Clone(t *testing.T) {
	original := NewPathsConfig().
		WithPluginInstallDir("/custom/plugins").
		WithServiceDeployDir("/custom/services")

	cloned := original.Clone()

	// Verify values are copied
	if cloned.PluginInstallDir != original.PluginInstallDir {
		t.Error("Clone did not copy PluginInstallDir")
	}

	if cloned.ServiceDeployDir != original.ServiceDeployDir {
		t.Error("Clone did not copy ServiceDeployDir")
	}

	// Verify it's a deep copy
	cloned.PluginInstallDir = "/different/path"
	if original.PluginInstallDir == cloned.PluginInstallDir {
		t.Error("Clone is not a deep copy")
	}
}

func TestPathsConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *PathsConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  NewPathsConfig(),
			wantErr: false,
		},
		{
			name: "empty plugin install dir",
			config: &PathsConfig{
				PluginInstallDir: "",
				ServiceDeployDir: "/data/services",
			},
			wantErr: true,
		},
		{
			name: "empty service deploy dir",
			config: &PathsConfig{
				PluginInstallDir: "/plugins",
				ServiceDeployDir: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPathsConfig_CustomPaths(t *testing.T) {
	config := NewPathsConfig().
		WithServiceDeployDir("/opt/apps").
		WithServiceBinDir("executable").
		WithServiceConfigDir("config").
		WithServiceLogDir("log").
		WithServiceDataDir("storage")

	serviceName := "testservice"

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"root", config.GetServiceRoot(serviceName), "/opt/apps/testservice"},
		{"bin", config.GetServiceBinPath(serviceName), "/opt/apps/testservice/executable"},
		{"config", config.GetServiceConfigPath(serviceName), "/opt/apps/testservice/config"},
		{"log", config.GetServiceLogPath(serviceName), "/opt/apps/testservice/log"},
		{"data", config.GetServiceDataPath(serviceName), "/opt/apps/testservice/storage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %v, want %v", tt.got, tt.expected)
			}
		})
	}
}
