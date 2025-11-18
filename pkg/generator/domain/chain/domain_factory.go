package chain

// DomainFactory is the interface for creating domain-specific handlers
type DomainFactory interface {
	// GetName returns the domain name (used for registration and ordering)
	GetName() string

	// GetPriority returns the default priority for this domain
	// Lower values = higher priority (executed first)
	GetPriority() int

	// IsEnabled returns whether this domain is enabled
	IsEnabled() bool

	// CreateParserHandler creates a parser handler for this domain
	CreateParserHandler() ParserHandler

	// CreateValidatorHandler creates a validator handler for this domain
	CreateValidatorHandler() ValidatorHandler

	// CreateGeneratorHandler creates a generator handler for this domain
	CreateGeneratorHandler() GeneratorHandler

	// GetDependencies returns the list of domain names this domain depends on
	// These dependencies will be processed before this domain
	GetDependencies() []string
}

// BaseDomainFactory provides default implementation for DomainFactory
type BaseDomainFactory struct {
	name         string
	priority     int
	enabled      bool
	dependencies []string
}

// NewBaseDomainFactory creates a new base domain factory
func NewBaseDomainFactory(name string, priority int) *BaseDomainFactory {
	return &BaseDomainFactory{
		name:         name,
		priority:     priority,
		enabled:      true,
		dependencies: []string{},
	}
}

// GetName returns the domain name
func (f *BaseDomainFactory) GetName() string {
	return f.name
}

// GetPriority returns the priority
func (f *BaseDomainFactory) GetPriority() int {
	return f.priority
}

// IsEnabled returns whether the domain is enabled
func (f *BaseDomainFactory) IsEnabled() bool {
	return f.enabled
}

// SetEnabled sets the enabled state
func (f *BaseDomainFactory) SetEnabled(enabled bool) {
	f.enabled = enabled
}

// GetDependencies returns the dependencies
func (f *BaseDomainFactory) GetDependencies() []string {
	return f.dependencies
}

// SetDependencies sets the dependencies
func (f *BaseDomainFactory) SetDependencies(dependencies []string) {
	f.dependencies = dependencies
}

// CreateParserHandler should be overridden by concrete factories
func (f *BaseDomainFactory) CreateParserHandler() ParserHandler {
	return nil
}

// CreateValidatorHandler should be overridden by concrete factories
func (f *BaseDomainFactory) CreateValidatorHandler() ValidatorHandler {
	return nil
}

// CreateGeneratorHandler should be overridden by concrete factories
func (f *BaseDomainFactory) CreateGeneratorHandler() GeneratorHandler {
	return nil
}
