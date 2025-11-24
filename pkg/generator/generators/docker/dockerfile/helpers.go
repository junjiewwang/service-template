package dockerfile

import (
	"github.com/junjiewwang/service-template/pkg/generator/domain/services/languageservice"
	"strings"

	"github.com/junjiewwang/service-template/pkg/generator/context"
)

// getDependencyFilesList returns list of dependency files
// projectDir: the root directory of the project to scan for dependency files
func getDependencyFilesList(ctx *context.GeneratorContext, projectDir string) []string {
	cfg := ctx.Config
	if cfg.Build.DependencyFiles.AutoDetect {
		// Use language service to detect actual dependency files in the project
		langService := languageservice.NewLanguageService(ctx)
		return langService.GetDependencyFilesWithDetection(
			cfg.Language.Type,
			projectDir,
			true,
			nil,
		)
	}
	return cfg.Build.DependencyFiles.Files
}

// getDepsInstallCommand generates dependency installation command
// Deprecated: Use LanguageService.GetDepsInstallCommand instead
func getDepsInstallCommand(language string) string {
	switch language {
	case "go":
		return "go mod download"
	case "python":
		return "pip install -r requirements.txt"
	case "nodejs":
		return "npm install"
	case "java":
		return "mvn dependency:go-offline"
	default:
		return "echo 'No dependency installation needed'"
	}
}

// detectPackageManager detects the package manager from the image name
func detectPackageManager(image string) string {
	imageLower := strings.ToLower(image)

	if strings.Contains(imageLower, "alpine") {
		return "apk"
	} else if strings.Contains(imageLower, "debian") || strings.Contains(imageLower, "ubuntu") {
		return "apt-get"
	} else if strings.Contains(imageLower, "centos") || strings.Contains(imageLower, "rhel") || strings.Contains(imageLower, "tencentos") {
		return "yum"
	} else if strings.Contains(imageLower, "fedora") {
		return "dnf"
	}

	return "yum"
}
