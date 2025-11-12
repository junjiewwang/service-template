package dockerfile

import (
	"strings"

	"github.com/junjiewwang/service-template/pkg/config"
)

// getDependencyFilesList returns list of dependency files
func getDependencyFilesList(cfg *config.ServiceConfig) []string {
	if cfg.Build.DependencyFiles.AutoDetect {
		return getDefaultDependencyFiles(cfg.Language.Type)
	}
	return cfg.Build.DependencyFiles.Files
}

// getDepsInstallCommand generates dependency installation command
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

// getDefaultDependencyFiles returns default dependency files for a language
func getDefaultDependencyFiles(language string) []string {
	switch language {
	case "go":
		return []string{"go.mod", "go.sum"}
	case "python":
		return []string{"requirements.txt"}
	case "nodejs":
		return []string{"package.json", "package-lock.json"}
	case "java":
		return []string{"pom.xml"}
	default:
		return []string{}
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
