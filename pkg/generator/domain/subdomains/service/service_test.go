package service

import (
	"context"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

func TestServiceConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServiceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ServiceConfig{
				Name:        "test-service",
				Description: "Test Service",
				Ports: []PortConfig{
					{Name: "http", Port: 8080, Protocol: "TCP"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: ServiceConfig{
				Ports: []PortConfig{
					{Name: "http", Port: 8080},
				},
			},
			wantErr: true,
		},
		{
			name: "missing ports",
			config: ServiceConfig{
				Name: "test-service",
			},
			wantErr: true,
		},
		{
			name: "invalid port number",
			config: ServiceConfig{
				Name: "test-service",
				Ports: []PortConfig{
					{Name: "http", Port: 70000},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceConfig_GetMainPort(t *testing.T) {
	config := ServiceConfig{
		Ports: []PortConfig{
			{Name: "http", Port: 8080},
			{Name: "metrics", Port: 9090},
		},
	}

	port := config.GetMainPort()
	if port != 8080 {
		t.Errorf("GetMainPort() = %d, want 8080", port)
	}
}

func TestServiceParserHandler(t *testing.T) {
	rawConfig := map[string]interface{}{
		"service": map[string]interface{}{
			"name":        "test-service",
			"description": "Test Service",
			"ports": []interface{}{
				map[string]interface{}{
					"name": "http",
					"port": 8080,
				},
			},
		},
	}

	ctx := chain.NewProcessingContext(context.Background(), rawConfig)
	handler := NewServiceParserHandler()

	err := handler.Handle(ctx)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	model, ok := ctx.GetDomainModel("service")
	if !ok {
		t.Fatal("service model not found")
	}

	config, ok := model.(*ServiceConfig)
	if !ok {
		t.Fatal("invalid service model type")
	}

	if config.Name != "test-service" {
		t.Errorf("Name = %s, want test-service", config.Name)
	}
}

func TestServiceValidatorHandler(t *testing.T) {
	config := &ServiceConfig{
		Name: "test-service",
		Ports: []PortConfig{
			{Name: "http", Port: 8080},
		},
	}

	ctx := chain.NewProcessingContext(context.Background(), nil)
	ctx.SetDomainModel("service", config)

	handler := NewServiceValidatorHandler()
	err := handler.Handle(ctx)

	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}

func TestServiceGeneratorHandler(t *testing.T) {
	config := &ServiceConfig{
		Name:        "test-service",
		Description: "Test Service",
		Ports: []PortConfig{
			{Name: "http", Port: 8080},
		},
	}

	ctx := chain.NewProcessingContext(context.Background(), nil)
	ctx.SetDomainModel("service", config)

	handler := NewServiceGeneratorHandler()
	err := handler.Handle(ctx)

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check generated file
	content, ok := ctx.GetGeneratedFile("service-metadata.txt")
	if !ok {
		t.Fatal("service-metadata.txt not generated")
	}

	if len(content) == 0 {
		t.Error("generated file is empty")
	}

	// Check metadata
	name, ok := ctx.GetMetadata("service_name")
	if !ok || name != "test-service" {
		t.Errorf("service_name metadata = %v, want test-service", name)
	}
}

func TestServiceDomainFactory(t *testing.T) {
	factory := NewServiceDomainFactory()

	if factory.GetName() != "service" {
		t.Errorf("GetName() = %s, want service", factory.GetName())
	}

	if factory.GetPriority() != 10 {
		t.Errorf("GetPriority() = %d, want 10", factory.GetPriority())
	}

	parser := factory.CreateParserHandler()
	if parser == nil {
		t.Error("CreateParserHandler() returned nil")
	}

	validator := factory.CreateValidatorHandler()
	if validator == nil {
		t.Error("CreateValidatorHandler() returned nil")
	}

	generator := factory.CreateGeneratorHandler()
	if generator == nil {
		t.Error("CreateGeneratorHandler() returned nil")
	}
}
