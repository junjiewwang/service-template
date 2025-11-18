package runtime

// RuntimeDomain 运行时子域领域模型
type RuntimeDomain struct {
	// 运行时系统依赖
	SystemDependencies *SystemDependencies `yaml:"system_dependencies,omitempty"`

	// 健康检查配置
	Healthcheck *HealthcheckConfig `yaml:"healthcheck,omitempty"`

	// 启动配置
	Startup *StartupConfig `yaml:"startup"`
}

// SystemDependencies 系统依赖配置
type SystemDependencies struct {
	Packages []string `yaml:"packages,omitempty"`
}

// HealthcheckConfig 健康检查配置
type HealthcheckConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Type         string `yaml:"type"` // default | custom
	CustomScript string `yaml:"custom_script,omitempty"`
}

// StartupConfig 启动配置
type StartupConfig struct {
	Command string   `yaml:"command"`
	Env     []EnvVar `yaml:"env,omitempty"`
}

// EnvVar 环境变量
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// HasSystemDependencies 是否有系统依赖
func (d *RuntimeDomain) HasSystemDependencies() bool {
	return d.SystemDependencies != nil && len(d.SystemDependencies.Packages) > 0
}

// IsHealthcheckEnabled 是否启用健康检查
func (d *RuntimeDomain) IsHealthcheckEnabled() bool {
	return d.Healthcheck != nil && d.Healthcheck.Enabled
}

// IsCustomHealthcheck 是否使用自定义健康检查
func (d *RuntimeDomain) IsCustomHealthcheck() bool {
	return d.IsHealthcheckEnabled() && d.Healthcheck.Type == "custom"
}

// HasStartupEnv 是否有启动环境变量
func (d *RuntimeDomain) HasStartupEnv() bool {
	return d.Startup != nil && len(d.Startup.Env) > 0
}
