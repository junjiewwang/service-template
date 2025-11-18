package language

import "github.com/junjiewwang/service-template/pkg/generator/domain/chain"

type LanguageDomainFactory struct {
	*chain.BaseDomainFactory
}

func NewLanguageDomainFactory() chain.DomainFactory {
	factory := &LanguageDomainFactory{
		BaseDomainFactory: chain.NewBaseDomainFactory("language", 20),
	}
	// Language depends on service
	factory.SetDependencies([]string{"service"})
	return factory
}

func (f *LanguageDomainFactory) CreateParserHandler() chain.ParserHandler {
	return NewLanguageParserHandler()
}

func (f *LanguageDomainFactory) CreateValidatorHandler() chain.ValidatorHandler {
	return NewLanguageValidatorHandler()
}

func (f *LanguageDomainFactory) CreateGeneratorHandler() chain.GeneratorHandler {
	return NewLanguageGeneratorHandler()
}
