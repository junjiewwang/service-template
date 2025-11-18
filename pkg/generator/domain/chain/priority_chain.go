package chain

import "fmt"

// PriorityChain provides a fluent API for defining domain processing order
type PriorityChain struct {
	orderedFactories []DomainFactory
	factoryMap       map[string]DomainFactory
}

// NewPriorityChain creates a new priority chain
func NewPriorityChain() *PriorityChain {
	return &PriorityChain{
		orderedFactories: make([]DomainFactory, 0),
		factoryMap:       make(map[string]DomainFactory),
	}
}

// First adds the first domain factory to the chain
func (c *PriorityChain) First(factory DomainFactory) *PriorityChain {
	c.orderedFactories = append(c.orderedFactories, factory)
	c.factoryMap[factory.GetName()] = factory
	return c
}

// Then adds the next domain factory to the chain
func (c *PriorityChain) Then(factory DomainFactory) *PriorityChain {
	c.orderedFactories = append(c.orderedFactories, factory)
	c.factoryMap[factory.GetName()] = factory
	return c
}

// ThenAll adds multiple domain factories that can be processed in parallel
func (c *PriorityChain) ThenAll(factories ...DomainFactory) *PriorityChain {
	for _, factory := range factories {
		c.orderedFactories = append(c.orderedFactories, factory)
		c.factoryMap[factory.GetName()] = factory
	}
	return c
}

// Finally adds the last domain factory to the chain
func (c *PriorityChain) Finally(factory DomainFactory) *PriorityChain {
	c.orderedFactories = append(c.orderedFactories, factory)
	c.factoryMap[factory.GetName()] = factory
	return c
}

// GetOrderedFactories returns the factories in the defined order
func (c *PriorityChain) GetOrderedFactories() []DomainFactory {
	return c.orderedFactories
}

// GetFactory retrieves a factory by name
func (c *PriorityChain) GetFactory(name string) (DomainFactory, bool) {
	factory, ok := c.factoryMap[name]
	return factory, ok
}

// BuildParseChain builds a parser chain from the priority chain
func (c *PriorityChain) BuildParseChain() Handler {
	builder := NewChainBuilder()
	for _, factory := range c.orderedFactories {
		if factory.IsEnabled() {
			if handler := factory.CreateParserHandler(); handler != nil {
				builder.Add(handler)
			}
		}
	}
	return builder.BuildWithLogging()
}

// BuildValidateChain builds a validator chain from the priority chain
func (c *PriorityChain) BuildValidateChain() Handler {
	builder := NewChainBuilder()
	for _, factory := range c.orderedFactories {
		if factory.IsEnabled() {
			if handler := factory.CreateValidatorHandler(); handler != nil {
				builder.Add(handler)
			}
		}
	}
	return builder.BuildWithLogging()
}

// BuildGenerateChain builds a generator chain from the priority chain
func (c *PriorityChain) BuildGenerateChain() Handler {
	builder := NewChainBuilder()
	for _, factory := range c.orderedFactories {
		if factory.IsEnabled() {
			if handler := factory.CreateGeneratorHandler(); handler != nil {
				builder.Add(handler)
			}
		}
	}
	return builder.BuildWithLogging()
}

// Validate validates the priority chain configuration
func (c *PriorityChain) Validate() error {
	if len(c.orderedFactories) == 0 {
		return fmt.Errorf("priority chain is empty")
	}

	// Check for duplicate names
	seen := make(map[string]bool)
	for _, factory := range c.orderedFactories {
		name := factory.GetName()
		if seen[name] {
			return fmt.Errorf("duplicate domain factory: %s", name)
		}
		seen[name] = true
	}

	return nil
}
