# Playwright Test Runner
# Used for executing generated Playwright UI tests
FROM mcr.microsoft.com/playwright:v1.40.0-jammy

# Set up working directory
WORKDIR /app

# Initialize a basic package.json for test execution
RUN echo '{"name":"test-runner","type":"module","devDependencies":{"@playwright/test":"^1.40.0"}}' > package.json

# Install dependencies
RUN npm install

# Create basic playwright config
RUN echo 'import { defineConfig } from "@playwright/test"; \n\
export default defineConfig({ \n\
  testDir: ".", \n\
  timeout: 30000, \n\
  retries: 0, \n\
  use: { \n\
    headless: true, \n\
    screenshot: "only-on-failure", \n\
  }, \n\
});' > playwright.config.ts

# Default command
CMD ["npx", "playwright", "test"]
