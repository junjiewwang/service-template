package localdev

// LocalDevDomain 本地开发子域领域模型
type LocalDevDomain struct {
	// Docker Compose 配置
	Compose *ComposeConfig `yaml:"compose"`

	// Kubernetes 本地部署配置
	Kubernetes *KubernetesConfig `yaml:"kubernetes,omitempty"`
}

// ComposeConfig Docker Compose配置
type ComposeConfig struct {
	// 资源限制
	Resources *ResourcesConfig `yaml:"resources,omitempty"`

	// 环境变量配置
	Environment []EnvVar `yaml:"environment,omitempty"`

	// Entrypoint 配置
	Entrypoint []string `yaml:"entrypoint,omitempty"`

	// 卷挂载配置
	Volumes []VolumeMount `yaml:"volumes,omitempty"`

	// 健康检查配置
	Healthcheck *ComposeHealthcheck `yaml:"healthcheck,omitempty"`

	// 标签配置
	Labels map[string]string `yaml:"labels,omitempty"`
}

// ResourcesConfig 资源配置
type ResourcesConfig struct {
	Limits       *ResourceLimit `yaml:"limits,omitempty"`
	Reservations *ResourceLimit `yaml:"reservations,omitempty"`
}

// ResourceLimit 资源限制
type ResourceLimit struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

// EnvVar 环境变量
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// VolumeMount 卷挂载
type VolumeMount struct {
	Source      string `yaml:"source"`
	Target      string `yaml:"target"`
	Type        string `yaml:"type"` // bind | volume
	Description string `yaml:"description,omitempty"`
}

// ComposeHealthcheck Compose健康检查配置
type ComposeHealthcheck struct {
	Interval    string `yaml:"interval,omitempty"`
	Timeout     string `yaml:"timeout,omitempty"`
	Retries     int    `yaml:"retries,omitempty"`
	StartPeriod string `yaml:"start_period,omitempty"`
}

// KubernetesConfig Kubernetes配置
type KubernetesConfig struct {
	Enabled    bool           `yaml:"enabled"`
	Namespace  string         `yaml:"namespace,omitempty"`
	OutputDir  string         `yaml:"output_dir,omitempty"`
	VolumeType string         `yaml:"volume_type,omitempty"` // configMap | persistentVolumeClaim | emptyDir | hostPath
	Wait       *K8sWaitConfig `yaml:"wait,omitempty"`
}

// K8sWaitConfig K8s等待配置
type K8sWaitConfig struct {
	Enabled bool   `yaml:"enabled"`
	Timeout string `yaml:"timeout,omitempty"`
}

// HasVolumes 是否有卷挂载
func (c *ComposeConfig) HasVolumes() bool {
	return len(c.Volumes) > 0
}

// HasEnvironment 是否有环境变量
func (c *ComposeConfig) HasEnvironment() bool {
	return len(c.Environment) > 0
}

// HasEntrypoint 是否有自定义Entrypoint
func (c *ComposeConfig) HasEntrypoint() bool {
	return len(c.Entrypoint) > 0
}

// HasLabels 是否有标签
func (c *ComposeConfig) HasLabels() bool {
	return len(c.Labels) > 0
}

// IsK8sEnabled 是否启用Kubernetes
func (d *LocalDevDomain) IsK8sEnabled() bool {
	return d.Kubernetes != nil && d.Kubernetes.Enabled
}
