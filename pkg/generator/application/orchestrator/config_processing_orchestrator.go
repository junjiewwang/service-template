package orchestrator

import (
	"context"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/build"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/crossdomain"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/language"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/localdev"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/plugin"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/runtime"
	"github.com/junjiewwang/service-template/pkg/generator/domain/subdomains/service"
)

// ConfigProcessingOrchestrator orchestrates the configuration processing workflow
type ConfigProcessingOrchestrator struct {
	registry      *chain.DomainRegistry
	parseChain    chain.Handler
	validateChain chain.Handler
	generateChain chain.Handler
}

// NewConfigProcessingOrchestrator creates a new orchestrator
func NewConfigProcessingOrchestrator() *ConfigProcessingOrchestrator {
	return &ConfigProcessingOrchestrator{
		registry: chain.NewDomainRegistry(),
	}
}

// Initialize initializes the orchestrator with domain factories
func (o *ConfigProcessingOrchestrator) Initialize() error {
	// Register domain factories
	// Note: Using factory instances, not strings - zero hardcoding!
	err := o.registry.RegisterAll(
		service.NewServiceDomainFactory(),
		language.NewLanguageDomainFactory(),
		build.NewFactory(),
		plugin.NewFactory(),
		runtime.NewFactory(),
		localdev.NewFactory(),
		crossdomain.NewFactory(),
	)
	if err != nil {
		return fmt.Errorf("failed to register domain factories: %w", err)
	}

	// Build chains automatically based on factory priorities
	o.parseChain = o.registry.BuildParseChain()
	o.validateChain = o.registry.BuildValidateChain()
	o.generateChain = o.registry.BuildGenerateChain()

	return nil
}

// Process processes the configuration through all chains
func (o *ConfigProcessingOrchestrator) Process(ctx context.Context, rawConfig map[string]interface{}) (*chain.ProcessingContext, error) {
	// Create processing context
	procCtx := chain.NewProcessingContext(ctx, rawConfig)

	// Execute parse chain
	if err := o.parseChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("parse phase failed: %w", err)
	}

	// Execute validate chain
	if err := o.validateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("validate phase failed: %w", err)
	}

	// Execute generate chain
	if err := o.generateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("generate phase failed: %w", err)
	}

	return procCtx, nil
}

// ProcessWithPriorityChain processes using explicit priority chain
func (o *ConfigProcessingOrchestrator) ProcessWithPriorityChain(ctx context.Context, rawConfig map[string]interface{}) (*chain.ProcessingContext, error) {
	// Define priority chain with explicit ordering
	priorityChain := chain.NewPriorityChain().
		First(service.NewServiceDomainFactory()).
		Then(language.NewLanguageDomainFactory()).
		Then(build.NewFactory()).
		ThenAll(plugin.NewFactory(), runtime.NewFactory()).
		Then(localdev.NewFactory()).
		Finally(crossdomain.NewFactory())

	// Validate priority chain
	if err := priorityChain.Validate(); err != nil {
		return nil, fmt.Errorf("priority chain validation failed: %w", err)
	}

	// Build chains
	parseChain := priorityChain.BuildParseChain()
	validateChain := priorityChain.BuildValidateChain()
	generateChain := priorityChain.BuildGenerateChain()

	// Create processing context
	procCtx := chain.NewProcessingContext(ctx, rawConfig)

	// Execute chains
	if err := parseChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("parse phase failed: %w", err)
	}

	if err := validateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("validate phase failed: %w", err)
	}

	if err := generateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("generate phase failed: %w", err)
	}

	return procCtx, nil
}

// ProcessWithDependencyGraph processes using dependency graph
func (o *ConfigProcessingOrchestrator) ProcessWithDependencyGraph(ctx context.Context, rawConfig map[string]interface{}) (*chain.ProcessingContext, error) {
	// Define dependency graph
	serviceFactory := service.NewServiceDomainFactory()
	languageFactory := language.NewLanguageDomainFactory()
	buildFactory := build.NewFactory()
	pluginFactory := plugin.NewFactory()
	runtimeFactory := runtime.NewFactory()
	localdevFactory := localdev.NewFactory()
	crossdomainFactory := crossdomain.NewFactory()

	graph := chain.NewDependencyGraph().
		AddNode(serviceFactory).
		AddNode(languageFactory, "service").
		AddNode(buildFactory, "service", "language").
		AddNode(pluginFactory, "build").
		AddNode(runtimeFactory, "build").
		AddNode(localdevFactory, "plugin", "runtime").
		AddNode(crossdomainFactory, "localdev")

	// Validate dependency graph
	if err := graph.Validate(); err != nil {
		return nil, fmt.Errorf("dependency graph validation failed: %w", err)
	}

	// Build chains
	parseChain, err := graph.BuildParseChain()
	if err != nil {
		return nil, fmt.Errorf("failed to build parse chain: %w", err)
	}

	validateChain, err := graph.BuildValidateChain()
	if err != nil {
		return nil, fmt.Errorf("failed to build validate chain: %w", err)
	}

	generateChain, err := graph.BuildGenerateChain()
	if err != nil {
		return nil, fmt.Errorf("failed to build generate chain: %w", err)
	}

	// Create processing context
	procCtx := chain.NewProcessingContext(ctx, rawConfig)

	// Execute chains
	if err := parseChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("parse phase failed: %w", err)
	}

	if err := validateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("validate phase failed: %w", err)
	}

	if err := generateChain.Handle(procCtx); err != nil {
		return procCtx, fmt.Errorf("generate phase failed: %w", err)
	}

	return procCtx, nil
}
