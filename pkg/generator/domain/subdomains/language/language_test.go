package language

import (
	"context"
	"testing"

	"github.com/junjiewwang/service-template/pkg/generator/domain/chain"
)

func TestLanguageConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  LanguageConfig
		wantErr bool
	}{
		{
			name:    "valid go config",
			config:  LanguageConfig{Type: "go"},
			wantErr: false,
		},
		{
			name:    "missing type",
			config:  LanguageConfig{},
			wantErr: true,
		},
		{
			name:    "unsupported language",
			config:  LanguageConfig{Type: "cobol"},
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

func TestLanguageParserHandler(t *testing.T) {
	rawConfig := map[string]interface{}{
		"language": map[string]interface{}{
			"type": "go",
			"config": map[string]interface{}{
				"goproxy": "https://goproxy.cn",
			},
		},
	}

	ctx := chain.NewProcessingContext(context.Background(), rawConfig)
	handler := NewLanguageParserHandler()

	err := handler.Handle(ctx)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	model, ok := ctx.GetDomainModel("language")
	if !ok {
		t.Fatal("language model not found")
	}

	config, ok := model.(*LanguageConfig)
	if !ok {
		t.Fatal("invalid language model type")
	}

	if config.Type != "go" {
		t.Errorf("Type = %s, want go", config.Type)
	}
}

func TestLanguageDomainFactory(t *testing.T) {
	factory := NewLanguageDomainFactory()

	if factory.GetName() != "language" {
		t.Errorf("GetName() = %s, want language", factory.GetName())
	}

	if factory.GetPriority() != 20 {
		t.Errorf("GetPriority() = %d, want 20", factory.GetPriority())
	}

	deps := factory.GetDependencies()
	if len(deps) != 1 || deps[0] != "service" {
		t.Errorf("GetDependencies() = %v, want [service]", deps)
	}
}
