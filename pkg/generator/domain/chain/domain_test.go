package chain

import (
	"context"
	"testing"
)

// Mock domain factory for testing
type mockDomainFactory struct {
	*BaseDomainFactory
	parserHandler    ParserHandler
	validatorHandler ValidatorHandler
	generatorHandler GeneratorHandler
}

func newMockDomainFactory(name string, priority int) *mockDomainFactory {
	return &mockDomainFactory{
		BaseDomainFactory: NewBaseDomainFactory(name, priority),
		parserHandler:     newMockParserHandler(name + "-parser"),
		validatorHandler:  newMockValidatorHandler(name + "-validator"),
		generatorHandler:  newMockGeneratorHandler(name + "-generator"),
	}
}

func (f *mockDomainFactory) CreateParserHandler() ParserHandler {
	return f.parserHandler
}

func (f *mockDomainFactory) CreateValidatorHandler() ValidatorHandler {
	return f.validatorHandler
}

func (f *mockDomainFactory) CreateGeneratorHandler() GeneratorHandler {
	return f.generatorHandler
}

// TestDomainRegistry tests the domain registry
func TestDomainRegistry(t *testing.T) {
	t.Run("Register and get factory", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory := newMockDomainFactory("service", 10)

		err := registry.Register(factory)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		retrieved, ok := registry.Get("service")
		if !ok {
			t.Error("Expected to find factory")
		}

		if retrieved.GetName() != "service" {
			t.Error("Factory name mismatch")
		}
	})

	t.Run("Register duplicate factory", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("service", 20)

		registry.Register(factory1)
		err := registry.Register(factory2)

		if err == nil {
			t.Error("Expected error for duplicate registration")
		}
	})

	t.Run("Register multiple factories", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)
		factory3 := newMockDomainFactory("build", 30)

		err := registry.RegisterAll(factory1, factory2, factory3)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if registry.Count() != 3 {
			t.Errorf("Expected 3 factories, got %d", registry.Count())
		}
	})

	t.Run("Build parse chain", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)

		registry.RegisterAll(factory1, factory2)

		chain := registry.BuildParseChain()
		if chain == nil {
			t.Error("Expected non-nil chain")
		}

		ctx := NewProcessingContext(context.Background(), nil)
		err := chain.Handle(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Get enabled factories", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)
		factory2.SetEnabled(false)

		registry.RegisterAll(factory1, factory2)

		enabled := registry.GetEnabled()
		if len(enabled) != 1 {
			t.Errorf("Expected 1 enabled factory, got %d", len(enabled))
		}
	})

	t.Run("Unregister factory", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory := newMockDomainFactory("service", 10)

		registry.Register(factory)
		registry.Unregister("service")

		_, ok := registry.Get("service")
		if ok {
			t.Error("Expected factory to be unregistered")
		}
	})

	t.Run("Clear registry", func(t *testing.T) {
		registry := NewDomainRegistry()
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)

		registry.RegisterAll(factory1, factory2)
		registry.Clear()

		if registry.Count() != 0 {
			t.Errorf("Expected 0 factories after clear, got %d", registry.Count())
		}
	})
}

// TestPriorityChain tests the priority chain
func TestPriorityChain(t *testing.T) {
	t.Run("Build priority chain", func(t *testing.T) {
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)
		factory3 := newMockDomainFactory("build", 30)

		chain := NewPriorityChain().
			First(factory1).
			Then(factory2).
			Finally(factory3)

		factories := chain.GetOrderedFactories()
		if len(factories) != 3 {
			t.Errorf("Expected 3 factories, got %d", len(factories))
		}

		if factories[0].GetName() != "service" {
			t.Error("First factory should be service")
		}

		if factories[2].GetName() != "build" {
			t.Error("Last factory should be build")
		}
	})

	t.Run("Build parse chain from priority chain", func(t *testing.T) {
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)

		priorityChain := NewPriorityChain().
			First(factory1).
			Then(factory2)

		parseChain := priorityChain.BuildParseChain()
		if parseChain == nil {
			t.Error("Expected non-nil parse chain")
		}

		ctx := NewProcessingContext(context.Background(), nil)
		err := parseChain.Handle(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("ThenAll adds multiple factories", func(t *testing.T) {
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("plugin", 40)
		factory3 := newMockDomainFactory("runtime", 50)

		chain := NewPriorityChain().
			First(factory1).
			ThenAll(factory2, factory3)

		factories := chain.GetOrderedFactories()
		if len(factories) != 3 {
			t.Errorf("Expected 3 factories, got %d", len(factories))
		}
	})

	t.Run("Validate priority chain", func(t *testing.T) {
		factory1 := newMockDomainFactory("service", 10)
		factory2 := newMockDomainFactory("language", 20)

		chain := NewPriorityChain().
			First(factory1).
			Then(factory2)

		err := chain.Validate()
		if err != nil {
			t.Errorf("Unexpected validation error: %v", err)
		}
	})

	t.Run("Validate empty chain", func(t *testing.T) {
		chain := NewPriorityChain()

		err := chain.Validate()
		if err == nil {
			t.Error("Expected validation error for empty chain")
		}
	})
}

// TestDependencyGraph tests the dependency graph
func TestDependencyGraph(t *testing.T) {
	t.Run("Add nodes and topological sort", func(t *testing.T) {
		graph := NewDependencyGraph()

		serviceFactory := newMockDomainFactory("service", 10)
		languageFactory := newMockDomainFactory("language", 20)
		buildFactory := newMockDomainFactory("build", 30)

		graph.AddNode(serviceFactory)
		graph.AddNode(languageFactory, "service")
		graph.AddNode(buildFactory, "service", "language")

		sorted, err := graph.TopologicalSort()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(sorted) != 3 {
			t.Errorf("Expected 3 factories, got %d", len(sorted))
		}

		// Service should be first
		if sorted[0].GetName() != "service" {
			t.Error("Service should be first in topological order")
		}

		// Build should be last
		if sorted[2].GetName() != "build" {
			t.Error("Build should be last in topological order")
		}
	})

	t.Run("Detect circular dependency", func(t *testing.T) {
		graph := NewDependencyGraph()

		factory1 := newMockDomainFactory("a", 10)
		factory2 := newMockDomainFactory("b", 20)
		factory3 := newMockDomainFactory("c", 30)

		graph.AddNode(factory1, "c")
		graph.AddNode(factory2, "a")
		graph.AddNode(factory3, "b")

		_, err := graph.TopologicalSort()
		if err == nil {
			t.Error("Expected error for circular dependency")
		}
	})

	t.Run("Build parse chain from dependency graph", func(t *testing.T) {
		graph := NewDependencyGraph()

		serviceFactory := newMockDomainFactory("service", 10)
		languageFactory := newMockDomainFactory("language", 20)

		graph.AddNode(serviceFactory)
		graph.AddNode(languageFactory, "service")

		chain, err := graph.BuildParseChain()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if chain == nil {
			t.Error("Expected non-nil chain")
		}

		ctx := NewProcessingContext(context.Background(), nil)
		err = chain.Handle(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Validate dependency graph", func(t *testing.T) {
		graph := NewDependencyGraph()

		serviceFactory := newMockDomainFactory("service", 10)
		languageFactory := newMockDomainFactory("language", 20)

		graph.AddNode(serviceFactory)
		graph.AddNode(languageFactory, "service")

		err := graph.Validate()
		if err != nil {
			t.Errorf("Unexpected validation error: %v", err)
		}
	})

	t.Run("Validate missing dependency", func(t *testing.T) {
		graph := NewDependencyGraph()

		factory := newMockDomainFactory("language", 20)
		graph.AddNode(factory, "service") // service doesn't exist

		err := graph.Validate()
		if err == nil {
			t.Error("Expected validation error for missing dependency")
		}
	})

	t.Run("Get dependencies and dependents", func(t *testing.T) {
		graph := NewDependencyGraph()

		serviceFactory := newMockDomainFactory("service", 10)
		languageFactory := newMockDomainFactory("language", 20)
		buildFactory := newMockDomainFactory("build", 30)

		graph.AddNode(serviceFactory)
		graph.AddNode(languageFactory, "service")
		graph.AddNode(buildFactory, "service", "language")

		// Check dependencies
		buildDeps := graph.GetDependencies("build")
		if len(buildDeps) != 2 {
			t.Errorf("Expected 2 dependencies for build, got %d", len(buildDeps))
		}

		// Check dependents
		serviceDeps := graph.GetDependents("service")
		if len(serviceDeps) != 2 {
			t.Errorf("Expected 2 dependents for service, got %d", len(serviceDeps))
		}
	})
}
