# LocalSprite

LocalSprite is an autonomous, Go-based QA agent designed for modularity and multi-environment flexibility. It orchestrates the process of analyzing a repository, planning test strategies, generating test code, and executing tests across different environments.

## Core Architecture

The system is built on dependency injection, allowing you to swap out LLM providers and execution environments seamlessly via profiles.

- **Planner**: Analyzes repository context to generate a test strategy.
- **Coder**: Generates executable test code based on the plan.
- **Executor**: Runs the generated tests in isolated Docker containers.

## Project Structure

```text
.
├── cmd/
│   └── localsprite/
│       └── main.go              # Entry point: CLI flags & Dependency Injection
├── internal/
│   ├── agent/
│   │   └── interfaces.go        # Core interfaces: Planner, Coder, Executor
│   └── config/
│       └── config.go            # Viper configuration & profile loading
├── pkg/
│   └── providers/
│       ├── coder/               # Implementations: Bedrock, Anthropic, Local LLM
│       ├── executor/            # Docker executors + ARCHITECTURE.md roadmap
│       └── planner/             # Implementations: Gemini, Local LLM
├── docker/                      # Test runner Dockerfiles
│   ├── go-test-runner.Dockerfile
│   ├── playwright-test-runner.Dockerfile
│   └── cypress-test-runner.Dockerfile
├── config.yaml                  # Profile definitions
└── README.md
```

## Setup

### Prerequisites
- Go 1.24+
- Docker (local or remote access via SSH)
- API Keys for cloud providers (optional - can run fully local)

### Building
```bash
go build -o localsprite ./cmd/localsprite
```

### Running Tests
```bash
# Unit tests
go test -v ./...

# Integration tests (requires Docker)
go test -v -tags=integration ./...
```

## Usage

LocalSprite uses profiles defined in `config.yaml`. Switch between environments using the `--profile` flag.

### Available Profiles

| Profile | Planner | Coder | Executor | Use Case |
|---------|---------|-------|----------|----------|
| `work` | Gemini | Bedrock Claude | Local Docker | Work environment with cloud APIs |
| `home` | Ollama (gemma3:12b) | Ollama (qwen2.5-coder) | Remote Docker | Cost-free local execution |
| `home-playwright` | Ollama | Ollama | Remote Docker + Playwright | UI testing |
| `home-cypress` | Ollama | Ollama | Remote Docker + Cypress | E2E testing |

### Examples

```bash
# Work profile - cloud APIs, local Docker
./localsprite --profile=work

# Home profile - fully local via Ollama
./localsprite --profile=home

# UI testing with Playwright
./localsprite --profile=home-playwright

# E2E testing with Cypress
./localsprite --profile=home-cypress
```

## Configuration

Profiles are managed in `config.yaml`. Each profile configures a planner, coder, and executor.

### Full Example

```yaml
profiles:
  home:
    planner:
      type: "local"
      model: "gemma3:12b"
      params:
        endpoint: "http://imperial-construct:11434/v1"
    coder:
      type: "local"
      model: "qwen2.5-coder:7b"
      params:
        endpoint: "http://imperial-construct:11434/v1"
    executor:
      type: "remote_docker"
      params:
        host: "ssh://imperial-construct"
        image: "golang:1.24-alpine"
        command: "go,test,-v,./..."
        workdir: "/app"
        test_file_pattern: "generated_test.go"
```

### Executor Configuration

| Parameter | Description | Default |
|-----------|-------------|---------|
| `host` | Docker host (e.g., `ssh://hostname`) | Local socket |
| `image` | Docker image for test execution | `golang:1.24-alpine` |
| `command` | Test command (comma-separated) | `go,test,-v,./...` |
| `workdir` | Working directory in container | `/app` |
| `test_file_pattern` | Generated test filename | `generated_test.go` |

### Provider Types

**Planner:**
- `gemini` - Google Gemini API
- `local` - Ollama/OpenAI-compatible endpoint

**Coder:**
- `bedrock` - AWS Bedrock (Claude)
- `anthropic` - Anthropic API (Claude)
- `local` - Ollama/OpenAI-compatible endpoint

**Executor:**
- `local_docker` - Local Docker daemon
- `remote_docker` - Remote Docker via SSH

## Test Runner Images

Pre-built Dockerfiles are available in `docker/`:

```bash
# Build Go test runner
docker build -t localsprite/go-test-runner:latest -f docker/go-test-runner.Dockerfile docker/

# Build Playwright test runner
docker build -t localsprite/playwright-test-runner:latest -f docker/playwright-test-runner.Dockerfile docker/

# Build Cypress test runner
docker build -t localsprite/cypress-test-runner:latest -f docker/cypress-test-runner.Dockerfile docker/
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ANTHROPIC_API_KEY` | API key for Anthropic provider |
| `GEMINI_API_KEY` | API key for Google Gemini |
| `AWS_REGION` | AWS region for Bedrock |

## Roadmap

See `pkg/providers/executor/ARCHITECTURE.md` for planned enhancements:
- Job queue architecture (NATS/Redis)
- Multi-provider worker pools
- n8n integration for notifications
