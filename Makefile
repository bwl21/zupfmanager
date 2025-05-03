# Output binary name
BINARY_NAME=zupfmanager

# Default target: build for the current OS/Arch
all: build-linux build-macos build-windows

# Build for the current system
build:
	@echo "Building $(BINARY_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@go build -o dist/$(BINARY_NAME) main.go
	@echo "Build complete: dist/$(BINARY_NAME)"

# Build for Linux (amd64)
build-linux:
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/$(BINARY_NAME)-linux-amd64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-linux-amd64"

# Build for macOS (amd64)
build-macos-amd64:
	@echo "Building $(BINARY_NAME) for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-amd64"

# Build for macOS (arm64 - Apple Silicon)
build-macos-arm64:
	@echo "Building $(BINARY_NAME) for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-arm64"

build-windows:
	@echo "Building $(BINARY_NAME) for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/$(BINARY_NAME)-windows-amd64.exe main.go
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