package languageservice

import (
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurableStrategy(t *testing.T) {
	tests := []struct {
		name           string
		language       string
		customCommand  string
		expectedResult string
	}{
		{
			name:           "python with custom command",
			language:       "python",
			customCommand:  "pip install -r requirements.txt -t /custom/dir",
			expectedResult: "pip install -r requirements.txt -t /custom/dir",
		},
		{
			name:           "python without custom command",
			language:       "python",
			customCommand:  "",
			expectedResult: "pip install -r requirements.txt",
		},
		{
			name:           "go with custom command",
			language:       "go",
			customCommand:  "GOMODCACHE=/cache go mod download",
			expectedResult: "GOMODCACHE=/cache go mod download",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create base strategy
			var baseStrategy LanguageStrategy
			switch tt.language {
			case "python":
				baseStrategy = NewPythonStrategy()
			case "go":
				baseStrategy = NewGoStrategy()
			}

			// Create config
			cfg := &config.LanguageConfig{
				Type: tt.language,
				Config: map[string]interface{}{
					"deps_install_command": tt.customCommand,
				},
			}

			// Decorate with configurable strategy
			decorated := NewConfigurableStrategy(baseStrategy, cfg)

			// Test
			result := decorated.GetDepsInstallCommand()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestVariableSubstitutor(t *testing.T) {
	tests := []struct {
		name           string
		language       string
		customCommand  string
		expectedResult string
	}{
		{
			name:           "python with BUILD_OUTPUT_DIR variable",
			language:       "python",
			customCommand:  "pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/bin",
			expectedResult: "pip install -r requirements.txt -t /opt/dist/bin",
		},
		{
			name:           "go with BUILD_OUTPUT_DIR variable",
			language:       "go",
			customCommand:  "GOMODCACHE=${BUILD_OUTPUT_DIR}/.cache go mod download",
			expectedResult: "GOMODCACHE=/opt/dist/.cache go mod download",
		},
		{
			name:           "multiple variables",
			language:       "python",
			customCommand:  "pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/bin --cache-dir ${BUILD_OUTPUT_DIR}/.cache",
			expectedResult: "pip install -r requirements.txt -t /opt/dist/bin --cache-dir /opt/dist/.cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			cfg := testutil.NewTestConfig()
			cfg.Language.Type = tt.language
			cfg.Language.Config = map[string]interface{}{
				"deps_install_command": tt.customCommand,
			}
			ctx := context.NewGeneratorContext(cfg, "/tmp/output")

			// Create base strategy
			var baseStrategy LanguageStrategy
			switch tt.language {
			case "python":
				baseStrategy = NewPythonStrategy()
			case "go":
				baseStrategy = NewGoStrategy()
			}

			// Decorate with configurable strategy
			configurable := NewConfigurableStrategy(baseStrategy, &cfg.Language)

			// Decorate with variable substitutor
			decorated := NewVariableSubstitutor(configurable, ctx)

			// Test
			result := decorated.GetDepsInstallCommand()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestStrategyFactory(t *testing.T) {
	tests := []struct {
		name           string
		language       string
		customCommand  string
		expectedResult string
	}{
		{
			name:           "python with custom command and variables",
			language:       "python",
			customCommand:  "pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/bin",
			expectedResult: "pip install -r requirements.txt -t /opt/dist/bin",
		},
		{
			name:           "go without custom command",
			language:       "go",
			customCommand:  "",
			expectedResult: "go mod download",
		},
		{
			name:           "nodejs with custom command",
			language:       "nodejs",
			customCommand:  "npm install --prefix ${BUILD_OUTPUT_DIR}",
			expectedResult: "npm install --prefix /opt/dist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			cfg := testutil.NewTestConfig()
			cfg.Language.Type = tt.language
			cfg.Language.Config = map[string]interface{}{
				"deps_install_command": tt.customCommand,
			}
			ctx := context.NewGeneratorContext(cfg, "/tmp/output")

			// Create factory
			factory := NewStrategyFactory(ctx)

			// Create strategy
			strategy, err := factory.CreateStrategy(tt.language, &cfg.Language)
			require.NoError(t, err)
			require.NotNil(t, strategy)

			// Test
			result := strategy.GetDepsInstallCommand()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestLanguageService_WithDecorators(t *testing.T) {
	tests := []struct {
		name           string
		language       string
		customCommand  string
		expectedResult string
	}{
		{
			name:           "python with custom command and variables",
			language:       "python",
			customCommand:  "pip install -r requirements.txt -t ${BUILD_OUTPUT_DIR}/bin",
			expectedResult: "pip install -r requirements.txt -t /opt/dist/bin",
		},
		{
			name:           "go without custom command",
			language:       "go",
			customCommand:  "",
			expectedResult: "go mod download",
		},
		{
			name:           "nodejs with SERVICE_NAME variable",
			language:       "nodejs",
			customCommand:  "npm install --prefix ${BUILD_OUTPUT_DIR}/${SERVICE_NAME}",
			expectedResult: "npm install --prefix /opt/dist/test-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			cfg := testutil.NewTestConfig()
			cfg.Language.Type = tt.language
			cfg.Language.Config = map[string]interface{}{
				"deps_install_command": tt.customCommand,
			}
			ctx := context.NewGeneratorContext(cfg, "/tmp/output")

			// Create language service
			service := NewLanguageService(ctx)

			// Test
			result := service.GetDepsInstallCommand(tt.language)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestLanguageService_UnsupportedLanguage(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	service := NewLanguageService(ctx)

	result := service.GetDepsInstallCommand("unsupported")
	assert.Equal(t, "echo 'No dependency installation needed'", result)
}

func TestStrategyDecorator_Unwrap(t *testing.T) {
	baseStrategy := NewPythonStrategy()
	decorator := NewStrategyDecorator(baseStrategy)

	unwrapped := decorator.Unwrap()
	assert.Equal(t, baseStrategy, unwrapped)
}
