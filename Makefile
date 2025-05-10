# Output binary name
BINARY_NAME=zupfmanager


# Makefile
BINARY_NAME := zupfmanager
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
MODULE := $(shell grep "^module " go.mod | awk '{print $$2}')
GO_LDFLAGS := -X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.GitCommit=$(COMMIT)

# Standardziel
all: build

# install completions
completion:
	@echo "building completions"	
	@zupfmanager completion > ~/.zsh/completions/_zupfmanager
	@echo please execute 'source ~/.zshrc to activate the completions'

# Build für aktuelles System
build:
	@echo "Building $(BINARY_NAME) für $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)
	@echo "Build abgeschlossen: dist/$(BINARY_NAME) $(VERSION) $(COMMIT)"

# Build for Linux (amd64)
build-linux:
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o dist/$(BINARY_NAME)-linux-amd64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-linux-amd64"

# Build for macOS (amd64)
build-macos-amd64:
	@echo "Building $(BINARY_NAME) for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-amd64"

# Build for macOS (arm64 - Apple Silicon)
build-macos-arm64:
	@echo "Building $(BINARY_NAME) for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-arm64"

build-windows:
	@echo "Building $(BINARY_NAME) for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o dist/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Build complete: dist/$(BINARY_NAME)-windows-amd64.exe"

# Build for both macOS architectures
build-macos: build-macos-amd64 build-macos-arm64

# Build for both macOS and Windows
build-all: build-macos build-windows

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f dist/$(BINARY_NAME)*
	@echo "Clean complete."

# Phony targets
.PHONY: all build build-linux build-macos-amd64 build-macos-arm64 build-macos clean 