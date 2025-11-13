package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/stretchr/testify/assert"
)

func TestNewVariables(t *testing.T) {
	// Arrange: Setup service configuration
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
			Type:    "golang",
		},
		Build: config.BuildConfig{},
	}

	// Act: Create variables from configuration
	paths := context.NewPaths(cfg)
	vars := context.NewVariables(cfg, paths)

	// Assert: Verify variables are correctly set
	assert.Equal(t, "test-service", vars.ServiceName, "ServiceName should match")
	assert.Equal(t, 8080, vars.ServicePort, "ServicePort should be 8080")
	assert.Equal(t, "/usr/local/services/test-service", vars.ServiceRoot, "ServiceRoot should match")
	assert.Equal(t, "8080,9090", vars.PortsList, "PortsList should match")
	assert.Len(t, vars.Ports, 2, "Ports should have 2 items")

	t.Logf("✓ Variables created: ServiceName=%s, ServicePort=%d, PortsList=%s",
		vars.ServiceName, vars.ServicePort, vars.PortsList)
}

func TestVariables_WithArchitecture(t *testing.T) {
	// Arrange: Setup service configuration
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "golang"},
	}

	paths := context.NewPaths(cfg)
	vars := context.NewVariables(cfg, paths)

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
			// Act: Create architecture-specific variables
			archVars := vars.WithArchitecture(tt.arch)

			// Assert: Verify architecture variables
			assert.Equal(t, tt.wantGOARCH, archVars.GOARCH, "GOARCH should match")
			assert.Equal(t, tt.wantGOOS, archVars.GOOS, "GOOS should match")
			t.Logf("✓ Architecture %s: GOARCH=%s, GOOS=%s", tt.arch, archVars.GOARCH, archVars.GOOS)
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
		Language: config.LanguageConfig{Type: "go"},
	}

	paths := context.NewPaths(cfg)
	vars := context.NewVariables(cfg, paths)

	plugin := config.PluginConfig{
		Name:        "selfMonitor",
		Description: "TCE Self Monitor",
		DownloadURL: "https://example.com/download.sh",
	}

	installDir := "/tce"
	pluginVars := vars.WithPlugin(plugin, installDir)

	assert.Equal(t, "selfMonitor", pluginVars.PluginName, "PluginName should match")
	assert.Equal(t, "TCE Self Monitor", pluginVars.PluginDescription, "PluginDescription should match")
	assert.Equal(t, "https://example.com/download.sh", pluginVars.PluginDownloadURL, "PluginDownloadURL should match")
	assert.Equal(t, "/tce", pluginVars.PluginInstallDir, "PluginInstallDir should match")
}

func TestVariables_ToMap(t *testing.T) {
	// Arrange: Setup service configuration
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			Ports:     []config.PortConfig{{Name: "http", Port: 8080, Protocol: "TCP"}},
			DeployDir: "/usr/local/services",
		},
		Language: config.LanguageConfig{Type: "golang"},
	}

	paths := context.NewPaths(cfg)
	vars := context.NewVariables(cfg, paths)

	// Act: Convert variables to map
	m := vars.ToMap()

	// Assert: Verify map contains expected keys and values
	assert.Equal(t, "test-service", m["ServiceName"], "Map ServiceName should match")
	assert.Equal(t, "test-service", m["SERVICE_NAME"], "Map SERVICE_NAME should match")
	assert.Equal(t, 8080, m["ServicePort"], "Map ServicePort should be 8080")
	assert.Equal(t, 8080, m["SERVICE_PORT"], "Map SERVICE_PORT should be 8080")

	t.Logf("✓ Variables map contains %d keys", len(m))
	t.Logf("✓ Verified both camelCase and UPPER_CASE keys present")
}
