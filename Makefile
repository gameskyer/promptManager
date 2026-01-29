# Makefile for PromptMaster

.PHONY: all build dev clean install test frontend backend

# Default target
all: build

# Install dependencies
install:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Development mode
dev:
	wails dev

# Build for current platform
build:
	wails build

# Build for Windows
build-windows:
	wails build -platform windows/amd64

# Build for macOS
build-macos:
	wails build -platform darwin/universal

# Build for Linux
build-linux:
	wails build -platform linux/amd64

# Build all platforms
build-all: build-windows build-macos build-linux

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf frontend/dist/
	cd frontend && rm -rf node_modules/

# Run tests
test:
	go test ./...

# Run frontend lint
lint-frontend:
	cd frontend && npm run lint

# Run Go vet
vet:
	go vet ./...

# Update Wails bindings
bindings:
	wails generate module

# Package the application
package: build
	@echo "Packaging application..."
	mkdir -p dist
	cp build/bin/PromptMaster* dist/ 2>/dev/null || cp build/bin/promptmaster* dist/

# Help
help:
	@echo "Available targets:"
	@echo "  make install      - Install all dependencies"
	@echo "  make dev          - Run in development mode"
	@echo "  make build        - Build for current platform"
	@echo "  make build-all    - Build for all platforms"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make bindings     - Update Wails bindings"
