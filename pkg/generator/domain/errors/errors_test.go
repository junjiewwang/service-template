package errors

import (
	"errors"
	"testing"
)

func TestGeneratorError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *GeneratorError
		expected string
	}{
		{
			name: "error without cause",
			err: &GeneratorError{
				Code:      ErrCodeValidation,
				Generator: "TestGenerator",
				Message:   "validation failed",
			},
			expected: "[VALIDATION_ERROR] TestGenerator: validation failed",
		},
		{
			name: "error with cause",
			err: &GeneratorError{
				Code:      ErrCodeTemplate,
				Generator: "TestGenerator",
				Message:   "template rendering failed",
				Cause:     errors.New("template not found"),
			},
			expected: "[TEMPLATE_ERROR] TestGenerator: template rendering failed (caused by: template not found)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGeneratorError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := &GeneratorError{
		Code:      ErrCodePlugin,
		Generator: "TestGenerator",
		Message:   "plugin error",
		Cause:     cause,
	}

	if unwrapped := err.Unwrap(); unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestGeneratorError_WithContext(t *testing.T) {
	err := &GeneratorError{
		Code:      ErrCodeConfiguration,
		Generator: "TestGenerator",
		Message:   "config error",
	}

	err.WithContext("key1", "value1").WithContext("key2", 123)

	if len(err.Context) != 2 {
		t.Errorf("Context length = %d, want 2", len(err.Context))
	}

	if err.Context["key1"] != "value1" {
		t.Errorf("Context[key1] = %v, want value1", err.Context["key1"])
	}

	if err.Context["key2"] != 123 {
		t.Errorf("Context[key2] = %v, want 123", err.Context["key2"])
	}
}

func TestNewValidationError(t *testing.T) {
	cause := errors.New("invalid input")
	err := NewValidationError("TestGenerator", "validation failed", cause)

	if err.Code != ErrCodeValidation {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeValidation)
	}

	if err.Generator != "TestGenerator" {
		t.Errorf("Generator = %v, want TestGenerator", err.Generator)
	}

	if err.Message != "validation failed" {
		t.Errorf("Message = %v, want validation failed", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestNewTemplateError(t *testing.T) {
	err := NewTemplateError("TestGenerator", "template error", nil)

	if err.Code != ErrCodeTemplate {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeTemplate)
	}
}

func TestNewPluginError(t *testing.T) {
	err := NewPluginError("TestGenerator", "plugin error", nil)

	if err.Code != ErrCodePlugin {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodePlugin)
	}
}

func TestNewLanguageError(t *testing.T) {
	err := NewLanguageError("TestGenerator", "language error", nil)

	if err.Code != ErrCodeLanguage {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeLanguage)
	}
}

func TestNewFileSystemError(t *testing.T) {
	err := NewFileSystemError("TestGenerator", "filesystem error", nil)

	if err.Code != ErrCodeFileSystem {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeFileSystem)
	}
}

func TestNewConfigurationError(t *testing.T) {
	err := NewConfigurationError("TestGenerator", "config error", nil)

	if err.Code != ErrCodeConfiguration {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeConfiguration)
	}
}
