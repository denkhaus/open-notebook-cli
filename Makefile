# Makefile for OpenNotebook CLI

# Variables
BINARY_NAME=onb
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build
DIST_DIR=dist

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Platform variables
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: help build clean test lint fmt vet deps run install docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary for current platform"
	@echo "  build-all   - Build binaries for all platforms"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  test-cover  - Run tests with coverage"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format Go code"
	@echo "  vet         - Run go vet"
	@echo "  deps        - Download dependencies"
	@echo "  run         - Run the application"
	@echo "  install     - Install binary to system"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/onb

# Build for all platforms
build-all:
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(DIST_DIR)
	@$(foreach platform,$(PLATFORMS), \
		echo "Building for $(platform)..."; \
		GOOS=$(word 1,$(subst /, ,$(platform))) GOARCH=$(word 2,$(subst /, ,$(platform))) \
		$(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$(platform) ./cmd/onb; \
	)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f coverage.out

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.54.2"; \
		exit 1; \
	fi

# Format Go code
fmt:
	@echo "Formatting Go code..."
	$(GOCMD) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/onb
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t open-notebook-cli:$(VERSION) .
	docker tag open-notebook-cli:$(VERSION) open-notebook-cli:latest

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it \
		-e OPEN_NOTEBOOK_API_URL=$(OPEN_NOTEBOOK_API_URL) \
		-e OPEN_NOTEBOOK_PASSWORD=$(OPEN_NOTEBOOK_PASSWORD) \
		-v $(PWD)/config:/app/config \
		open-notebook-cli:latest

# Development targets
dev-setup:
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/air-verse/air
	@if [ ! -f .air.toml ]; then \
		echo "Creating air.toml for hot reload..."; \
		echo 'root = "."' > .air.toml; \
		echo 'testdata_dir = "testdata"' >> .air.toml; \
		echo 'tmp_dir = "tmp"' >> .air.toml; \
		echo '[build]' >> .air.toml; \
		echo '  args_bin = []' >> .air.toml; \
		echo '  bin = "./tmp/onb"' >> .air.toml; \
		echo '  cmd = "go build -o ./tmp/onb ./cmd/onb"' >> .air.toml; \
		echo '  delay = 1000' >> .air.toml; \
		echo '  exclude_dir = ["assets", "tmp", "vendor", "testdata"]' >> .air.toml; \
		echo '  exclude_file = []' >> .air.toml; \
		echo '  exclude_regex = ["_test.go"]' >> .air.toml; \
		echo '  exclude_unchanged = false' >> .air.toml; \
		echo '  follow_symlink = false' >> .air.toml; \
		echo '  full_bin = ""' >> .air.toml; \
		echo '  include_dir = []' >> .air.toml; \
		echo '  include_ext = ["go", "tpl", "tmpl", "html"]' >> .air.toml; \
		echo '  include_file = []' >> .air.toml; \
		echo '  kill_delay = "0s"' >> .air.toml; \
		echo '  log = "build-errors.log"' >> .air.toml; \
		echo '  poll = false' >> .air.toml; \
		echo '  poll_interval = 0' >> .air.toml; \
		echo '  rerun = false' >> .air.toml; \
		echo '  rerun_delay = 500' >> .air.toml; \
		echo '  send_interrupt = false' >> .air.toml; \
		echo '  stop_on_root = false' >> .air.toml; \
	fi

# Hot reload development server
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Run 'make dev-setup' first."; \
		exit 1; \
	fi

# CI target
ci: fmt vet test lint
	@echo "CI checks passed!"