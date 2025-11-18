package build

import "errors"

var (
	// ErrBuilderImageRequired 构建镜像必需
	ErrBuilderImageRequired = errors.New("builder image is required")

	// ErrRuntimeImageRequired 运行时镜像必需
	ErrRuntimeImageRequired = errors.New("runtime image is required")

	// ErrBuildCommandRequired 构建命令必需
	ErrBuildCommandRequired = errors.New("build command is required")

	// ErrInvalidArchitecture 无效的架构
	ErrInvalidArchitecture = errors.New("invalid architecture")

	// ErrCustomPackageNameRequired 自定义包名称必需
	ErrCustomPackageNameRequired = errors.New("custom package name is required")

	// ErrCustomPackageInstallCommandRequired 自定义包安装命令必需
	ErrCustomPackageInstallCommandRequired = errors.New("custom package install command is required")

	// ErrDependencyFilesRequired 依赖文件列表必需
	ErrDependencyFilesRequired = errors.New("dependency files list is required when auto_detect is false")
)
