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

	// Add a second port
	cfg.Service.Ports = append(cfg.Service.Ports, config.PortConfig{
		Name:     "metrics",
		Port:     9090,
		Protocol: "TCP",
		Expose:   false,
	})

	ctx := context.NewGeneratorContext(cfg, "/tmp/output")
	gen, err := New(ctx)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	content, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Verify K8s Service structure
	if !strings.Contains(content, "apiVersion: v1") {
		t.Error("Expected apiVersion: v1 not found")
	}
	if !strings.Contains(content, "kind: Service") {
		t.Error("Expected kind: Service not found")
	}
	if !strings.Contains(content, "name: my-service") {
		t.Error("Expected service name not found")
	}
	if !strings.Contains(content, "app: my-service") {
		t.Error("Expected app label not found")
	}

	// Verify port with name (only expose=true ports should be included)
	if !strings.Contains(content, "name: http") {
		t.Error("Expected port name 'http' not found")
	}
	if !strings.Contains(content, "port: 8080") {
		t.Error("Expected port 8080 not found")
	}
	if !strings.Contains(content, "protocol: TCP") {
		t.Error("Expected protocol TCP not found")
	}

	// Verify non-exposed port is NOT included
	if strings.Contains(content, "name: metrics") {
		t.Error("Non-exposed port 'metrics' should NOT be in K8s Service")
	}
	if strings.Contains(content, "port: 9090") {
		t.Error("Non-exposed port 9090 should NOT be in K8s Service")
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

	// Verify exposed ports are included with correct names
	if !strings.Contains(content, "name: http") {
		t.Error("Expected port name 'http' not found")
	}
	if !strings.Contains(content, "port: 8080") {
		t.Error("Expected port 8080 not found")
	}
	if !strings.Contains(content, "name: grpc") {
		t.Error("Expected port name 'grpc' not found")
	}
	if !strings.Contains(content, "port: 9090") {
		t.Error("Expected port 9090 not found")
	}

	// Verify non-exposed ports are excluded
	if strings.Contains(content, "name: metrics") {
		t.Error("Non-exposed port 'metrics' should NOT be in K8s Service")
	}
	if strings.Contains(content, "name: debug") {
		t.Error("Non-exposed port 'debug' should NOT be in K8s Service")
	}
}

func TestGenerator_Generate_AllPortsExposed(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("full-service", "Full Service").
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
		{Name: "metrics", Port: 9090, Protocol: "TCP", Expose: true},
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

	if !strings.Contains(content, "name: http") {
		t.Error("Expected port name 'http' not found")
	}
	if !strings.Contains(content, "name: metrics") {
		t.Error("Expected port name 'metrics' not found")
	}
}

func TestGenerator_Generate_NoExposedPorts(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("internal-service", "Internal Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPort("http", 8080, "TCP", true).
		Build()

	// Override with no exposed ports
	cfg.Service.Ports = []config.PortConfig{
		{Name: "internal", Port: 8080, Protocol: "TCP", Expose: false},
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

	// Service should still be generated but without ports section
	if !strings.Contains(content, "kind: Service") {
		t.Error("Expected kind: Service not found")
	}
	// ports section should not be rendered when no exposed ports
	if strings.Contains(content, "port: 8080") {
		t.Error("Non-exposed port 8080 should NOT appear in ports section")
	}
	if strings.Contains(content, "targetPort:") {
		t.Error("targetPort should NOT appear when no ports are exposed")
	}
}

func TestGenerator_Generate_ProtocolUppercase(t *testing.T) {
	cfg := testutil.NewConfigBuilder().
		WithService("udp-service", "UDP Service").
		WithLanguage("go").
		WithBuilder("go_1.21", "golang:1.21", "golang:1.21").
		WithRuntime("alpine_3.18", "alpine:3.18", "alpine:3.18").
		WithBuilderImage("@builders.go_1.21").
		WithRuntimeImage("@runtimes.alpine_3.18").
		WithBuildCommand("go build -o bin/app").
		WithPort("http", 8080, "TCP", true).
		Build()

	// Use lowercase protocol in config
	cfg.Service.Ports = []config.PortConfig{
		{Name: "dns", Port: 53, Protocol: "udp", Expose: true},
		{Name: "http", Port: 8080, Protocol: "tcp", Expose: true},
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

	// Protocol should be uppercase in K8s Service
	if !strings.Contains(content, "protocol: UDP") {
		t.Error("Expected protocol UDP (uppercase) not found")
	}
	if !strings.Contains(content, "protocol: TCP") {
		t.Error("Expected protocol TCP (uppercase) not found")
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
