package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "test-service", vars.ServiceName, "ServiceName should match")
	assert.Equal(t, 8080, vars.ServicePort, "ServicePort should be 8080")
	assert.Equal(t, "/usr/local/services/test-service", vars.ServiceRoot, "ServiceRoot should match")
	assert.Equal(t, "8080,9090", vars.PortsList, "PortsList should match")
	assert.Len(t, vars.Ports, 2, "Ports should have 2 items")
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

			assert.Equal(t, tt.wantGOARCH, archVars.GOARCH, "GOARCH should match")
			assert.Equal(t, tt.wantGOOS, archVars.GOOS, "GOOS should match")
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

	assert.Equal(t, "selfMonitor", pluginVars.PluginName, "PluginName should match")
	assert.Equal(t, "TCE Self Monitor", pluginVars.PluginDescription, "PluginDescription should match")
	assert.Equal(t, "https://example.com/download.sh", pluginVars.PluginDownloadURL, "PluginDownloadURL should match")
	assert.Equal(t, "/tce", pluginVars.PluginInstallDir, "PluginInstallDir should match")
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

	assert.Equal(t, "test-service", m["ServiceName"], "Map ServiceName should match")
	assert.Equal(t, "test-service", m["SERVICE_NAME"], "Map SERVICE_NAME should match")
	assert.Equal(t, 8080, m["ServicePort"], "Map ServicePort should be 8080")
	assert.Equal(t, 8080, m["SERVICE_PORT"], "Map SERVICE_PORT should be 8080")
}
