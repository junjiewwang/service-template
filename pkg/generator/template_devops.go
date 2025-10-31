package generator

import (
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeDevOps = "devops"

// DevOpsTemplateGenerator generates TAD DevOps configuration using factory pattern
type DevOpsTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the DevOps generator
func init() {
	RegisterGenerator(GeneratorTypeDevOps, createDevOpsGenerator)
}

// createDevOpsGenerator is the creator function for DevOps generator
func createDevOpsGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewDevOpsTemplateGenerator(cfg, engine, vars), nil
}

// NewDevOpsTemplateGenerator creates a new DevOps template generator
func NewDevOpsTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DevOpsTemplateGenerator {
	return &DevOpsTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeDevOps,
		},
	}
}

// Generate generates devops.yaml content in TAD format
func (g *DevOpsTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for devops template
func (g *DevOpsTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt

	// Parse runtime images
	runtimeImageX86, runtimeTagX86 := parseImageAndTag(g.config.Build.RuntimeImage.AMD64)
	runtimeImageARM, runtimeTagARM := parseImageAndTag(g.config.Build.RuntimeImage.ARM64)

	vars["RUNTIME_IMAGE_X86"] = runtimeImageX86
	vars["RUNTIME_TAG_X86"] = runtimeTagX86
	vars["RUNTIME_IMAGE_ARM"] = runtimeImageARM
	vars["RUNTIME_TAG_ARM"] = runtimeTagARM

	// Builder images
	vars["BUILDER_IMAGE_X86"] = g.config.Build.BuilderImage.AMD64
	vars["BUILDER_IMAGE_ARM"] = g.config.Build.BuilderImage.ARM64

	// Language info
	vars["LANGUAGE_TYPE"] = g.config.Language.Type
	vars["LANGUAGE_VERSION"] = g.config.Language.Version
	vars["LANGUAGE_DISPLAY_NAME"] = getLanguageDisplayName(g.config.Language.Type, g.config.Language.Version)

	// Deployment config
	vars["DEPLOY_DIR"] = g.config.Service.DeployDir

	return vars
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
func getLanguageDisplayName(langType, version string) string {
	switch langType {
	case "go":
		return "Go " + version
	case "python":
		return "Python " + version
	case "nodejs":
		return "Node.js " + version
	case "java":
		return "Java " + version
	default:
		return langType
	}
}

func (g *DevOpsTemplateGenerator) getTemplate() string {

	return `# Auto-generated DevOps configuration
# Generated at: {{ .GENERATED_AT }}

tad:
  export_envs:
  # ============================================
  # Runtime Base Images (运行时基础镜像)
  # ============================================
  - name: "TLINUX_BASE_IMAGE_X86"  # 运行时基础镜像名称 (x86)
    value: "{{ .RUNTIME_IMAGE_X86 }}"
  - name: "TLINUX_TAG_X86"  # 运行时基础镜像tag (x86)
    value: "{{ .RUNTIME_TAG_X86 }}"
  - name: "TLINUX_BASE_IMAGE_ARM"  # 运行时基础镜像名称 (arm64)
    value: "{{ .RUNTIME_IMAGE_ARM }}"
  - name: "TLINUX_TAG_ARM"  # 运行时基础镜像tag (arm64)
    value: "{{ .RUNTIME_TAG_ARM }}"
  
  # ============================================
  # Builder Images (构建镜像)
  # ============================================
  # 使用 {{ .LANGUAGE_DISPLAY_NAME }} 构建镜像
  
  # {{ .LANGUAGE_DISPLAY_NAME }} 构建镜像 (默认)
  - name: "BUILDER_IMAGE_X86"  # x86 构建镜像
    value: "{{ .BUILDER_IMAGE_X86 }}"
  - name: "BUILDER_IMAGE_ARM"  # arm64 构建镜像
    value: "{{ .BUILDER_IMAGE_ARM }}"
  
  # ============================================
  # Deployment Configuration (部署配置)
  # ============================================
  - name: "DEPLOY_DIR"  # 服务部署目录
    value: "{{ .DEPLOY_DIR }}"
`
}
