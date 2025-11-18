package build

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ParserHandler Build子域解析处理器
type ParserHandler struct {
	*chain.BaseHandler
}

// NewParserHandler 创建Build解析处理器
func NewParserHandler() chain.ParserHandler {
	return &ParserHandler{
		BaseHandler: chain.NewBaseHandler("build-parser"),
	}
}

// Parse implements ParserHandler interface
func (h *ParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理解析逻辑
func (h *ParserHandler) Handle(ctx *chain.ProcessingContext) error {
	// 从原始配置中解析Build配置
	rawConfig, ok := ctx.RawConfig["build"]
	if !ok {
		// Build配置是可选的
		return h.CallNext(ctx)
	}

	buildMap, ok := rawConfig.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid build configuration format")
	}

	domain := &BuildDomain{}

	// 解析依赖文件配置
	if depFiles, ok := buildMap["dependency_files"].(map[string]interface{}); ok {
		domain.DependencyFiles = &DependencyFilesConfig{}
		if autoDetect, ok := depFiles["auto_detect"].(bool); ok {
			domain.DependencyFiles.AutoDetect = autoDetect
		}
		if files, ok := depFiles["files"].([]interface{}); ok {
			for _, f := range files {
				if file, ok := f.(string); ok {
					domain.DependencyFiles.Files = append(domain.DependencyFiles.Files, file)
				}
			}
		}
	}

	// 解析构建镜像
	if builderImage, ok := buildMap["builder_image"].(map[string]interface{}); ok {
		domain.BuilderImage = make(map[string]string)
		for arch, img := range builderImage {
			if imgStr, ok := img.(string); ok {
				domain.BuilderImage[arch] = imgStr
			}
		}
	}

	// 解析运行时镜像
	if runtimeImage, ok := buildMap["runtime_image"].(map[string]interface{}); ok {
		domain.RuntimeImage = make(map[string]string)
		for arch, img := range runtimeImage {
			if imgStr, ok := img.(string); ok {
				domain.RuntimeImage[arch] = imgStr
			}
		}
	}

	// 解析依赖配置
	if deps, ok := buildMap["dependencies"].(map[string]interface{}); ok {
		domain.Dependencies = &BuildDependencies{}

		// 系统包
		if sysPkgs, ok := deps["system_pkgs"].([]interface{}); ok {
			for _, pkg := range sysPkgs {
				if pkgStr, ok := pkg.(string); ok {
					domain.Dependencies.SystemPackages = append(domain.Dependencies.SystemPackages, pkgStr)
				}
			}
		}

		// 自定义包
		if customPkgs, ok := deps["custom_pkgs"].([]interface{}); ok {
			for _, pkg := range customPkgs {
				if pkgMap, ok := pkg.(map[string]interface{}); ok {
					customPkg := CustomPackage{}
					if name, ok := pkgMap["name"].(string); ok {
						customPkg.Name = name
					}
					if desc, ok := pkgMap["description"].(string); ok {
						customPkg.Description = desc
					}
					if cmd, ok := pkgMap["install_command"].(string); ok {
						customPkg.InstallCommand = cmd
					}
					if req, ok := pkgMap["required"].(bool); ok {
						customPkg.Required = req
					}
					domain.Dependencies.CustomPackages = append(domain.Dependencies.CustomPackages, customPkg)
				}
			}
		}
	}

	// 解析构建命令
	if commands, ok := buildMap["commands"].(map[string]interface{}); ok {
		domain.Commands = &BuildCommands{}
		if preBuild, ok := commands["pre_build"].(string); ok {
			domain.Commands.PreBuild = preBuild
		}
		if build, ok := commands["build"].(string); ok {
			domain.Commands.Build = build
		}
		if postBuild, ok := commands["post_build"].(string); ok {
			domain.Commands.PostBuild = postBuild
		}
	}

	// 存储到上下文
	ctx.SetDomainModel("build", domain)

	return h.CallNext(ctx)
}

// ValidatorHandler Build子域校验处理器
type ValidatorHandler struct {
	*chain.BaseHandler
}

// NewValidatorHandler 创建Build校验处理器
func NewValidatorHandler() chain.ValidatorHandler {
	return &ValidatorHandler{
		BaseHandler: chain.NewBaseHandler("build-validator"),
	}
}

// Validate implements ValidatorHandler interface
func (h *ValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理校验逻辑
func (h *ValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("build")
	if !ok {
		// Build配置是可选的
		return h.CallNext(ctx)
	}

	buildDomain, ok := domain.(*BuildDomain)
	if !ok {
		return fmt.Errorf("invalid build domain model type")
	}

	// 校验构建镜像
	if len(buildDomain.BuilderImage) == 0 {
		ctx.AddValidationError("build.builder_image", ErrBuilderImageRequired)
	} else {
		// 至少需要一个架构的镜像
		hasValidImage := false
		for arch, img := range buildDomain.BuilderImage {
			if img != "" {
				hasValidImage = true
			} else {
				ctx.AddValidationError(fmt.Sprintf("build.builder_image.%s", arch),
					fmt.Errorf("builder image for %s cannot be empty", arch))
			}
		}
		if !hasValidImage {
			ctx.AddValidationError("build.builder_image", ErrBuilderImageRequired)
		}
	}

	// 校验运行时镜像
	if len(buildDomain.RuntimeImage) == 0 {
		ctx.AddValidationError("build.runtime_image", ErrRuntimeImageRequired)
	} else {
		hasValidImage := false
		for arch, img := range buildDomain.RuntimeImage {
			if img != "" {
				hasValidImage = true
			} else {
				ctx.AddValidationError(fmt.Sprintf("build.runtime_image.%s", arch),
					fmt.Errorf("runtime image for %s cannot be empty", arch))
			}
		}
		if !hasValidImage {
			ctx.AddValidationError("build.runtime_image", ErrRuntimeImageRequired)
		}
	}

	// 校验构建命令
	if buildDomain.Commands == nil || buildDomain.Commands.Build == "" {
		ctx.AddValidationError("build.commands.build", ErrBuildCommandRequired)
	}

	// 校验依赖文件配置
	if buildDomain.DependencyFiles != nil && !buildDomain.DependencyFiles.AutoDetect {
		if len(buildDomain.DependencyFiles.Files) == 0 {
			ctx.AddValidationError("build.dependency_files.files", ErrDependencyFilesRequired)
		}
	}

	// 校验自定义包
	if buildDomain.HasCustomDependencies() {
		for i, pkg := range buildDomain.Dependencies.CustomPackages {
			if pkg.Name == "" {
				ctx.AddValidationError(fmt.Sprintf("build.dependencies.custom_pkgs[%d].name", i),
					ErrCustomPackageNameRequired)
			}
			if pkg.InstallCommand == "" {
				ctx.AddValidationError(fmt.Sprintf("build.dependencies.custom_pkgs[%d].install_command", i),
					ErrCustomPackageInstallCommandRequired)
			}
		}
	}

	return h.CallNext(ctx)
}

// GeneratorHandler Build子域生成处理器
type GeneratorHandler struct {
	*chain.BaseHandler
}

// NewGeneratorHandler 创建Build生成处理器
func NewGeneratorHandler() chain.GeneratorHandler {
	return &GeneratorHandler{
		BaseHandler: chain.NewBaseHandler("build-generator"),
	}
}

// Generate implements GeneratorHandler interface
func (h *GeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}

// Handle 处理生成逻辑
func (h *GeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	domain, ok := ctx.GetDomainModel("build")
	if !ok {
		// Build配置是可选的
		return h.CallNext(ctx)
	}

	buildDomain, ok := domain.(*BuildDomain)
	if !ok {
		return fmt.Errorf("invalid build domain model type")
	}

	// 记录生成的文件
	ctx.AddGeneratedFile("build.sh", []byte("# Build script"))
	ctx.AddGeneratedFile("build_deps_install.sh", []byte("# Build dependencies installation script"))

	// 为每个架构生成Dockerfile
	for arch := range buildDomain.BuilderImage {
		ctx.AddGeneratedFile(fmt.Sprintf("Dockerfile.%s", arch),
			[]byte(fmt.Sprintf("# Dockerfile for %s architecture", arch)))
	}

	// 添加元数据
	ctx.SetMetadata("build.has_system_deps", buildDomain.HasSystemDependencies())
	ctx.SetMetadata("build.has_custom_deps", buildDomain.HasCustomDependencies())
	ctx.SetMetadata("build.architectures", getArchitectures(buildDomain.BuilderImage))

	return h.CallNext(ctx)
}

// getArchitectures 获取支持的架构列表
func getArchitectures(images map[string]string) []string {
	archs := make([]string, 0, len(images))
	for arch := range images {
		if arch != "default" {
			archs = append(archs, arch)
		}
	}
	return archs
}
