package service

import (
	"strings"
	"testing"

	"github.com/junjiewwang/service-template/pkg/config"
	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/internal/testutil"
)

func TestGenerator_Generate(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("my-service", "My Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPort("http", 8080, "TCP", true).
		Build()

	cfg.Service.Ports = []config.PortConfig{
		{Name: "http", Port: 8080, Protocol: "TCP", Expose: true},
		{Name: "metrics", Port: 9090, Protocol: "TCP", Expose: false},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify it's a strategic merge patch (minimal K8s Service)
	if !strings.Contains(content, "apiVersion: v1") {
		t.Error("Expected apiVersion: v1 not found")
	}
	if !strings.Contains(content, "kind: Service") {
		t.Error("Expected kind: Service not found")
	}
	if !strings.Contains(content, "name: my-service") {
		t.Error("Expected service name not found")
	}

	// All ports should be included (patch needs to cover all ports kompose generates)
	if !strings.Contains(content, "name: http") {
		t.Error("Expected port name 'http' not found")
	}
	if !strings.Contains(content, "port: 8080") {
		t.Error("Expected port 8080 not found")
	}
	if !strings.Contains(content, "name: metrics") {
		t.Error("Expected port name 'metrics' not found")
	}
	if !strings.Contains(content, "port: 9090") {
		t.Error("Expected port 9090 not found")
	}

	// Should NOT include protocol/targetPort (not needed in patch, minimal change)
	if strings.Contains(content, "targetPort:") {
		t.Error("Patch should NOT include targetPort (minimal change)")
	}
	if strings.Contains(content, "protocol:") {
		t.Error("Patch should NOT include protocol (minimal change)")
	}
}

func TestGenerator_Generate_MultiplePorts(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("multi-port-service", "Multi Port Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPort("http", 8080, "TCP", true).
		Build()

	cfg.Service.Ports = []config.PortConfig{
		{Name: "http", Port: 8080, Protocol: "TCP", Expose: true},
		{Name: "grpc", Port: 9090, Protocol: "TCP", Expose: true},
		{Name: "metrics", Port: 9100, Protocol: "TCP", Expose: false},
		{Name: "debug", Port: 6060, Protocol: "TCP", Expose: false},
	}

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// All ports should be included in the patch
	if !strings.Contains(content, "name: http") {
		t.Error("Expected port name 'http' not found")
	}
	if !strings.Contains(content, "name: grpc") {
		t.Error("Expected port name 'grpc' not found")
	}
	if !strings.Contains(content, "name: metrics") {
		t.Error("Expected port name 'metrics' not found")
	}
	if !strings.Contains(content, "name: debug") {
		t.Error("Expected port name 'debug' not found")
	}
}

func TestGenerator_Generate_NoPorts(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("no-port-service", "No Port Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPort("http", 8080, "TCP", true).
		Build()

	cfg.Service.Ports = nil

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Should still be a valid patch with metadata
	if !strings.Contains(content, "kind: Service") {
		t.Error("Expected kind: Service not found")
	}
	// Should not have ports section
	if strings.Contains(content, "port:") {
		t.Error("Should NOT have ports when none configured")
	}
}

func TestGenerator_GetName(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if gen.GetName() != GeneratorType {
		t.Errorf("Expected name %s, got %s", GeneratorType, gen.GetName())
	}
}

func TestGenerator_Validate(t *testing.T) {
	cfg := testutil.NewTestConfig()
	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, _ := New(ctx)

	if err := gen.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}
