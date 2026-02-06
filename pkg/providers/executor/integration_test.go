// +build integration

package executor

import (
	"os"
	"strings"
	"testing"
)

// Integration tests require Docker to be running
// Run with: go test -tags=integration ./...

func TestLocalDockerExecutor_Integration(t *testing.T) {
	// Skip if Docker is not available
	if os.Getenv("DOCKER_HOST") == "" && !dockerAvailable() {
		t.Skip("Docker not available, skipping integration test")
	}

	cfg := ExecutorConfig{
		Image:           "golang:1.22-alpine",
		Command:         []string{"go", "version"},
		WorkDir:         "/app",
		TestFilePattern: "test.go",
		Timeout:         60,
	}

	exec := NewLocalDockerExecutor(cfg)

	// Execute with minimal Go code
	code := `package main

func main() {}
`

	output, err := exec.Execute(code)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should contain Go version info
	if !strings.Contains(output, "go version") {
		t.Errorf("expected output to contain 'go version', got: %s", output)
	}
}

func TestLocalDockerExecutor_GoTest_Integration(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" && !dockerAvailable() {
		t.Skip("Docker not available, skipping integration test")
	}

	cfg := DefaultGoConfig()
	cfg.Image = "golang:1.22-alpine" // Use a known stable version
	exec := NewLocalDockerExecutor(cfg)

	// A simple passing test
	code := `package main

import "testing"

func TestAdd(t *testing.T) {
	result := 1 + 1
	if result != 2 {
		t.Errorf("expected 2, got %d", result)
	}
}

func Add(a, b int) int {
	return a + b
}
`

	output, err := exec.Execute(code)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should indicate test passed
	if !strings.Contains(output, "PASS") && !strings.Contains(output, "ok") {
		t.Errorf("expected test to pass, got: %s", output)
	}
}

func TestLocalDockerExecutor_FailingTest_Integration(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" && !dockerAvailable() {
		t.Skip("Docker not available, skipping integration test")
	}

	cfg := DefaultGoConfig()
	cfg.Image = "golang:1.22-alpine"
	exec := NewLocalDockerExecutor(cfg)

	// A failing test
	code := `package main

import "testing"

func TestFail(t *testing.T) {
	t.Error("intentional failure")
}
`

	output, err := exec.Execute(code)
	// Execute should not error even if test fails - we capture the output
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should indicate test failed
	if !strings.Contains(output, "FAIL") {
		t.Errorf("expected test to fail, got: %s", output)
	}
}

// dockerAvailable checks if Docker daemon is accessible
func dockerAvailable() bool {
	// Try to stat the Docker socket
	_, err := os.Stat("/var/run/docker.sock")
	return err == nil
}
