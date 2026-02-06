package executor

// ExecutorConfig holds configuration for test execution
type ExecutorConfig struct {
	// Host is the Docker host (e.g., "ssh://imperial-construct" or empty for local)
	Host string

	// Image is the Docker image to use for running tests
	Image string

	// Command is the test command to run (e.g., ["go", "test", "-v", "./..."])
	Command []string

	// WorkDir is the working directory inside the container (default: /app)
	WorkDir string

	// TestFilePattern is the filename pattern for generated test files
	// e.g., "generated_test.go" for Go, "generated.spec.ts" for Playwright
	TestFilePattern string

	// Timeout in seconds for test execution (default: 300)
	Timeout int
}

// DefaultGoConfig returns default configuration for Go tests
func DefaultGoConfig() ExecutorConfig {
	return ExecutorConfig{
		Image:           "golang:1.24-alpine",
		Command:         []string{"go", "test", "-v", "./..."},
		WorkDir:         "/app",
		TestFilePattern: "generated_test.go",
		Timeout:         300,
	}
}

// DefaultPlaywrightConfig returns default configuration for Playwright tests
func DefaultPlaywrightConfig() ExecutorConfig {
	return ExecutorConfig{
		Image:           "mcr.microsoft.com/playwright:v1.40.0-jammy",
		Command:         []string{"npx", "playwright", "test"},
		WorkDir:         "/app",
		TestFilePattern: "generated.spec.ts",
		Timeout:         600,
	}
}

// DefaultCypressConfig returns default configuration for Cypress tests
func DefaultCypressConfig() ExecutorConfig {
	return ExecutorConfig{
		Image:           "cypress/included:13.6.0",
		Command:         []string{"cypress", "run"},
		WorkDir:         "/e2e",
		TestFilePattern: "generated.cy.ts",
		Timeout:         600,
	}
}
