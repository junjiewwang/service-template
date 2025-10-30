package generator

import (
	"strings"
	"testing"
)

func TestTemplateEngine_Render(t *testing.T) {
	engine := NewTemplateEngine()

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
			got, err := engine.Render(tt.template, tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
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
			got := SubstituteVariables(tt.text, tt.vars)
			if got != tt.want {
				t.Errorf("SubstituteVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateEngine_CustomFunctions(t *testing.T) {
	engine := NewTemplateEngine()

	t.Run("indentLines", func(t *testing.T) {
		template := `{{indentLines 2 "line1\nline2\nline3"}}`
		vars := map[string]interface{}{}

		got, err := engine.Render(template, vars)
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		lines := strings.Split(got, "\n")
		for i, line := range lines {
			if line != "" && !strings.HasPrefix(line, "  ") {
				t.Errorf("Line %d not indented: %q", i, line)
			}
		}
	})
}
