package service

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
	"gopkg.in/yaml.v3"
)

// ServiceParserHandler parses service configuration
type ServiceParserHandler struct {
	*chain.BaseHandler
}

// NewServiceParserHandler creates a new service parser handler
func NewServiceParserHandler() chain.ParserHandler {
	return &ServiceParserHandler{
		BaseHandler: chain.NewBaseHandler("service-parser"),
	}
}

// Handle processes the service configuration parsing
func (h *ServiceParserHandler) Handle(ctx *chain.ProcessingContext) error {
	// Get service config from raw data
	rawConfig := ctx.RawConfig
	serviceData, ok := rawConfig["service"]
	if !ok {
		return fmt.Errorf("service configuration not found")
	}

	// Convert to YAML and parse
	yamlData, err := yaml.Marshal(serviceData)
	if err != nil {
		return fmt.Errorf("failed to marshal service config: %w", err)
	}

	var config ServiceConfig
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return fmt.Errorf("failed to unmarshal service config: %w", err)
	}

	// Store parsed model
	ctx.SetDomainModel("service", &config)

	return h.CallNext(ctx)
}

// Parse implements ParserHandler interface
func (h *ServiceParserHandler) Parse(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}
