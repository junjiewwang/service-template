package core

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// TemplateEngine handles template rendering
type TemplateEngine struct {
	funcMap template.FuncMap
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		funcMap: sprig.TxtFuncMap(),
	}

	// Add custom functions
	engine.addCustomFunctions()

	return engine
}

// addCustomFunctions adds custom template functions
func (e *TemplateEngine) addCustomFunctions() {
	// Add variable substitution function
	e.funcMap["substitute"] = func(text string, vars map[string]interface{}) string {
		result := text
		for key, value := range vars {
			placeholder := fmt.Sprintf("${%s}", key)
			result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
		}
		return result
	}

	// Add join function for ports
	e.funcMap["joinPorts"] = func(ports []interface{}, sep string) string {
		var result []string
		for _, p := range ports {
			result = append(result, fmt.Sprint(p))
		}
		return strings.Join(result, sep)
	}

	// Add indent function
	e.funcMap["indentLines"] = func(spaces int, text string) string {
		indent := strings.Repeat(" ", spaces)
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if line != "" {
				lines[i] = indent + line
			}
		}
		return strings.Join(lines, "\n")
	}
}

// Render renders a template with the given variables
func (e *TemplateEngine) Render(templateContent string, vars map[string]interface{}) (string, error) {
	tmpl, err := template.New("template").Funcs(e.funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderWithName renders a named template
func (e *TemplateEngine) RenderWithName(name, templateContent string, vars map[string]interface{}) (string, error) {
	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// SubstituteVariables performs simple variable substitution in text
func SubstituteVariables(text string, vars map[string]interface{}) string {
	result := text
	for key, value := range vars {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}
	return result
}

// ReplaceVariables replaces variables in text with string values
func (e *TemplateEngine) ReplaceVariables(text string, vars map[string]string) string {
	result := text
	for key, value := range vars {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
