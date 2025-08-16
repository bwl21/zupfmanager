# Makefile
BINARY_NAME := zupfmanager
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
MODULE := $(shell grep "^module " go.mod | awk '{print $$2}')
GO_LDFLAGS := -X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.GitCommit=$(COMMIT)

# Standardziel
all: build

# Frontend build
frontend:
	@echo "Building frontend..."
	@cd frontend && npm ci && npm run build
	@echo "Frontend build complete: frontend/dist/"

# Frontend development dependencies
frontend-deps:
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Frontend dependencies installed."

# install completions
completion:
	@echo "building completions"	
	@zupfmanager completion > ~/.zsh/completions/_zupfmanager
	@echo please execute 'source ~/.zshrc to activate the completions'

# Build für aktuelles System (mit Frontend)
build: frontend
	@echo "Building $(BINARY_NAME) für $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p dist
	@go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)
	@cp -r frontend/dist dist/frontend
	@echo "Build abgeschlossen: dist/$(BINARY_NAME) $(VERSION) $(COMMIT)"
	@echo "Frontend included: dist/frontend/"

# Build nur Backend (ohne Frontend)
build-backend:
	@echo "Building $(BINARY_NAME) backend only für $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p dist
	@go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)
	@echo "Backend build abgeschlossen: dist/$(BINARY_NAME) $(VERSION) $(COMMIT)"

# Build for Linux (amd64) with frontend
build-linux: frontend
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64
	@cp -r frontend/dist dist/frontend-linux
	@echo "Build complete: dist/$(BINARY_NAME)-linux-amd64 with frontend"

# Build for macOS (amd64) with frontend
build-macos-amd64: frontend
	@echo "Building $(BINARY_NAME) for darwin/amd64..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64
	@cp -r frontend/dist dist/frontend-darwin-amd64
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-amd64 with frontend"

# Build for macOS (arm64 - Apple Silicon) with frontend
build-macos-arm64: frontend
	@echo "Building $(BINARY_NAME) for darwin/arm64..."
	@mkdir -p dist
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64
	@cp -r frontend/dist dist/frontend-darwin-arm64
	@echo "Build complete: dist/$(BINARY_NAME)-darwin-arm64 with frontend"

# Build for Windows with frontend
build-windows: frontend
	@echo "Building $(BINARY_NAME) for windows/amd64..."
	@mkdir -p dist
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -ldflags "$(GO_LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe
	@cp -r frontend/dist dist/frontend-windows
	@echo "Build complete: dist/$(BINARY_NAME)-windows-amd64.exe with frontend"

# Build for both macOS architectures
build-macos: build-macos-amd64 build-macos-arm64

# Build for both macOS and Windows
build-all: build-macos build-windows

# Create release packages
package-linux: build-linux
	@echo "Creating Linux package..."
	@cd dist && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 frontend-linux/
	@echo "Linux package created: dist/$(BINARY_NAME)-linux-amd64.tar.gz"

package-macos-amd64: build-macos-amd64
	@echo "Creating macOS Intel package..."
	@cd dist && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 frontend-darwin-amd64/
	@echo "macOS Intel package created: dist/$(BINARY_NAME)-darwin-amd64.tar.gz"

package-macos-arm64: build-macos-arm64
	@echo "Creating macOS Apple Silicon package..."
	@cd dist && tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 frontend-darwin-arm64/
	@echo "macOS Apple Silicon package created: dist/$(BINARY_NAME)-darwin-arm64.tar.gz"

package-windows: build-windows
	@echo "Creating Windows package..."
	@cd dist && zip -r $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe frontend-windows/
	@echo "Windows package created: dist/$(BINARY_NAME)-windows-amd64.zip"

# Package all platforms
package-all: package-linux package-macos-amd64 package-macos-arm64 package-windows

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/
	@rm -rf frontend/dist/
	@echo "Clean complete."

# Development server
dev:
	@echo "Starting development server..."
	@go run . api --port 8080 --frontend frontend/dist

# Phony targets
.PHONY: all build build-backend frontend frontend-deps build-linux build-macos-amd64 build-macos-arm64 build-macos build-windows build-all package-linux package-macos-amd64 package-macos-arm64 package-windows package-all clean dev completion
