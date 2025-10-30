package generator

import (
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// DevOpsGenerator generates TAD DevOps configuration
type DevOpsGenerator struct {
	config         *config.ServiceConfig
	templateEngine *TemplateEngine
	variables      *Variables
}

// NewDevOpsGenerator creates a new DevOps generator
func NewDevOpsGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *DevOpsGenerator {
	return &DevOpsGenerator{
		config:         cfg,
		templateEngine: engine,
		variables:      vars,
	}
}

// Generate generates devops.yaml content in TAD format
func (g *DevOpsGenerator) Generate() (string, error) {
	// Use embedded template
	templateContent := g.getDefaultDevOpsTemplate()

	// Prepare template variables
	vars := g.prepareTemplateVars()

	return g.templateEngine.Render(templateContent, vars)
}

// prepareTemplateVars prepares variables for devops template
func (g *DevOpsGenerator) prepareTemplateVars() map[string]interface{} {
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

	// Language examples
	vars["SHOW_PYTHON_EXAMPLE"] = g.config.Language.Type != "python"
	vars["SHOW_NODEJS_EXAMPLE"] = g.config.Language.Type != "nodejs"
	vars["SHOW_JAVA_EXAMPLE"] = g.config.Language.Type != "java"

	return vars
}

// getDefaultDevOpsTemplate returns embedded default devops template
func (g *DevOpsGenerator) getDefaultDevOpsTemplate() string {
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
  # 根据项目语言选择合适的构建镜像
  # 默认使用 {{ .LANGUAGE_DISPLAY_NAME }} 构建镜像，如需其他语言请修改以下配置
  
  # {{ .LANGUAGE_DISPLAY_NAME }} 构建镜像 (默认)
  - name: "BUILDER_IMAGE_X86"  # x86 构建镜像
    value: "{{ .BUILDER_IMAGE_X86 }}"
  - name: "BUILDER_IMAGE_ARM"  # arm64 构建镜像
    value: "{{ .BUILDER_IMAGE_ARM }}"
{{- if .SHOW_PYTHON_EXAMPLE }}
  
  # Python 构建镜像示例 (如需使用 Python，请取消注释并注释掉上面的构建镜像)
  # - name: "BUILDER_IMAGE_X86"
  #   value: "mirrors.tencent.com/tcs-infra/python:3.11-slim"
  # - name: "BUILDER_IMAGE_ARM"
  #   value: "mirrors.tencent.com/tcs-infra/python:3.11-slim"
{{- end }}
{{- if .SHOW_NODEJS_EXAMPLE }}
  
  # Node.js 构建镜像示例 (如需使用 Node.js，请取消注释并注释掉上面的构建镜像)
  # - name: "BUILDER_IMAGE_X86"
  #   value: "mirrors.tencent.com/tcs-infra/node:18-alpine"
  # - name: "BUILDER_IMAGE_ARM"
  #   value: "mirrors.tencent.com/tcs-infra/node:18-alpine"
{{- end }}
{{- if .SHOW_JAVA_EXAMPLE }}
  
  # Java 构建镜像示例 (如需使用 Java，请取消注释并注释掉上面的构建镜像)
  # - name: "BUILDER_IMAGE_X86"
  #   value: "mirrors.tencent.com/tcs-infra/maven:3.8-openjdk-17"
  # - name: "BUILDER_IMAGE_ARM"
  #   value: "mirrors.tencent.com/tcs-infra/maven:3.8-openjdk-17"
{{- end }}
  
  # ============================================
  # Deployment Configuration (部署配置)
  # ============================================
  - name: "DEPLOY_DIR"  # 服务部署目录
    value: "{{ .DEPLOY_DIR }}"
`
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


