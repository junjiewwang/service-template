package languageservice

import (
	"fmt"
	"strings"

	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// VariableSubstitutor decorates a strategy with variable substitution capability
// It replaces variables like ${BUILD_OUTPUT_DIR} in the install command
type VariableSubstitutor struct {
	*StrategyDecorator
	ctx *context.GeneratorContext
}

// NewVariableSubstitutor creates a new variable substitutor decorator
func NewVariableSubstitutor(strategy LanguageStrategy, ctx *context.GeneratorContext) *VariableSubstitutor {
	return &VariableSubstitutor{
		StrategyDecorator: NewStrategyDecorator(strategy),
		ctx:               ctx,
	}
}

// GetDepsInstallCommand returns the install command with variables substituted
func (s *VariableSubstitutor) GetDepsInstallCommand() string {
	// Get command from wrapped strategy
	command := s.wrapped.GetDepsInstallCommand()

	// Perform variable substitution
	return s.substituteVariables(command)
}

// substituteVariables replaces variables in the command string
// Supports ${VAR} format
func (s *VariableSubstitutor) substituteVariables(command string) string {
	if s.ctx == nil {
		return command
	}

	// Get all common variables from context
	composer := s.ctx.GetVariableComposer().WithCommon()
	variables := composer.Build()

	result := command

	// Replace ${VAR} format
	for key, value := range variables {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprint(value))
	}

	return result
}
