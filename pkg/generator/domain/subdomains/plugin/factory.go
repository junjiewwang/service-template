package plugin

import "github.com/junjiewwang/service-template/pkg/generator/domain/chain"

// Factory Plugin子域工厂
type Factory struct{}

// NewFactory 创建Plugin子域工厂
func NewFactory() chain.DomainFactory {
	return &Factory{}
}

// GetName 获取子域名称
func (f *Factory) GetName() string {
	return "plugin"
}

// GetPriority 获取优先级
func (f *Factory) GetPriority() int {
	return 40 // 在build之后
}

// IsEnabled 是否启用
func (f *Factory) IsEnabled() bool {
	return true
}

// CreateParserHandler 创建解析处理器
func (f *Factory) CreateParserHandler() chain.ParserHandler {
	return NewParserHandler()
}

// CreateValidatorHandler 创建校验处理器
func (f *Factory) CreateValidatorHandler() chain.ValidatorHandler {
	return NewValidatorHandler()
}

// CreateGeneratorHandler 创建生成处理器
func (f *Factory) CreateGeneratorHandler() chain.GeneratorHandler {
	return NewGeneratorHandler()
}

// GetDependencies 获取依赖列表
func (f *Factory) GetDependencies() []string {
	return []string{"build"}
}
