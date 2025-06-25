# Keygeist Makefile

# Variables
BINARY_NAME=keygeist
BUILD_DIR=build
MAIN_PATH=cmd/operator/main.go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
.DEFAULT_GOAL := build

# Build the main application (keygeist)
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for current platform
.PHONY: build-local
build-local:
	@echo "Building $(BINARY_NAME) for local platform..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_UNIX) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_UNIX)"

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe $(MAIN_PATH)
	@echo "Build complete for all platforms"

# Build debugging tools
.PHONY: build-tools
build-tools:
	@echo "Building debugging tools..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/emulator cmd/emulator/main.go
	$(GOBUILD) -o $(BUILD_DIR)/listener cmd/listener/main.go
	@echo "Debugging tools build complete: $(BUILD_DIR)/emulator, $(BUILD_DIR)/listener"

# Build the emulator (debugging tool)
.PHONY: build-emulator
build-emulator:
	@echo "Building emulator (debugging tool)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/emulator cmd/emulator/main.go
	@echo "Build complete: $(BUILD_DIR)/emulator"

# Build the listener (debugging tool)
.PHONY: build-listener
build-listener:
	@echo "Building listener (debugging tool)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/listener cmd/listener/main.go
	@echo "Build complete: $(BUILD_DIR)/listener"

# Run the main application
.PHONY: run
run: build-local
	@echo "Running $(BINARY_NAME)..."
	@if [ -z "$$OPENAI_API_KEY" ]; then \
		echo "Error: OPENAI_API_KEY environment variable is not set"; \
		echo "Please set it with: export OPENAI_API_KEY=your_api_key_here"; \
		exit 1; \
	fi
	@if [ -z "$$OPENAI_MODEL" ]; then \
		echo "Error: OPENAI_MODEL environment variable is not set"; \
		echo "Please set it with: export OPENAI_MODEL=gpt-3.5-turbo"; \
		exit 1; \
	fi
	@echo "Using model: $$OPENAI_MODEL"
	@if [ -n "$$OPENAI_BASE_URL" ]; then \
		echo "Using custom base URL: $$OPENAI_BASE_URL"; \
	fi
	./$(BINARY_NAME)

# Run with sudo (required for uinput access)
.PHONY: run-sudo
run-sudo: build-local
	@echo "Running $(BINARY_NAME) with sudo..."
	@if [ -z "$$OPENAI_API_KEY" ]; then \
		echo "Error: OPENAI_API_KEY environment variable is not set"; \
		echo "Please set it with: export OPENAI_API_KEY=your_api_key_here"; \
		exit 1; \
	fi
	@if [ -z "$$OPENAI_MODEL" ]; then \
		echo "Error: OPENAI_MODEL environment variable is not set"; \
		echo "Please set it with: export OPENAI_MODEL=gpt-3.5-turbo"; \
		exit 1; \
	fi
	@echo "Using model: $$OPENAI_MODEL"
	@if [ -n "$$OPENAI_BASE_URL" ]; then \
		echo "Using custom base URL: $$OPENAI_BASE_URL"; \
	fi
	sudo -E ./$(BINARY_NAME)

# Run debugging tools
.PHONY: run-emulator
run-emulator: build-emulator
	@echo "Running emulator (debugging tool)..."
	sudo ./$(BUILD_DIR)/emulator

.PHONY: run-listener
run-listener: build-listener
	@echo "Running listener (debugging tool)..."
	sudo ./$(BUILD_DIR)/listener

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html
	@echo "Clean complete"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Install the binary to system
.PHONY: install
install: build-local
	@echo "Installing $(BINARY_NAME)..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Install the binary to user's home directory
.PHONY: user-install
user-install: build-local
	@echo "Installing $(BINARY_NAME) to user's home directory..."
	cp -rfv $(BINARY_NAME) ~/bin/
	@echo "Installation complete"

# Uninstall the binary from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the main keygeist to build/ directory"
	@echo "  build-local  - Build keygeist for current platform"
	@echo "  build-linux  - Build keygeist for Linux"
	@echo "  build-all    - Build keygeist for all platforms (Linux, macOS, Windows)"
	@echo "  build-tools  - Build all debugging tools (emulator, listener)"
	@echo "  build-emulator- Build the emulator debugging tool"
	@echo "  build-listener- Build the listener debugging tool"
	@echo "  run          - Build and run the keygeist"
	@echo "  run-sudo     - Build and run keygeist with sudo"
	@echo "  run-emulator - Run the emulator debugging tool"
	@echo "  run-listener - Run the listener debugging tool"
	@echo "  deps         - Install dependencies"
	@echo "  deps-update  - Update dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  install      - Install keygeist to /usr/local/bin"
	@echo "  uninstall    - Remove keygeist from /usr/local/bin"
	@echo "  help         - Show this help message"

# Development helpers
.PHONY: dev
dev: deps fmt lint test build

# Release preparation
.PHONY: release
release: clean deps fmt lint test build-all
	@echo "Release build complete"
	@ls -la $(BUILD_DIR)/

# Check if running as root (required for uinput)
.PHONY: check-root
check-root:
	@if [ "$$(id -u)" != "0" ]; then \
		echo "Warning: This application requires root privileges to access /dev/uinput"; \
		echo "Run with: sudo make run-sudo"; \
	fi 