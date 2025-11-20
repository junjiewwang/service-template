package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// BuildConfigBuilder 构建配置构建器
type BuildConfigBuilder struct {
	build *config.BuildConfig
}

// BuilderImage 设置构建镜像引用
func (b *BuildConfigBuilder) BuilderImage(ref string) *BuildConfigBuilder {
	b.build.BuilderImage = config.ImageRef(ref)
	return b
}

// RuntimeImage 设置运行时镜像引用
func (b *BuildConfigBuilder) RuntimeImage(ref string) *BuildConfigBuilder {
	b.build.RuntimeImage = config.ImageRef(ref)
	return b
}

// BuildCommand 设置构建命令
func (b *BuildConfigBuilder) BuildCommand(cmd string) *BuildConfigBuilder {
	b.build.Commands.Build = cmd
	return b
}

// PreBuildCommand 设置预构建命令
func (b *BuildConfigBuilder) PreBuildCommand(cmd string) *BuildConfigBuilder {
	b.build.Commands.PreBuild = cmd
	return b
}

// PostBuildCommand 设置后构建命令
func (b *BuildConfigBuilder) PostBuildCommand(cmd string) *BuildConfigBuilder {
	b.build.Commands.PostBuild = cmd
	return b
}

// AddSystemPackage 添加系统包依赖
func (b *BuildConfigBuilder) AddSystemPackage(pkg string) *BuildConfigBuilder {
	b.build.Dependencies.SystemPkgs = append(b.build.Dependencies.SystemPkgs, pkg)
	return b
}

// SetSystemPackages 设置系统包依赖列表
func (b *BuildConfigBuilder) SetSystemPackages(pkgs []string) *BuildConfigBuilder {
	b.build.Dependencies.SystemPkgs = pkgs
	return b
}

// AutoDetectDependencies 设置是否自动检测依赖文件
func (b *BuildConfigBuilder) AutoDetectDependencies(autoDetect bool) *BuildConfigBuilder {
	b.build.DependencyFiles.AutoDetect = autoDetect
	return b
}

// AddDependencyFile 添加依赖文件
func (b *BuildConfigBuilder) AddDependencyFile(file string) *BuildConfigBuilder {
	b.build.DependencyFiles.Files = append(b.build.DependencyFiles.Files, file)
	return b
}

// SetDependencyFiles 设置依赖文件列表
func (b *BuildConfigBuilder) SetDependencyFiles(files []string) *BuildConfigBuilder {
	b.build.DependencyFiles.Files = files
	return b
}
