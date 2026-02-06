# Go Test Runner
# Used for executing generated Go tests
FROM golang:1.24-alpine

# Install common testing dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set up working directory
WORKDIR /app

# Create a basic go.mod for test execution
RUN go mod init testrunner

# Default command - can be overridden
CMD ["go", "test", "-v", "./..."]
