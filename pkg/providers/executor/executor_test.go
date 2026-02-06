package executor

import (
	"strings"
	"testing"
)

func TestDefaultGoConfig(t *testing.T) {
	cfg := DefaultGoConfig()

	if cfg.Image != "golang:1.24-alpine" {
		t.Errorf("expected image golang:1.24-alpine, got %s", cfg.Image)
	}
	if cfg.WorkDir != "/app" {
		t.Errorf("expected workdir /app, got %s", cfg.WorkDir)
	}
	if cfg.TestFilePattern != "generated_test.go" {
		t.Errorf("expected test file pattern generated_test.go, got %s", cfg.TestFilePattern)
	}
	if len(cfg.Command) != 4 || cfg.Command[0] != "go" {
		t.Errorf("expected go test command, got %v", cfg.Command)
	}
}

func TestDefaultPlaywrightConfig(t *testing.T) {
	cfg := DefaultPlaywrightConfig()

	if !strings.Contains(cfg.Image, "playwright") {
		t.Errorf("expected playwright image, got %s", cfg.Image)
	}
	if cfg.TestFilePattern != "generated.spec.ts" {
		t.Errorf("expected test file pattern generated.spec.ts, got %s", cfg.TestFilePattern)
	}
	if cfg.Timeout != 600 {
		t.Errorf("expected timeout 600, got %d", cfg.Timeout)
	}
}

func TestDefaultCypressConfig(t *testing.T) {
	cfg := DefaultCypressConfig()

	if !strings.Contains(cfg.Image, "cypress") {
		t.Errorf("expected cypress image, got %s", cfg.Image)
	}
	if cfg.WorkDir != "/e2e" {
		t.Errorf("expected workdir /e2e, got %s", cfg.WorkDir)
	}
	if cfg.TestFilePattern != "generated.cy.ts" {
		t.Errorf("expected test file pattern generated.cy.ts, got %s", cfg.TestFilePattern)
	}
}

func TestNewLocalDockerExecutor_AppliesDefaults(t *testing.T) {
	cfg := ExecutorConfig{
		Image: "test-image:latest",
	}

	exec := NewLocalDockerExecutor(cfg)

	if exec.Config.WorkDir != "/app" {
		t.Errorf("expected default workdir /app, got %s", exec.Config.WorkDir)
	}
	if exec.Config.TestFilePattern != "generated_test.go" {
		t.Errorf("expected default test file pattern, got %s", exec.Config.TestFilePattern)
	}
	if exec.Config.Timeout != 300 {
		t.Errorf("expected default timeout 300, got %d", exec.Config.Timeout)
	}
	if len(exec.Config.Command) == 0 {
		t.Error("expected default command to be set")
	}
}

func TestNewRemoteDockerExecutor_AppliesDefaults(t *testing.T) {
	cfg := ExecutorConfig{
		Host:  "ssh://test-host",
		Image: "test-image:latest",
	}

	exec := NewRemoteDockerExecutor(cfg)

	if exec.Config.Host != "ssh://test-host" {
		t.Errorf("expected host ssh://test-host, got %s", exec.Config.Host)
	}
	if exec.Config.WorkDir != "/app" {
		t.Errorf("expected default workdir /app, got %s", exec.Config.WorkDir)
	}
}

func TestNewLocalDockerExecutor_PreservesCustomConfig(t *testing.T) {
	cfg := ExecutorConfig{
		Image:           "custom-image:latest",
		Command:         []string{"npm", "test"},
		WorkDir:         "/custom",
		TestFilePattern: "custom.test.js",
		Timeout:         600,
	}

	exec := NewLocalDockerExecutor(cfg)

	if exec.Config.Image != "custom-image:latest" {
		t.Errorf("expected custom image, got %s", exec.Config.Image)
	}
	if exec.Config.WorkDir != "/custom" {
		t.Errorf("expected custom workdir, got %s", exec.Config.WorkDir)
	}
	if exec.Config.TestFilePattern != "custom.test.js" {
		t.Errorf("expected custom test file pattern, got %s", exec.Config.TestFilePattern)
	}
	if exec.Config.Timeout != 600 {
		t.Errorf("expected custom timeout, got %d", exec.Config.Timeout)
	}
	if len(exec.Config.Command) != 2 || exec.Config.Command[0] != "npm" {
		t.Errorf("expected custom command, got %v", exec.Config.Command)
	}
}
