package service

import (
	"fmt"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

// ServiceValidatorHandler validates service configuration
type ServiceValidatorHandler struct {
	*chain.BaseHandler
}

// NewServiceValidatorHandler creates a new service validator handler
func NewServiceValidatorHandler() chain.ValidatorHandler {
	return &ServiceValidatorHandler{
		BaseHandler: chain.NewBaseHandler("service-validator"),
	}
}

// Handle processes the service configuration validation
func (h *ServiceValidatorHandler) Handle(ctx *chain.ProcessingContext) error {
	// Get parsed service model
	model, ok := ctx.GetDomainModel("service")
	if !ok {
		return fmt.Errorf("service model not found")
	}

	config, ok := model.(*ServiceConfig)
	if !ok {
		return fmt.Errorf("invalid service model type")
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		ctx.AddValidationError("service", err)
		return err
	}

	return h.CallNext(ctx)
}

// Validate implements ValidatorHandler interface
func (h *ServiceValidatorHandler) Validate(ctx *chain.ProcessingContext) error {
	return h.Handle(ctx)
}
