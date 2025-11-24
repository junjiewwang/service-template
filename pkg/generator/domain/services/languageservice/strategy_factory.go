package languageservice

import (
	"fmt"
	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// strategyConstructor defines a factory function that creates a base LanguageStrategy
type strategyConstructor func() LanguageStrategy

// LanguageRegistryEntry holds metadata about a supported language
type LanguageRegistryEntry struct {
	// Name is the language identifier (e.g., "go", "python")
	Name string
	// Constructor creates a new instance of the base strategy
	Constructor strategyConstructor
	// Description provides human-readable information about the language
	Description string
}

// strategyRegistry is the single source of truth for all supported languages.
// To add a new language:
// 1. Define a constant in language_service.go (e.g., LangNewLang = "newlang")
// 2. Implement a strategy struct that satisfies LanguageStrategy interface
// 3. Add an entry to this registry
var strategyRegistry = map[string]LanguageRegistryEntry{
	LangGo: {
		Name:        LangGo,
		Constructor: func() LanguageStrategy { return NewGoStrategy() },
		Description: "Go programming language with go.mod support",
	},
	LangPython: {
		Name:        LangPython,
		Constructor: func() LanguageStrategy { return NewPythonStrategy() },
		Description: "Python with pip and requirements.txt",
	},
	LangNodeJS: {
		Name:        LangNodeJS,
		Constructor: func() LanguageStrategy { return NewNodeJSStrategy() },
		Description: "Node.js with npm and package.json",
	},
	LangJava: {
		Name:        LangJava,
		Constructor: func() LanguageStrategy { return NewJavaStrategy() },
		Description: "Java with Maven (pom.xml) or Gradle (build.gradle)",
	},
	LangRust: {
		Name:        LangRust,
		Constructor: func() LanguageStrategy { return NewRustStrategy() },
		Description: "Rust with Cargo.toml",
	},
}

// StrategyFactory creates decorated language strategies
type StrategyFactory struct {
	ctx *context.GeneratorContext
}

// NewStrategyFactory creates a new strategy factory
func NewStrategyFactory(ctx *context.GeneratorContext) *StrategyFactory {
	return &StrategyFactory{
		ctx: ctx,
	}
}

// CreateStrategy creates a fully decorated strategy for the given language
// Decoration chain: BaseStrategy -> ConfigurableStrategy -> VariableSubstitutor
func (f *StrategyFactory) CreateStrategy(language string, config *config.LanguageConfig) (LanguageStrategy, error) {
	// 1. Look up base strategy from registry
	entry, ok := strategyRegistry[language]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	baseStrategy := entry.Constructor()

	// 2. Decorate with configurable strategy (custom command support)
	decorated := NewConfigurableStrategy(baseStrategy, config)

	// 3. Decorate with variable substitutor (variable replacement)
	result := NewVariableSubstitutor(decorated, f.ctx)

	return result, nil
}
