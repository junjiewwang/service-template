package generator

import (
	"github.com/junjiewwang/service-template/pkg/config"
)

// Generator type constant
const GeneratorTypeMakefile = "makefile"

// MakefileTemplateGenerator generates Makefile using factory pattern
type MakefileTemplateGenerator struct {
	BaseTemplateGenerator
}

// init registers the Makefile generator
func init() {
	RegisterGenerator(GeneratorTypeMakefile, createMakefileGenerator)
}

// createMakefileGenerator is the creator function for Makefile generator
func createMakefileGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables, options ...interface{}) (TemplateGenerator, error) {
	return NewMakefileTemplateGenerator(cfg, engine, vars), nil
}

// NewMakefileTemplateGenerator creates a new Makefile template generator
func NewMakefileTemplateGenerator(cfg *config.ServiceConfig, engine *TemplateEngine, vars *Variables) *MakefileTemplateGenerator {
	return &MakefileTemplateGenerator{
		BaseTemplateGenerator: BaseTemplateGenerator{
			config:         cfg,
			templateEngine: engine,
			variables:      vars,
			name:           GeneratorTypeMakefile,
		},
	}
}

// Generate generates Makefile content
func (g *MakefileTemplateGenerator) Generate() (string, error) {
	vars := g.prepareTemplateVars()
	return g.RenderTemplate(g.getTemplate(), vars)
}

// prepareTemplateVars prepares variables for Makefile template
func (g *MakefileTemplateGenerator) prepareTemplateVars() map[string]interface{} {
	vars := make(map[string]interface{})

	// Basic info
	vars["GENERATED_AT"] = g.config.Metadata.GeneratedAt
	vars["SERVICE_NAME"] = g.config.Service.Name
	vars["SERVICE_PORT"] = g.variables.ServicePort
	vars["OUTPUT_DIR"] = g.config.Build.OutputDir

	// Kubernetes config
	vars["K8S_ENABLED"] = g.config.LocalDev.Kubernetes.Enabled
	vars["K8S_NAMESPACE"] = g.config.LocalDev.Kubernetes.Namespace
	vars["K8S_OUTPUT_DIR"] = g.config.LocalDev.Kubernetes.OutputDir
	vars["K8S_WAIT_ENABLED"] = g.config.LocalDev.Kubernetes.Wait.Enabled
	vars["K8S_WAIT_TIMEOUT"] = g.config.LocalDev.Kubernetes.Wait.Timeout

	// Custom targets
	vars["CUSTOM_TARGETS"] = g.config.Makefile.CustomTargets

	return vars
}

// getTemplate returns the Makefile template
func (g *MakefileTemplateGenerator) getTemplate() string {
	return `# Auto-generated Makefile
# Generated at: {{ .GENERATED_AT }}

SERVICE_NAME := {{ .SERVICE_NAME }}
IMAGE_TAG := latest
DOCKER_REGISTRY := 

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: docker-build
docker-build: docker-build-amd64 ## Build Docker images for all architectures

.PHONY: docker-build-amd64
docker-build-amd64: ## Build Docker image for amd64
	@echo "Building Docker image for amd64..."
	docker build -f .tad/build/$(SERVICE_NAME)/Dockerfile.$(SERVICE_NAME).amd64 -t $(SERVICE_NAME):$(IMAGE_TAG)-amd64 .
	@echo "✓ Image built: $(SERVICE_NAME):$(IMAGE_TAG)-amd64"

.PHONY: docker-build-arm64
docker-build-arm64: ## Build Docker image for arm64
	@echo "Building Docker image for arm64..."
	docker build -f .tad/build/$(SERVICE_NAME)/Dockerfile.$(SERVICE_NAME).arm64 -t $(SERVICE_NAME):$(IMAGE_TAG)-arm64 .
	@echo "✓ Image built: $(SERVICE_NAME):$(IMAGE_TAG)-arm64"

.PHONY: docker-up
docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	docker-compose -f compose.yaml up -d
	@echo "✓ Service is running on http://localhost:{{ .SERVICE_PORT }}"

.PHONY: docker-down
docker-down: ## Stop services with Docker Compose
	@echo "Stopping services..."
	docker-compose -f compose.yaml down
	@echo "✓ Services stopped"

.PHONY: docker-logs
docker-logs: ## Show Docker Compose logs
	docker-compose -f compose.yaml logs -f
{{- if .K8S_ENABLED }}

.PHONY: k8s-deploy
k8s-deploy: k8s-configmap ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f {{ .K8S_OUTPUT_DIR }}/ -n {{ .K8S_NAMESPACE }}
{{- if .K8S_WAIT_ENABLED }}
	kubectl wait --for=condition=ready pod -l app=$(SERVICE_NAME) -n {{ .K8S_NAMESPACE }} --timeout={{ .K8S_WAIT_TIMEOUT }}
{{- end }}
	@echo "✓ Deployed to Kubernetes"

.PHONY: k8s-configmap
k8s-configmap: ## Create Kubernetes ConfigMap
	@echo "Creating ConfigMap..."
	kubectl create configmap {{ .SERVICE_NAME }}-config --from-file=configs/ -n {{ .K8S_NAMESPACE }} --dry-run=client -o yaml | kubectl apply -f -
	@echo "✓ ConfigMap created"

.PHONY: k8s-delete
k8s-delete: ## Delete from Kubernetes
	@echo "Deleting from Kubernetes..."
	kubectl delete -f {{ .K8S_OUTPUT_DIR }}/ -n {{ .K8S_NAMESPACE }} --ignore-not-found=true
	@echo "✓ Deleted from Kubernetes"
{{- end }}

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning {{ .OUTPUT_DIR }}..."
	rm -rf {{ .OUTPUT_DIR }}
	@echo "✓ Cleaned"
{{- range .CUSTOM_TARGETS }}

.PHONY: {{ .Name }}
{{ .Name }}: ## {{ .Description }}
{{- range .Commands }}
	{{ . }}
{{- end }}
{{- end }}
`
}
