package service

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ServiceGeneratorHandler generates service-related files
type ServiceGeneratorHandler struct {
	*chain.BaseHandler
}

// NewServiceGeneratorHandler creates a new service generator handler
func NewServiceGeneratorHandler() chain.GeneratorHandler {
	return &ServiceGeneratorHandler{
		BaseHandler: chain.NewBaseHandler("service-generator"),
	}
}

// Handle processes the service file generation
func (h *ServiceGeneratorHandler) Handle(ctx *chain.ProcessingContext) error {
	// Get parsed service model
	model, ok := ctx.GetDomainModel("service")
	if !ok {
		return fmt.Errorf("service model not found")
	}

	config, ok := model.(*ServiceConfig)
	if !ok {
		return fmt.Errorf("invalid service model type")
	}

	// Generate service metadata file (example)
	metadata := fmt.Sprintf(`# Service Metadata
Name: %s
Description: %s
Main Port: %d
Deploy Directory: %s
`, config.Name, config.Description, config.GetMainPort(), config.GetDefaultDeployDir())

	ctx.AddGeneratedFile("service-metadata.txt", []byte(metadata))

	// Store service config for other domains to use
	ctx.SetMetadata("service_name", config.Name)
	ctx.SetMetadata("service_port", config.GetMainPort())
	ctx.SetMetadata("deploy_dir", config.GetDefaultDeployDir())

	return h.CallNext(ctx)
}

// Generate implements GeneratorHandler interface
func (h *ServiceGeneratorHandler) Generate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}
