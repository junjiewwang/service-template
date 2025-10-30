package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
)

func TestNewVariables(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name: "test-service",
			Ports: []config.PortConfig{
				{Name: "http", Port: 8080, Protocol: "TCP"},
				{Name: "metrics", Port: 9090, Protocol: "TCP"},
			},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{
			Type:    "go",
			Version: "1.23",
		},
		Build: config.BuildConfig{
			OutputDir: "dist",
		},
	}

	vars := NewVariables(cfg)

	if vars.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service", vars.ServiceName)
	}

	if vars.ServicePort != 8080 {
		t.Errorf("ServicePort = %v, want 8080", vars.ServicePort)
	}

	if vars.ServiceRoot != "/usr/local/services/test-service" {
		t.Errorf("ServiceRoot = %v, want /usr/local/services/test-service", vars.ServiceRoot)
	}

	if vars.PortsList != "8080,9090" {
		t.Errorf("PortsList = %v, want 8080,9090", vars.PortsList)
	}

	if len(vars.Ports) != 2 {
		t.Errorf("Ports length = %v, want 2", len(vars.Ports))
	}
}

func TestVariables_WithArchitecture(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "go", Version: "1.23"},
		Build:    config.BuildConfig{OutputDir: "dist"},
	}

	vars := NewVariables(cfg)

	tests := []struct {
		arch       string
		wantGOARCH string
		wantGOOS   string
	}{
		{"amd64", "amd64", "linux"},
		{"arm64", "arm64", "linux"},
	}

	for _, tt := range tests {
		t.Run(tt.arch, func(t *testing.T) {
			archVars := vars.WithArchitecture(tt.arch)

			if archVars.GOARCH != tt.wantGOARCH {
				t.Errorf("GOARCH = %v, want %v", archVars.GOARCH, tt.wantGOARCH)
			}

			if archVars.GOOS != tt.wantGOOS {
				t.Errorf("GOOS = %v, want %v", archVars.GOOS, tt.wantGOOS)
			}
		})
	}
}

func TestVariables_WithPlugin(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "go", Version: "1.23"},
		Build:    config.BuildConfig{OutputDir: "dist"},
	}

	vars := NewVariables(cfg)

	plugin := config.PluginConfig{
		Name:        "selfMonitor",
		Description: "TCE Self Monitor",
		DownloadURL: "https://example.com/download.sh",
		InstallDir:  "/tce",
	}

	pluginVars := vars.WithPlugin(plugin)

	if pluginVars.PluginName != "selfMonitor" {
		t.Errorf("PluginName = %v, want selfMonitor", pluginVars.PluginName)
	}

	if pluginVars.PluginDescription != "TCE Self Monitor" {
		t.Errorf("PluginDescription = %v, want TCE Self Monitor", pluginVars.PluginDescription)
	}

	if pluginVars.PluginDownloadURL != "https://example.com/download.sh" {
		t.Errorf("PluginDownloadURL = %v, want https://example.com/download.sh", pluginVars.PluginDownloadURL)
	}

	if pluginVars.PluginInstallDir != "/tce" {
		t.Errorf("PluginInstallDir = %v, want /tce", pluginVars.PluginInstallDir)
	}
}

func TestVariables_ToMap(t *testing.T) {
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "go", Version: "1.23"},
		Build:    config.BuildConfig{OutputDir: "dist"},
	}

	vars := NewVariables(cfg)
	m := vars.ToMap()

	if m["ServiceName"] != "test-service" {
		t.Errorf("Map ServiceName = %v, want test-service", m["ServiceName"])
	}

	if m["SERVICE_NAME"] != "test-service" {
		t.Errorf("Map SERVICE_NAME = %v, want test-service", m["SERVICE_NAME"])
	}

	if m["ServicePort"] != 8080 {
		t.Errorf("Map ServicePort = %v, want 8080", m["ServicePort"])
	}

	if m["SERVICE_PORT"] != 8080 {
		t.Errorf("Map SERVICE_PORT = %v, want 8080", m["SERVICE_PORT"])
	}
}
