# Cypress Test Runner
# Used for executing generated Cypress E2E tests
FROM cypress/included:13.6.0

# Set up working directory
WORKDIR /e2e

# Initialize a basic package.json for test execution
RUN echo '{"name":"test-runner","devDependencies":{"cypress":"^13.6.0"}}' > package.json

# Create basic cypress config
RUN echo 'const { defineConfig } = require("cypress"); \n\
module.exports = defineConfig({ \n\
  e2e: { \n\
    supportFile: false, \n\
    specPattern: "**/*.cy.{js,ts}", \n\
  }, \n\
  video: false, \n\
  screenshotOnRunFailure: true, \n\
});' > cypress.config.js

# Create cypress directory structure
RUN mkdir -p cypress/e2e

# Default command
CMD ["cypress", "run"]
