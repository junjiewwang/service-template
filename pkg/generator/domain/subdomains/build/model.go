package build

// BuildDomain 构建子域领域模型
type BuildDomain struct {
	// 依赖文件配置
	DependencyFiles *DependencyFilesConfig `yaml:"dependency_files,omitempty"`

	// 构建基础镜像（按架构）
	BuilderImage map[string]string `yaml:"builder_image"`

	// 运行时基础镜像（按架构）
	RuntimeImage map[string]string `yaml:"runtime_image"`

	// 构建阶段依赖配置
	Dependencies *BuildDependencies `yaml:"dependencies,omitempty"`

	// 构建命令（三阶段）
	Commands *BuildCommands `yaml:"commands"`
}

// DependencyFilesConfig 依赖文件配置
type DependencyFilesConfig struct {
	// 是否自动检测依赖文件
	AutoDetect bool `yaml:"auto_detect"`

	// 自定义依赖文件列表（当 auto_detect=false 时使用）
	Files []string `yaml:"files,omitempty"`
}

// BuildDependencies 构建依赖配置
type BuildDependencies struct {
	// 系统包列表（通过包管理器安装）
	SystemPackages []string `yaml:"system_pkgs,omitempty"`

	// 自定义包列表（通过自定义命令安装）
	CustomPackages []CustomPackage `yaml:"custom_pkgs,omitempty"`
}

// CustomPackage 自定义包配置
type CustomPackage struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description"`
	InstallCommand string `yaml:"install_command"`
	Required       bool   `yaml:"required"`
}

// BuildCommands 构建命令配置
type BuildCommands struct {
	PreBuild  string `yaml:"pre_build,omitempty"`
	Build     string `yaml:"build"`
	PostBuild string `yaml:"post_build,omitempty"`
}

// GetBuilderImage 获取指定架构的构建镜像
func (d *BuildDomain) GetBuilderImage(arch string) string {
	if img, ok := d.BuilderImage[arch]; ok {
		return img
	}
	// 尝试默认架构
	if img, ok := d.BuilderImage["default"]; ok {
		return img
	}
	return ""
}

// GetRuntimeImage 获取指定架构的运行时镜像
func (d *BuildDomain) GetRuntimeImage(arch string) string {
	if img, ok := d.RuntimeImage[arch]; ok {
		return img
	}
	// 尝试默认架构
	if img, ok := d.RuntimeImage["default"]; ok {
		return img
	}
	return ""
}

// HasSystemDependencies 是否有系统依赖
func (d *BuildDomain) HasSystemDependencies() bool {
	return d.Dependencies != nil && len(d.Dependencies.SystemPackages) > 0
}

// HasCustomDependencies 是否有自定义依赖
func (d *BuildDomain) HasCustomDependencies() bool {
	return d.Dependencies != nil && len(d.Dependencies.CustomPackages) > 0
}

// GetRequiredCustomPackages 获取必需的自定义包
func (d *BuildDomain) GetRequiredCustomPackages() []CustomPackage {
	if !d.HasCustomDependencies() {
		return nil
	}

	var required []CustomPackage
	for _, pkg := range d.Dependencies.CustomPackages {
		if pkg.Required {
			required = append(required, pkg)
		}
	}
	return required
}
