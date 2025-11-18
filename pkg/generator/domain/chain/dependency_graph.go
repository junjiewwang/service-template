package chain

import (
	"fmt"
	"sort"
)

// DependencyGraph manages domain dependencies and topological sorting
type DependencyGraph struct {
	nodes map[string]*GraphNode
}

// GraphNode represents a node in the dependency graph
type GraphNode struct {
	Factory      DomainFactory
	Dependencies []string
	Dependents   []string
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string]*GraphNode),
	}
}

// AddNode adds a domain factory with its dependencies
func (g *DependencyGraph) AddNode(factory DomainFactory, dependencies ...string) *DependencyGraph {
	name := factory.GetName()

	node := &GraphNode{
		Factory:      factory,
		Dependencies: dependencies,
		Dependents:   make([]string, 0),
	}

	g.nodes[name] = node

	// Update dependents
	for _, dep := range dependencies {
		if depNode, exists := g.nodes[dep]; exists {
			depNode.Dependents = append(depNode.Dependents, name)
		}
	}

	return g
}

// GetNode retrieves a node by name
func (g *DependencyGraph) GetNode(name string) (*GraphNode, bool) {
	node, ok := g.nodes[name]
	return node, ok
}

// TopologicalSort performs topological sort on the dependency graph
// Returns factories in dependency order (dependencies first)
func (g *DependencyGraph) TopologicalSort() ([]DomainFactory, error) {
	// Create a copy of in-degree map
	inDegree := make(map[string]int)
	for name, node := range g.nodes {
		inDegree[name] = len(node.Dependencies)
	}

	// Queue for nodes with no dependencies
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// Sort queue to ensure deterministic order
	sort.Strings(queue)

	// Result list
	result := make([]DomainFactory, 0, len(g.nodes))

	// Process nodes
	for len(queue) > 0 {
		// Pop from queue
		current := queue[0]
		queue = queue[1:]

		node := g.nodes[current]
		if node.Factory.IsEnabled() {
			result = append(result, node.Factory)
		}

		// Reduce in-degree for dependents
		for _, dependent := range node.Dependents {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
				sort.Strings(queue) // Keep sorted for determinism
			}
		}
	}

	// Check for cycles
	if len(result) < len(g.getEnabledNodes()) {
		return nil, fmt.Errorf("circular dependency detected in domain graph")
	}

	return result, nil
}

// BuildParseChain builds a parser chain using topological sort
func (g *DependencyGraph) BuildParseChain() (Handler, error) {
	factories, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateParserHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging(), nil
}

// BuildValidateChain builds a validator chain using topological sort
func (g *DependencyGraph) BuildValidateChain() (Handler, error) {
	factories, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateValidatorHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging(), nil
}

// BuildGenerateChain builds a generator chain using topological sort
func (g *DependencyGraph) BuildGenerateChain() (Handler, error) {
	factories, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	builder := NewChainBuilder()
	for _, factory := range factories {
		if handler := factory.CreateGeneratorHandler(); handler != nil {
			builder.Add(handler)
		}
	}

	return builder.BuildWithLogging(), nil
}

// Validate validates the dependency graph
func (g *DependencyGraph) Validate() error {
	// Check for missing dependencies
	for name, node := range g.nodes {
		for _, dep := range node.Dependencies {
			if _, exists := g.nodes[dep]; !exists {
				return fmt.Errorf("domain %s depends on non-existent domain %s", name, dep)
			}
		}
	}

	// Check for cycles by attempting topological sort
	_, err := g.TopologicalSort()
	return err
}

// getEnabledNodes returns all enabled nodes
func (g *DependencyGraph) getEnabledNodes() []*GraphNode {
	nodes := make([]*GraphNode, 0)
	for _, node := range g.nodes {
		if node.Factory.IsEnabled() {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetDependencies returns all dependencies for a domain
func (g *DependencyGraph) GetDependencies(name string) []string {
	if node, exists := g.nodes[name]; exists {
		return node.Dependencies
	}
	return nil
}

// GetDependents returns all dependents for a domain
func (g *DependencyGraph) GetDependents(name string) []string {
	if node, exists := g.nodes[name]; exists {
		return node.Dependents
	}
	return nil
}
