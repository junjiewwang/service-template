package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDownloadURLConfig_UnmarshalYAML_StaticString(t *testing.T) {
	yamlContent := `download_url: "https://example.com/plugin.tar.gz"`

	var config struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}

	err := yaml.Unmarshal([]byte(yamlContent), &config)
	require.NoError(t, err)

	assert.True(t, config.DownloadURL.IsStatic())
	assert.False(t, config.DownloadURL.IsArchMapping())
	assert.False(t, config.DownloadURL.IsEmpty())

	url, err := config.DownloadURL.GetStaticURL()
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/plugin.tar.gz", url)
}

func TestDownloadURLConfig_UnmarshalYAML_ArchMapping(t *testing.T) {
	yamlContent := `
download_url:
  x86_64: "https://example.com/plugin-x86_64.tar.gz"
  aarch64: "https://example.com/plugin-aarch64.tar.gz"
  default: "https://example.com/plugin-generic.tar.gz"
`

	var config struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}

	err := yaml.Unmarshal([]byte(yamlContent), &config)
	require.NoError(t, err)

	assert.False(t, config.DownloadURL.IsStatic())
	assert.True(t, config.DownloadURL.IsArchMapping())
	assert.False(t, config.DownloadURL.IsEmpty())

	urls, err := config.DownloadURL.GetArchURLs()
	require.NoError(t, err)
	assert.Len(t, urls, 3)
	assert.Equal(t, "https://example.com/plugin-x86_64.tar.gz", urls["x86_64"])
	assert.Equal(t, "https://example.com/plugin-aarch64.tar.gz", urls["aarch64"])
	assert.Equal(t, "https://example.com/plugin-generic.tar.gz", urls["default"])
}

func TestDownloadURLConfig_UnmarshalYAML_EmptyString(t *testing.T) {
	yamlContent := `download_url: ""`

	var config struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}

	err := yaml.Unmarshal([]byte(yamlContent), &config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty string")
}

func TestDownloadURLConfig_UnmarshalYAML_EmptyMap(t *testing.T) {
	yamlContent := `download_url: {}`

	var config struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}

	err := yaml.Unmarshal([]byte(yamlContent), &config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestDownloadURLConfig_MarshalYAML_StaticString(t *testing.T) {
	config := struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}{
		DownloadURL: NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
	}

	data, err := yaml.Marshal(&config)
	require.NoError(t, err)

	assert.Contains(t, string(data), "https://example.com/plugin.tar.gz")
}

func TestDownloadURLConfig_MarshalYAML_ArchMapping(t *testing.T) {
	config := struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}{
		DownloadURL: NewArchMappingDownloadURL(map[string]string{
			"x86_64":  "https://example.com/plugin-x86_64.tar.gz",
			"aarch64": "https://example.com/plugin-aarch64.tar.gz",
		}),
	}

	data, err := yaml.Marshal(&config)
	require.NoError(t, err)

	// Unmarshal back to verify
	var result struct {
		DownloadURL DownloadURLConfig `yaml:"download_url"`
	}
	err = yaml.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.True(t, result.DownloadURL.IsArchMapping())
	urls, err := result.DownloadURL.GetArchURLs()
	require.NoError(t, err)
	assert.Len(t, urls, 2)
}

func TestDownloadURLConfig_GetStaticURL_Error(t *testing.T) {
	config := NewArchMappingDownloadURL(map[string]string{
		"x86_64": "https://example.com/plugin.tar.gz",
	})

	_, err := config.GetStaticURL()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a static string")
}

func TestDownloadURLConfig_GetArchURLs_Error(t *testing.T) {
	config := NewStaticDownloadURL("https://example.com/plugin.tar.gz")

	_, err := config.GetArchURLs()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not an architecture mapping")
}

func TestDownloadURLConfig_String(t *testing.T) {
	tests := []struct {
		name     string
		config   DownloadURLConfig
		expected string
	}{
		{
			name:     "static URL",
			config:   NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
			expected: "https://example.com/plugin.tar.gz",
		},
		{
			name: "arch mapping",
			config: NewArchMappingDownloadURL(map[string]string{
				"x86_64":  "https://example.com/plugin-x86_64.tar.gz",
				"aarch64": "https://example.com/plugin-aarch64.tar.gz",
			}),
			expected: "arch_mapping(2 archs)",
		},
		{
			name:     "empty",
			config:   DownloadURLConfig{},
			expected: "<empty>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDownloadURLConfig_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		config   DownloadURLConfig
		expected bool
	}{
		{
			name:     "static URL is not empty",
			config:   NewStaticDownloadURL("https://example.com/plugin.tar.gz"),
			expected: false,
		},
		{
			name: "arch mapping is not empty",
			config: NewArchMappingDownloadURL(map[string]string{
				"x86_64": "https://example.com/plugin.tar.gz",
			}),
			expected: false,
		},
		{
			name:     "default config is empty",
			config:   DownloadURLConfig{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPluginConfig_WithDownloadURL(t *testing.T) {
	// Test complete plugin configuration with static URL
	yamlContent := `
name: selfMonitor
description: "TCE Self Monitor Tool"
download_url: "https://example.com/tool.sh"
install_command: |
  echo "Installing..."
required: true
`

	var plugin PluginConfig
	err := yaml.Unmarshal([]byte(yamlContent), &plugin)
	require.NoError(t, err)

	assert.Equal(t, "selfMonitor", plugin.Name)
	assert.Equal(t, "TCE Self Monitor Tool", plugin.Description)
	assert.True(t, plugin.DownloadURL.IsStatic())
	assert.True(t, plugin.Required)

	url, err := plugin.DownloadURL.GetStaticURL()
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/tool.sh", url)
}

func TestPluginConfig_WithArchMappingURL(t *testing.T) {
	// Test complete plugin configuration with arch mapping
	yamlContent := `
name: jre
description: "Java Runtime Environment"
download_url:
  x86_64: "https://example.com/jdk-x86_64.tar.gz"
  aarch64: "https://example.com/jdk-aarch64.tar.gz"
install_command: |
  echo "Installing JDK..."
required: false
`

	var plugin PluginConfig
	err := yaml.Unmarshal([]byte(yamlContent), &plugin)
	require.NoError(t, err)

	assert.Equal(t, "jre", plugin.Name)
	assert.Equal(t, "Java Runtime Environment", plugin.Description)
	assert.True(t, plugin.DownloadURL.IsArchMapping())
	assert.False(t, plugin.Required)

	urls, err := plugin.DownloadURL.GetArchURLs()
	require.NoError(t, err)
	assert.Len(t, urls, 2)
	assert.Equal(t, "https://example.com/jdk-x86_64.tar.gz", urls["x86_64"])
	assert.Equal(t, "https://example.com/jdk-aarch64.tar.gz", urls["aarch64"])
}
