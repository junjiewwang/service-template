package errors

import (
	"fmt"
)

// ErrorCode represents error type code
type ErrorCode string

const (
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeTemplate      ErrorCode = "TEMPLATE_ERROR"
	ErrCodePlugin        ErrorCode = "PLUGIN_ERROR"
	ErrCodeLanguage      ErrorCode = "LANGUAGE_ERROR"
	ErrCodeFileSystem    ErrorCode = "FILESYSTEM_ERROR"
	ErrCodeConfiguration ErrorCode = "CONFIGURATION_ERROR"
)

// GeneratorError represents a structured error in generator
type GeneratorError struct {
	Code      ErrorCode
	Message   string
	Generator string
	Cause     error
	Context   map[string]interface{}
}

// Error implements error interface
func (e *GeneratorError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Code, e.Generator, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Generator, e.Message)
}

// Unwrap implements errors.Unwrap
func (e *GeneratorError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to error
func (e *GeneratorError) WithContext(key string, value interface{}) *GeneratorError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewValidationError creates a validation error
func NewValidationError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodeValidation,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// NewTemplateError creates a template error
func NewTemplateError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodeTemplate,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// NewPluginError creates a plugin error
func NewPluginError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodePlugin,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// NewLanguageError creates a language error
func NewLanguageError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodeLanguage,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// NewFileSystemError creates a filesystem error
func NewFileSystemError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodeFileSystem,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// NewConfigurationError creates a configuration error
func NewConfigurationError(generator, message string, cause error) *GeneratorError {
	return &GeneratorError{
		Code:      ErrCodeConfiguration,
		Generator: generator,
		Message:   message,
		Cause:     cause,
	}
}

// Path configuration errors
var (
	ErrInvalidPluginInstallDir = &GeneratorError{Code: ErrCodeConfiguration, Message: "plugin install directory cannot be empty"}
	ErrInvalidServiceDeployDir = &GeneratorError{Code: ErrCodeConfiguration, Message: "service deploy directory cannot be empty"}
	ErrInvalidPathsConfig      = &GeneratorError{Code: ErrCodeConfiguration, Message: "invalid paths configuration"}
)
