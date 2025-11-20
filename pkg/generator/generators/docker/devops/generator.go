package devops

import (
	_ "embed"
	"strings"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "devops"

// init registers the devops generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates TAD DevOps configuration
type Generator struct {
	core.BaseGenerator
}

// New creates a new devops generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates devops.yaml content in TAD format
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	vars := g.prepareTemplateVars()
	return g.RenderTemplate(template, vars)
}

// prepareTemplateVars prepares variables for devops template
func (g *Generator) prepareTemplateVars() map[string]interface{} {
	ctx := g.GetContext()

	// 创建镜像解析器
	resolver := ctx.Config.NewImageResolver()

	// 解析镜像
	builderImages := resolver.MustResolveBuilderImage()
	runtimeImages := resolver.MustResolveRuntimeImage()

	// Use preset for DevOps
	composer := ctx.GetVariablePreset().ForDevOps()

	// Parse runtime images
	runtimeImageX86, runtimeTagX86 := parseImageAndTag(runtimeImages.AMD64)
	runtimeImageARM, runtimeTagARM := parseImageAndTag(runtimeImages.ARM64)

	// Add DevOps-specific custom variables
	composer.
		WithCustom("RUNTIME_IMAGE_X86", runtimeImageX86).
		WithCustom("RUNTIME_TAG_X86", runtimeTagX86).
		WithCustom("RUNTIME_IMAGE_ARM", runtimeImageARM).
		WithCustom("RUNTIME_TAG_ARM", runtimeTagARM).
		WithCustom("BUILDER_IMAGE_X86", builderImages.AMD64).
		WithCustom("BUILDER_IMAGE_ARM", builderImages.ARM64).
		WithCustom("LANGUAGE_TYPE", ctx.Config.Language.Type).
		WithCustom("LANGUAGE_DISPLAY_NAME", getLanguageDisplayName(ctx.Config.Language.Type))

	return composer.Build()
}

// parseImageAndTag parses image name and tag from full image string
func parseImageAndTag(fullImage string) (string, string) {
	parts := strings.Split(fullImage, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return fullImage, "latest"
}

// getLanguageDisplayName returns display name for the language
func getLanguageDisplayName(langType string) string {
	switch langType {
	case "go":
		return "Go"
	case "python":
		return "Python"
	case "nodejs":
		return "Node.js"
	case "java":
		return "Java"
	case "rust":
		return "Rust"
	default:
		return langType
	}
}

//go:embed templates/devops.yaml.tmpl
var template string
