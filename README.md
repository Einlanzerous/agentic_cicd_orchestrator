# LocalSprite

LocalSprite is an autonomous, Go-based QA agent designed for modularity and multi-environment flexibility. It orchestrates the process of analyzing a repository, planning test strategies, generating test code, and executing tests across different environments.

## Core Architecture

The system is built on dependency injection, allowing you to swap out LLM providers and execution environments seamlessly via profiles.

- **Planner**: Analyzes repository context to generate a test strategy.
- **Coder**: Generates executable test code based on the plan.
- **Executor**: Runs the generated tests in isolated environments (Docker).

## Project Structure

```text
.
├── cmd/
│   └── localsprite/
│       └── main.go           # Entry point: handles CLI flags & Dependency Injection
├── internal/
│   ├── agent/
│   │   └── interfaces.go     # Core interfaces: Planner, Coder, Executor
│   └── config/
│       └── config.go         # Viper configuration & profile loading
├── pkg/
│   └── providers/
│       ├── coder/            # Implementations: Bedrock, Anthropic, Local LLM
│       ├── executor/         # Implementations: Local Docker, Remote Docker
│       └── planner/          # Implementations: Gemini
├── config.yaml               # Profile definitions (Work, Home, LowCost)
└── README.md                 # This file
```

## Setup

### Prerequisites
- Go 1.24+
- Docker (local or remote access)
- API Keys for providers (Google GenAI, Anthropic, AWS Bedrock)

### Building
```bash
go build -o localsprite ./cmd/localsprite
```

## Usage

LocalSprite uses profiles defined in `config.yaml`. You can switch between environments using the `--profile` flag.

### Work Profile
Uses AWS Bedrock for coding and runs tests on the local machine.
```bash
./localsprite --profile=work
```

### Home Profile
Uses high-complexity Anthropic models and executes tests on a remote home server via SSH.
```bash
export ANTHROPIC_API_KEY=your-api-key
./localsprite --profile=home
```

### Home Low-Cost Profile
Uses a local LLM (Ollama/OpenAI compatible) and remote execution.
```bash
./localsprite --profile=home-lowcost
```

## Configuration

Profiles are managed in `config.yaml`. Example structure:

```yaml
profiles:
  work:
    planner:
      type: "gemini"
      model: "gemini-3.0-pro-latest"
    coder:
      type: "bedrock"
      model: "anthropic.claude-3-5-sonnet-20241022-v2:0"
    executor:
      type: "local_docker"
      params:
        image: "golang:1.24-alpine"
```

### Remote Execution Note
For `remote_docker` executors, ensure your `DOCKER_HOST` is configured (e.g., `ssh://user@home-server`). The system is designed to mount test-runner images specifically configured for your local hardware.
