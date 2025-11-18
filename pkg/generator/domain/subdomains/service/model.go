package service

// ServiceConfig represents the service configuration domain model
type ServiceConfig struct {
	Name        string       `yaml:"name" json:"name"`
	Description string       `yaml:"description" json:"description"`
	Ports       []PortConfig `yaml:"ports" json:"ports"`
	DeployDir   string       `yaml:"deploy_dir" json:"deploy_dir"`
}

// PortConfig represents a service port configuration
type PortConfig struct {
	Name        string `yaml:"name" json:"name"`
	Port        int    `yaml:"port" json:"port"`
	Protocol    string `yaml:"protocol" json:"protocol"`
	Expose      bool   `yaml:"expose" json:"expose"`
	Description string `yaml:"description" json:"description"`
}

// GetMainPort returns the first port (main service port)
func (c *ServiceConfig) GetMainPort() int {
	if len(c.Ports) > 0 {
		return c.Ports[0].Port
	}
	return 0
}

// GetDefaultDeployDir returns the default deployment directory
func (c *ServiceConfig) GetDefaultDeployDir() string {
	if c.DeployDir != "" {
		return c.DeployDir
	}
	return "/usr/local/services"
}

// HasExposedPorts checks if the service has any exposed ports
func (c *ServiceConfig) HasExposedPorts() bool {
	for _, port := range c.Ports {
		if port.Expose {
			return true
		}
	}
	return false
}

// Validate validates the service configuration
func (c *ServiceConfig) Validate() error {
	if c.Name == "" {
		return ErrServiceNameRequired
	}

	if len(c.Ports) == 0 {
		return ErrPortsRequired
	}

	for i, port := range c.Ports {
		if port.Name == "" {
			return NewErrInvalidPort(i, "port name is required")
		}
		if port.Port <= 0 || port.Port > 65535 {
			return NewErrInvalidPort(i, "port must be between 1 and 65535")
		}
		if port.Protocol == "" {
			port.Protocol = "TCP"
		}
	}

	return nil
}
