# Git-Graft Makefile

BINARY_NAME=graft
BUILD_DIR=build
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
INSTALL_DIR?=$(shell go env GOPATH 2>/dev/null || echo "$(HOME)/go")/bin

.PHONY: all build clean install run test lint fmt help

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/graft

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/graft
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/graft
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/graft
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/graft
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/graft

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	go clean

# Install to GOPATH/bin (or ~/go/bin)
install: build
	@mkdir -p $(INSTALL_DIR)
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Done! Make sure $(INSTALL_DIR) is in your PATH"

# Install to /usr/local/bin (requires sudo)
install-global: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)

# Run the application
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run without building
run-dev:
	go run ./cmd/graft

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	@which goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Check for updates
check-updates:
	@which go-mod-outdated > /dev/null || go install github.com/psampaz/go-mod-outdated@latest
	go list -u -m -json all | go-mod-outdated -direct

# Development mode with auto-reload
dev:
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	air

# Show help
help:
	@echo "Git-Graft Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build         Build the binary"
	@echo "  make build-all     Build for all platforms"
	@echo "  make clean         Clean build artifacts"
	@echo "  make install       Install to GOPATH/bin"
	@echo "  make install-global Install to /usr/local/bin (sudo)"
	@echo "  make uninstall     Uninstall from GOPATH/bin"
	@echo "  make run           Build and run"
	@echo "  make run-dev       Run without building"
	@echo "  make test          Run tests"
	@echo "  make test-coverage Run tests with coverage"
	@echo "  make lint          Run linter"
	@echo "  make fmt           Format code"
	@echo "  make tidy          Tidy dependencies"
	@echo "  make deps          Download dependencies"
	@echo "  make dev           Development mode with auto-reload"
	@echo "  make help          Show this help"
