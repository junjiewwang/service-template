package language

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
	"gopkg.in/yaml.v3"
)

type LanguageParserHandler struct {
	*chain.BaseHandler
}

func NewLanguageParserHandler() chain.ParserHandler {
	return &LanguageParserHandler{
		BaseHandler: chain.NewBaseHandler("language-parser"),
	}
}

func (h *LanguageParserHandler) Handle(ctx *chain.ProcessingContext) error {
	rawConfig := ctx.RawConfig
	langData, ok := rawConfig["language"]
	if !ok {
		return fmt.Errorf("language configuration not found")
	}

	yamlData, err := yaml.Marshal(langData)
	if err != nil {
		return fmt.Errorf("failed to marshal language config: %w", err)
	}

	var config LanguageConfig
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return fmt.Errorf("failed to unmarshal language config: %w", err)
	}

	ctx.SetDomainModel("language", &config)
	return h.CallNext(ctx)
}

func (h *LanguageParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

type LanguageValidatorHandler struct {
	*chain.BaseHandler
}

func NewLanguageValidatorHandler() chain.ValidatorHandler {
	return &LanguageValidatorHandler{
		BaseHandler: chain.NewBaseHandler("language-validator"),
	}
}

func (h *LanguageValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	model, ok := ctx.GetDomainModel("language")
	if !ok {
		return fmt.Errorf("language model not found")
	}

	config, ok := model.(*LanguageConfig)
	if !ok {
		return fmt.Errorf("invalid language model type")
	}

	if err := config.Validate(); err != nil {
		ctx.AddValidationError("language", err)
		return err
	}

	return h.CallNext(ctx)
}

func (h *LanguageValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

type LanguageGeneratorHandler struct {
	*chain.BaseHandler
}

func NewLanguageGeneratorHandler() chain.GeneratorHandler {
	return &LanguageGeneratorHandler{
		BaseHandler: chain.NewBaseHandler("language-generator"),
	}
}

func (h *LanguageGeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	model, ok := ctx.GetDomainModel("language")
	if !ok {
		return fmt.Errorf("language model not found")
	}

	config, ok := model.(*LanguageConfig)
	if !ok {
		return fmt.Errorf("invalid language model type")
	}

	// Generate language-specific configuration
	content := fmt.Sprintf("# Language Configuration\nType: %s\n", config.Type)
	ctx.AddGeneratedFile("language-config.txt", []byte(content))
	ctx.SetMetadata("language_type", config.Type)

	return h.CallNext(ctx)
}

func (h *LanguageGeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}
