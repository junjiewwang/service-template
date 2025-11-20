package testutil

import "github.com/junjiewwang/service-template/pkg/config"

// ServiceBuilder 服务配置构建器
type ServiceBuilder struct {
	info *config.ServiceInfo
}

// Name 设置服务名称
func (b *ServiceBuilder) Name(name string) *ServiceBuilder {
	b.info.Name = name
	return b
}

// Description 设置服务描述
func (b *ServiceBuilder) Description(desc string) *ServiceBuilder {
	b.info.Description = desc
	return b
}

// DeployDir 设置部署目录
func (b *ServiceBuilder) DeployDir(dir string) *ServiceBuilder {
	b.info.DeployDir = dir
	return b
}

// AddPort 添加端口配置
func (b *ServiceBuilder) AddPort(port int, protocol string) *ServiceBuilder {
	b.info.Ports = append(b.info.Ports, config.PortConfig{
		Port:     port,
		Protocol: protocol,
		Name:     protocol,
	})
	return b
}

// AddPortWithName 添加带名称的端口配置
func (b *ServiceBuilder) AddPortWithName(port int, protocol, name string) *ServiceBuilder {
	b.info.Ports = append(b.info.Ports, config.PortConfig{
		Port:     port,
		Protocol: protocol,
		Name:     name,
	})
	return b
}

// AddPortConfig 添加完整的端口配置
func (b *ServiceBuilder) AddPortConfig(cfg config.PortConfig) *ServiceBuilder {
	b.info.Ports = append(b.info.Ports, cfg)
	return b
}

// SetPorts 设置端口列表（替换现有）
func (b *ServiceBuilder) SetPorts(ports []config.PortConfig) *ServiceBuilder {
	b.info.Ports = ports
	return b
}
