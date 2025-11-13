package config

import (
	"testing"
)

func TestLanguageConfig_GetString(t *testing.T) {
	tests := []struct {
		name         string
		config       map[string]interface{}
		key          string
		defaultValue string
		want         string
	}{
		{
			name: "get existing string value",
			config: map[string]interface{}{
				"goproxy": "https://goproxy.cn",
			},
			key:          "goproxy",
			defaultValue: "default",
			want:         "https://goproxy.cn",
		},
		{
			name:         "get non-existing key returns default",
			config:       map[string]interface{}{},
			key:          "missing",
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "nil config returns default",
			config:       nil,
			key:          "any",
			defaultValue: "default",
			want:         "default",
		},
		{
			name: "non-string value returns default",
			config: map[string]interface{}{
				"port": 8080,
			},
			key:          "port",
			defaultValue: "default",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LanguageConfig{
				Type:   "go",
				Config: tt.config,
			}
			got := lc.GetString(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguageConfig_GetInt(t *testing.T) {
	tests := []struct {
		name         string
		config       map[string]interface{}
		key          string
		defaultValue int
		want         int
	}{
		{
			name: "get existing int value",
			config: map[string]interface{}{
				"timeout": 30,
			},
			key:          "timeout",
			defaultValue: 10,
			want:         30,
		},
		{
			name: "get float64 value converts to int",
			config: map[string]interface{}{
				"timeout": 30.5,
			},
			key:          "timeout",
			defaultValue: 10,
			want:         30,
		},
		{
			name:         "get non-existing key returns default",
			config:       map[string]interface{}{},
			key:          "missing",
			defaultValue: 10,
			want:         10,
		},
		{
			name:         "nil config returns default",
			config:       nil,
			key:          "any",
			defaultValue: 10,
			want:         10,
		},
		{
			name: "non-numeric value returns default",
			config: map[string]interface{}{
				"timeout": "30",
			},
			key:          "timeout",
			defaultValue: 10,
			want:         10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LanguageConfig{
				Type:   "go",
				Config: tt.config,
			}
			got := lc.GetInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguageConfig_GetBool(t *testing.T) {
	tests := []struct {
		name         string
		config       map[string]interface{}
		key          string
		defaultValue bool
		want         bool
	}{
		{
			name: "get existing bool value true",
			config: map[string]interface{}{
				"enabled": true,
			},
			key:          "enabled",
			defaultValue: false,
			want:         true,
		},
		{
			name: "get existing bool value false",
			config: map[string]interface{}{
				"enabled": false,
			},
			key:          "enabled",
			defaultValue: true,
			want:         false,
		},
		{
			name:         "get non-existing key returns default",
			config:       map[string]interface{}{},
			key:          "missing",
			defaultValue: true,
			want:         true,
		},
		{
			name:         "nil config returns default",
			config:       nil,
			key:          "any",
			defaultValue: true,
			want:         true,
		},
		{
			name: "non-bool value returns default",
			config: map[string]interface{}{
				"enabled": "true",
			},
			key:          "enabled",
			defaultValue: false,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LanguageConfig{
				Type:   "go",
				Config: tt.config,
			}
			got := lc.GetBool(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguageConfig_GetStringSlice(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]interface{}
		key    string
		want   []string
	}{
		{
			name: "get existing string slice",
			config: map[string]interface{}{
				"tags": []string{"tag1", "tag2", "tag3"},
			},
			key:  "tags",
			want: []string{"tag1", "tag2", "tag3"},
		},
		{
			name: "get interface slice converts to string slice",
			config: map[string]interface{}{
				"tags": []interface{}{"tag1", "tag2", "tag3"},
			},
			key:  "tags",
			want: []string{"tag1", "tag2", "tag3"},
		},
		{
			name: "interface slice with non-string items filters them out",
			config: map[string]interface{}{
				"tags": []interface{}{"tag1", 123, "tag2"},
			},
			key:  "tags",
			want: []string{"tag1", "tag2"},
		},
		{
			name:   "get non-existing key returns nil",
			config: map[string]interface{}{},
			key:    "missing",
			want:   nil,
		},
		{
			name:   "nil config returns nil",
			config: nil,
			key:    "any",
			want:   nil,
		},
		{
			name: "non-slice value returns nil",
			config: map[string]interface{}{
				"tags": "tag1,tag2",
			},
			key:  "tags",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LanguageConfig{
				Type:   "go",
				Config: tt.config,
			}
			got := lc.GetStringSlice(tt.key)

			if (got == nil) != (tt.want == nil) {
				t.Errorf("GetStringSlice() = %v, want %v", got, tt.want)
				return
			}

			if got == nil {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetStringSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetStringSlice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestLanguageConfig_ComplexScenario(t *testing.T) {
	// Test a realistic Go configuration
	lc := &LanguageConfig{
		Type: "go",
		Config: map[string]interface{}{
			"goproxy":     "https://goproxy.cn,direct",
			"gosumdb":     "sum.golang.org",
			"goprivate":   []interface{}{"github.com/myorg/*", "gitlab.com/myteam/*"},
			"cgo_enabled": false,
			"build_tags":  []string{"integration", "e2e"},
			"timeout":     300,
		},
	}

	// Test string values
	if got := lc.GetString("goproxy", ""); got != "https://goproxy.cn,direct" {
		t.Errorf("goproxy = %v, want https://goproxy.cn,direct", got)
	}

	// Test bool values
	if got := lc.GetBool("cgo_enabled", true); got != false {
		t.Errorf("cgo_enabled = %v, want false", got)
	}

	// Test int values
	if got := lc.GetInt("timeout", 60); got != 300 {
		t.Errorf("timeout = %v, want 300", got)
	}

	// Test string slice
	goprivate := lc.GetStringSlice("goprivate")
	if len(goprivate) != 2 {
		t.Errorf("goprivate length = %v, want 2", len(goprivate))
	}

	// Test default values
	if got := lc.GetString("missing", "default"); got != "default" {
		t.Errorf("missing key = %v, want default", got)
	}
}
