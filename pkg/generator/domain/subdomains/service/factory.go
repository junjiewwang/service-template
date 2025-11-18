package service

import "github.com/junjiewwang/service-template/pkg/generator/domain/chain"

// ServiceDomainFactory creates service domain handlers
type ServiceDomainFactory struct {
	*chain.BaseDomainFactory
}

// NewServiceDomainFactory creates a new service domain factory
func NewServiceDomainFactory() chain.DomainFactory {
	factory := &ServiceDomainFactory{
		BaseDomainFactory: chain.NewBaseDomainFactory("service", 10),
	}
	return factory
}

// CreateParserHandler creates a service parser handler
func (f *ServiceDomainFactory) CreateParserHandler() chain.ParserHandler {
	return NewServiceParserHandler()
}

// CreateValidatorHandler creates a service validator handler
func (f *ServiceDomainFactory) CreateValidatorHandler() chain.ValidatorHandler {
	return NewServiceValidatorHandler()
}

// CreateGeneratorHandler creates a service generator handler
func (f *ServiceDomainFactory) CreateGeneratorHandler() chain.GeneratorHandler {
	return NewServiceGeneratorHandler()
}
