package commands

import (
	_ "embed"
	"fmt"

	"github.com/junjiewwang/service-template/pkg/utils"
	"github.com/spf13/cobra"
)

//go:embed templates/service.example.yaml
var serviceYamlExample string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new service.yaml configuration file",
	Long:  `Creates a new service.yaml configuration file with example values.`,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if service.yaml already exists
	if utils.FileExists(configFile) {
		return fmt.Errorf("configuration file %s already exists", configFile)
	}

	// Write example configuration
	if err := utils.WriteFile(configFile, getServiceYamlExample()); err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	fmt.Printf("✓ Created %s\n", configFile)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit service.yaml to configure your service")
	fmt.Println("  2. Run 'svcgen validate' to check your configuration")
	fmt.Println("  3. Run 'svcgen generate' to generate project files")

	return nil
}

func getServiceYamlExample() string {
	if serviceYamlExample != "" {
		return serviceYamlExample
	}

	// Fallback example if embed fails
	return `# ============================================
# Service Configuration
# Version: 2.0
# ============================================

service:
  name: my-service
  description: "My Service Description"
  ports:
    - name: http
      port: 8080
      protocol: TCP
      expose: true
      description: "HTTP API port"
  deploy_dir: /usr/local/services

language:
  type: go
  version: "1.23"
  config:
    goproxy: "https://goproxy.cn,direct"

build:
  dependency_files:
    auto_detect: true
  builder_image:
    amd64: "mirrors.tencent.com/tcs-infra/tceforqci_x86_go23:v1.0.0"
    arm64: "mirrors.tencent.com/tcs-infra/tceforqci_arm_go23:v1.0.0"
  runtime_image:
    amd64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"
    arm64: "mirrors.tencent.com/tencentos/tencentos3-minimal:latest"
  system_dependencies:
    build:
      packages:
        - git
        - make
  commands:
    build: |
      cd ${SERVICE_ROOT}
      go build -o ${BUILD_OUTPUT_DIR}/bin/${SERVICE_NAME} ./cmd/server
  output_dir: dist

runtime:
  healthcheck:
    enabled: true
    type: http
    http:
      path: /health
      port: 8080
      timeout: 3
  startup:
    command: |
      #!/bin/sh
      exec ${SERVICE_BIN_DIR}/${SERVICE_NAME}

local_dev:
  compose:
    volumes: []
  kubernetes:
    enabled: false
    namespace: default
    output_dir: k8s-manifests

metadata:
  template_version: "2.0.0"
generator: "svcgen"
`
}
