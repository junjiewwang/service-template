package generator

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateEngine_Render(t *testing.T) {
	// Arrange: Create template engine
	engine := core.NewTemplateEngine()
	require.NotNil(t, engine, "Template engine should be created")

	tests := []struct {
		name     string
		template string
		vars     map[string]interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple substitution",
			template: "Hello {{.Name}}!",
			vars:     map[string]interface{}{"Name": "World"},
			want:     "Hello World!",
			wantErr:  false,
		},
		{
			name:     "with sprig function",
			template: "{{.Name | upper}}",
			vars:     map[string]interface{}{"Name": "hello"},
			want:     "HELLO",
			wantErr:  false,
		},
		{
			name:     "with custom substitute function",
			template: `{{substitute "Port: ${PORT}" .}}`,
			vars:     map[string]interface{}{"PORT": 8080},
			want:     "Port: 8080",
			wantErr:  false,
		},
		{
			name:     "with conditional",
			template: "{{if .Enabled}}Enabled{{else}}Disabled{{end}}",
			vars:     map[string]interface{}{"Enabled": true},
			want:     "Enabled",
			wantErr:  false,
		},
		{
			name:     "with range",
			template: "{{range .Items}}{{.}},{{end}}",
			vars:     map[string]interface{}{"Items": []string{"a", "b", "c"}},
			want:     "a,b,c,",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			template: "{{.Name",
			vars:     map[string]interface{}{"Name": "World"},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Render template
			got, err := engine.Render(tt.template, tt.vars)

			// Assert: Check results
			if tt.wantErr {
				assert.Error(t, err, "Render() should return an error for invalid template")
				t.Logf("Expected error occurred: %v", err)
			} else {
				require.NoError(t, err, "Render() should not return an error")
				assert.Equal(t, tt.want, got, "Rendered content should match expected")
				t.Logf("✓ Template rendered successfully: %q", got)
			}
		})
	}
}

func TestSubstituteVariables(t *testing.T) {
	tests := []struct {
		name string
		text string
		vars map[string]interface{}
		want string
	}{
		{
			name: "single variable",
			text: "Port: ${PORT}",
			vars: map[string]interface{}{"PORT": 8080},
			want: "Port: 8080",
		},
		{
			name: "multiple variables",
			text: "${SERVICE_NAME} on ${SERVICE_PORT}",
			vars: map[string]interface{}{
				"SERVICE_NAME": "my-service",
				"SERVICE_PORT": 8080,
			},
			want: "my-service on 8080",
		},
		{
			name: "no variables",
			text: "Hello World",
			vars: map[string]interface{}{},
			want: "Hello World",
		},
		{
			name: "variable not found",
			text: "Port: ${PORT}",
			vars: map[string]interface{}{"OTHER": 8080},
			want: "Port: ${PORT}",
		},
		{
			name: "multiline text",
			text: `#!/bin/sh
cd ${SERVICE_ROOT}
exec ${SERVICE_BIN_DIR}/${SERVICE_NAME}`,
			vars: map[string]interface{}{
				"SERVICE_ROOT":    "/app",
				"SERVICE_BIN_DIR": "/app/bin",
				"SERVICE_NAME":    "myapp",
			},
			want: `#!/bin/sh
cd /app
exec /app/bin/myapp`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Substitute variables
			got := core.SubstituteVariables(tt.text, tt.vars)

			// Assert: Check result
			assert.Equal(t, tt.want, got, "SubstituteVariables() should return expected result")
			t.Logf("✓ Variable substitution successful: %q -> %q", tt.text, got)
		})
	}
}

func TestTemplateEngine_CustomFunctions(t *testing.T) {
	t.Run("indentLines", func(t *testing.T) {
		// Arrange: Create template engine and setup template with indentLines function
		template := `{{indentLines 2 "line1\nline2\nline3"}}`
		vars := map[string]interface{}{}
		engine := core.NewTemplateEngine()
		require.NotNil(t, engine, "Template engine should be created")

		// Act: Render template
		got, err := engine.Render(template, vars)
		require.NoError(t, err, "Render() should not return an error")

		// Assert: Check that all lines are indented
		lines := strings.Split(got, "\n")
		t.Logf("Rendered %d lines", len(lines))
		for i, line := range lines {
			if line != "" {
				assert.True(t, strings.HasPrefix(line, "  "),
					"Line %d should be indented with 2 spaces: %q", i, line)
				t.Logf("✓ Line %d correctly indented: %q", i, line)
			}
		}
	})
}
