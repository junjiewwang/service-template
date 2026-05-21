package service

import (
	_ "embed"
	"strings"

	"github.com/junjiewwang/service-template/pkg/generator/context"
	"github.com/junjiewwang/service-template/pkg/generator/core"
)

const GeneratorType = "k8s-service"

// init registers the k8s service generator
func init() {
	core.DefaultRegistry.Register(GeneratorType, New)
}

// Generator generates Kubernetes Service manifest
type Generator struct {
	core.BaseGenerator
}

// New creates a new k8s service generator
func New(ctx *context.GeneratorContext, options ...interface{}) (core.Generator, error) {
	engine := core.NewTemplateEngine()
	return &Generator{
		BaseGenerator: core.NewBaseGenerator(GeneratorType, ctx, engine),
	}, nil
}

// Generate generates Kubernetes Service manifest content
func (g *Generator) Generate() (string, error) {
	if err := g.Validate(); err != nil {
		return "", err
	}

	vars := g.prepareTemplateVars()
	return g.RenderTemplate(tmpl, vars)
}

// ServicePort represents a port entry in K8s Service spec
type ServicePort struct {
	Name     string
	Port     int
	Protocol string
}

// prepareTemplateVars prepares variables for k8s service template
func (g *Generator) prepareTemplateVars() map[string]interface{} {
	ctx := g.GetContext()

	// Collect ports that are marked as expose: true
	var servicePorts []ServicePort
	for _, port := range ctx.Config.Service.Ports {
		if !port.Expose {
			continue
		}
		servicePorts = append(servicePorts, ServicePort{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: strings.ToUpper(port.Protocol),
		})
	}

	return map[string]interface{}{
		"SERVICE_NAME":  ctx.Config.Service.Name,
		"SERVICE_PORTS": servicePorts,
	}
}

//go:embed templates/service.yaml.tmpl
var tmpl string
