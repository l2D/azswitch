.PHONY: build test lint clean install run fmt tidy coverage setup pre-commit

# Build variables
BINARY_NAME=azswitch
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-s -w -X github.com/l2D/azswitch/internal/version.Version=$(VERSION) -X github.com/l2D/azswitch/internal/version.CommitSHA=$(COMMIT_SHA) -X github.com/l2D/azswitch/internal/version.BuildTime=$(BUILD_TIME)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

# Default target
all: lint test build

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/azswitch

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/azswitch
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/azswitch
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/azswitch
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/azswitch
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/azswitch

# Run tests
test:
	$(GOTEST) -v -race ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GOFMT) ./...

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Install the binary
install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

# Development: run with live reload (requires air)
dev:
	air

# Setup pre-commit hooks
setup:
	pre-commit install

# Run pre-commit on all files
pre-commit:
	pre-commit run --all-files

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  build-all  - Build for all platforms"
	@echo "  test       - Run tests"
	@echo "  coverage   - Run tests with coverage report"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  tidy       - Tidy dependencies"
	@echo "  install    - Install binary to GOPATH/bin"
	@echo "  run        - Build and run"
	@echo "  clean      - Remove build artifacts"
	@echo "  dev        - Run with live reload (requires air)"
	@echo "  setup      - Install pre-commit hooks"
	@echo "  pre-commit - Run pre-commit on all files"
