package generator

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthcheckStrategyFactory_CreateStrategy(t *testing.T) {
	tests := []struct {
		name          string
		healthcheck   config.HealthcheckConfig
		expectedType  string
		expectError   bool
		errorContains string
	}{
		{
			name: "create default strategy when disabled",
			healthcheck: config.HealthcheckConfig{
				Enabled: false,
				Type:    "",
			},
			expectedType: "default",
			expectError:  false,
		},
		{
			name: "create default strategy when type is default",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
			expectedType: "default",
			expectError:  false,
		},
		{
			name: "create default strategy when type is empty",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "",
			},
			expectedType: "default",
			expectError:  false,
		},
		{
			name: "create custom strategy",
			healthcheck: config.HealthcheckConfig{
				Enabled:      true,
				Type:         "custom",
				CustomScript: "#!/bin/sh\nexit 0",
			},
			expectedType: "custom",
			expectError:  false,
		},
		{
			name: "error on unsupported type",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "http",
			},
			expectError:   true,
			errorContains: "unsupported healthcheck type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cfg := &config.ServiceConfig{
				Runtime: config.RuntimeConfig{
					Healthcheck: tt.healthcheck,
				},
			}
			factory := NewHealthcheckStrategyFactory(cfg)
			t.Logf("Testing healthcheck strategy factory with type: %s", tt.healthcheck.Type)

			// Act
			strategy, err := factory.CreateStrategy()

			// Assert
			if tt.expectError {
				require.Error(t, err, "Expected error but got none")
				assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				t.Logf("✓ Got expected error: %v", err)
			} else {
				require.NoError(t, err, "Should create strategy without error")
				require.NotNil(t, strategy, "Strategy should not be nil")
				assert.Equal(t, tt.expectedType, strategy.GetType(), "Strategy type should match expected")
				t.Logf("✓ Created strategy with type: %s", strategy.GetType())
			}
		})
	}
}

func TestDefaultHealthcheckStrategy(t *testing.T) {
	// Arrange
	cfg := &config.ServiceConfig{
		Service: config.ServiceInfo{
			Name:      "test-service",
			DeployDir: "/usr/local/services",
		},
		Runtime: config.RuntimeConfig{
			Healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
		},
	}
	strategy := NewDefaultHealthcheckStrategy(cfg)
	t.Log("Testing default healthcheck strategy")

	// Act - Test GetType
	strategyType := strategy.GetType()

	// Assert
	assert.Equal(t, "default", strategyType, "Strategy type should be 'default'")
	t.Logf("✓ Strategy type: %s", strategyType)

	// Act - Test Validate
	err := strategy.Validate()

	// Assert
	assert.NoError(t, err, "Default strategy validation should always pass")
	t.Log("✓ Validation passed")

	// Act - Test GenerateScript
	vars := map[string]interface{}{
		"SERVICE_NAME": cfg.Service.Name,
		"DEPLOY_DIR":   cfg.Service.DeployDir,
	}
	script, err := strategy.GenerateScript(vars)

	// Assert
	require.NoError(t, err, "Should generate script without error")
	assert.Contains(t, script, "#!/bin/sh", "Script should contain shebang")
	assert.Contains(t, script, "SERVICE_NAME", "Script should contain SERVICE_NAME variable")
	assert.Contains(t, script, "Default healthcheck", "Script should contain default healthcheck comment")
	assert.Contains(t, script, "ls -l /proc/*/exe", "Script should contain process check logic")
	t.Logf("✓ Generated script length: %d bytes", len(script))
}

func TestCustomHealthcheckStrategy(t *testing.T) {
	tests := []struct {
		name          string
		customScript  string
		expectError   bool
		errorContains string
	}{
		{
			name:         "valid custom script",
			customScript: "#!/bin/sh\ncurl -f http://localhost:8080/health || exit 1",
			expectError:  false,
		},
		{
			name:          "missing custom script",
			customScript:  "",
			expectError:   true,
			errorContains: "custom_script is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/usr/local/services",
				},
				Runtime: config.RuntimeConfig{
					Healthcheck: config.HealthcheckConfig{
						Enabled:      true,
						Type:         "custom",
						CustomScript: tt.customScript,
					},
				},
			}
			strategy := NewCustomHealthcheckStrategy(cfg)
			t.Logf("Testing custom healthcheck strategy with script length: %d", len(tt.customScript))

			// Act - Test GetType
			strategyType := strategy.GetType()

			// Assert
			assert.Equal(t, "custom", strategyType, "Strategy type should be 'custom'")
			t.Logf("✓ Strategy type: %s", strategyType)

			// Act - Test Validate
			err := strategy.Validate()

			// Assert
			if tt.expectError {
				require.Error(t, err, "Expected validation error")
				assert.Contains(t, err.Error(), tt.errorContains, "Error message should contain expected text")
				t.Logf("✓ Got expected validation error: %v", err)
				return
			}

			require.NoError(t, err, "Validation should pass")
			t.Log("✓ Validation passed")

			// Act - Test GenerateScript
			vars := map[string]interface{}{
				"SERVICE_NAME":  cfg.Service.Name,
				"DEPLOY_DIR":    cfg.Service.DeployDir,
				"CUSTOM_SCRIPT": tt.customScript,
			}
			script, err := strategy.GenerateScript(vars)

			// Assert
			require.NoError(t, err, "Should generate script without error")
			assert.Contains(t, script, "#!/bin/sh", "Script should contain shebang")
			assert.Contains(t, script, "SERVICE_NAME", "Script should contain SERVICE_NAME variable")
			assert.Contains(t, script, "Custom healthcheck", "Script should contain custom healthcheck comment")
			t.Logf("✓ Generated script length: %d bytes", len(script))
		})
	}
}

func TestHealthcheckScriptTemplateGenerator_WithStrategy(t *testing.T) {
	tests := []struct {
		name           string
		healthcheck    config.HealthcheckConfig
		expectedType   string
		expectError    bool
		scriptContains []string
	}{
		{
			name: "generate with default strategy",
			healthcheck: config.HealthcheckConfig{
				Enabled: true,
				Type:    "default",
			},
			expectedType: "default",
			expectError:  false,
			scriptContains: []string{
				"#!/bin/sh",
				"SERVICE_NAME",
				"Default healthcheck",
			},
		},
		{
			name: "generate with custom strategy",
			healthcheck: config.HealthcheckConfig{
				Enabled:      true,
				Type:         "custom",
				CustomScript: "#!/bin/sh\necho 'custom check'\nexit 0",
			},
			expectedType: "custom",
			expectError:  false,
			scriptContains: []string{
				"#!/bin/sh",
				"SERVICE_NAME",
				"Custom healthcheck",
			},
		},
		{
			name: "error on custom without script",
			healthcheck: config.HealthcheckConfig{
				Enabled:      true,
				Type:         "custom",
				CustomScript: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cfg := &config.ServiceConfig{
				Service: config.ServiceInfo{
					Name:      "test-service",
					DeployDir: "/usr/local/services",
				},
				Runtime: config.RuntimeConfig{
					Healthcheck: tt.healthcheck,
				},
			}
			engine := NewTemplateEngine()
			vars := NewVariables(cfg)
			t.Logf("Creating healthcheck generator with type: %s", tt.healthcheck.Type)

			// Act
			generator, err := NewHealthcheckScriptTemplateGenerator(cfg, engine, vars)

			// Assert
			if tt.expectError {
				require.Error(t, err, "Expected error during generator creation")
				t.Logf("✓ Got expected error: %v", err)
				return
			}

			require.NoError(t, err, "Should create generator without error")
			require.NotNil(t, generator, "Generator should not be nil")
			t.Log("✓ Generator created successfully")

			// Verify strategy type
			strategy := generator.GetStrategy()
			require.NotNil(t, strategy, "Strategy should not be nil")
			assert.Equal(t, tt.expectedType, strategy.GetType(), "Strategy type should match expected")
			t.Logf("✓ Strategy type: %s", strategy.GetType())

			// Act - Generate script
			script, err := generator.Generate()

			// Assert
			require.NoError(t, err, "Should generate script without error")
			assert.NotEmpty(t, script, "Generated script should not be empty")
			t.Logf("✓ Generated script length: %d bytes", len(script))

			// Verify script content
			for _, expected := range tt.scriptContains {
				assert.Contains(t, script, expected, "Script should contain expected content: %s", expected)
			}
			t.Logf("✓ Verified all %d expected content patterns", len(tt.scriptContains))
		})
	}
}
