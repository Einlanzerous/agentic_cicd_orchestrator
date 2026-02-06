# LocalSprite Test Runner Images

Pre-configured Docker images for running generated tests.

## Available Images

### Go Test Runner
For executing generated Go unit tests.

```bash
docker build -t localsprite/go-test-runner:latest -f go-test-runner.Dockerfile .
```

### Playwright Test Runner
For executing generated Playwright UI tests.

```bash
docker build -t localsprite/playwright-test-runner:latest -f playwright-test-runner.Dockerfile .
```

### Cypress Test Runner
For executing generated Cypress E2E tests.

```bash
docker build -t localsprite/cypress-test-runner:latest -f cypress-test-runner.Dockerfile .
```

## Building All Images

```bash
# Build all test runner images
docker build -t localsprite/go-test-runner:latest -f go-test-runner.Dockerfile .
docker build -t localsprite/playwright-test-runner:latest -f playwright-test-runner.Dockerfile .
docker build -t localsprite/cypress-test-runner:latest -f cypress-test-runner.Dockerfile .
```

## Usage with LocalSprite

Configure the image in your `config.yaml`:

```yaml
profiles:
  my-go-tests:
    executor:
      type: "remote_docker"
      params:
        host: "ssh://imperial-construct"
        image: "localsprite/go-test-runner:latest"
        command: "go,test,-v,./..."
```

## Pushing to Remote Registry (Optional)

If you want to use a registry:

```bash
# Tag and push to your registry
docker tag localsprite/go-test-runner:latest your-registry/localsprite/go-test-runner:latest
docker push your-registry/localsprite/go-test-runner:latest
```

## Notes

- Images are pre-configured with basic test frameworks
- Generated test files are mounted at the working directory
- Commands can be overridden via config.yaml
