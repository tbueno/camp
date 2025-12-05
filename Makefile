.PHONY: help build test test-unit test-integration test-integration-docker clean install fmt lint

# Default target
help:
	@echo "Camp - Development Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build                 - Build the camp binary"
	@echo "  make test                  - Run all tests (unit + integration)"
	@echo "  make test-unit             - Run unit tests only"
	@echo "  make test-integration      - Run integration tests"
	@echo "  make test-integration-docker - Run integration tests in Docker only"
	@echo "  make test-bootstrap        - Run bootstrap integration test"
	@echo "  make test-rebuild          - Run rebuild integration test"
	@echo "  make test-packages         - Run packages integration test"
	@echo "  make test-flakes           - Run flakes integration test"
	@echo "  make test-nuke             - Run nuke integration test"
	@echo "  make clean                 - Clean build artifacts and test files"
	@echo "  make install               - Install camp to /usr/local/bin"
	@echo "  make fmt                   - Format Go code"
	@echo "  make lint                  - Run linters"
	@echo ""

# Build the camp binary
build:
	@echo "Building camp..."
	go build -o camp main.go
	@echo "✓ Build complete: ./camp"

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.txt ./...
	@echo "✓ Unit tests complete"

# Run all integration tests
test-integration:
	@echo "Building camp for Linux (for Docker tests)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "✓ Build complete: ./camp (linux/amd64)"
	@echo "Running integration tests..."
	./test/integration/scripts/run-tests.sh

# Run integration tests in Docker only (skip macOS tests)
test-integration-docker:
	@echo "Building camp for Linux (for Docker tests)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "✓ Build complete: ./camp (linux/amd64)"
	@echo "Building Docker image..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	@echo "Running Docker integration tests..."
	./test/integration/scripts/run-tests.sh

# Run individual integration tests
test-bootstrap:
	@echo "Building camp for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "Running bootstrap test..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	docker run --rm \
		-v "$(PWD)/camp:/home/testuser/bin/camp:ro" \
		-v "$(PWD)/test/integration/scripts:/home/testuser/tests:ro" \
		-v "$(PWD)/test/integration/fixtures:/home/testuser/fixtures:ro" \
		-v "$(PWD)/templates:/home/testuser/templates:ro" \
		camp-integration-test:latest \
		/bin/bash /home/testuser/tests/test-bootstrap.sh

test-rebuild:
	@echo "Building camp for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "Running rebuild test..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	docker run --rm \
		-v "$(PWD)/camp:/home/testuser/bin/camp:ro" \
		-v "$(PWD)/test/integration/scripts:/home/testuser/tests:ro" \
		-v "$(PWD)/test/integration/fixtures:/home/testuser/fixtures:ro" \
		-v "$(PWD)/templates:/home/testuser/templates:ro" \
		camp-integration-test:latest \
		/bin/bash /home/testuser/tests/test-rebuild.sh

test-packages:
	@echo "Building camp for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "Running packages test..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	docker run --rm \
		-v "$(PWD)/camp:/home/testuser/bin/camp:ro" \
		-v "$(PWD)/test/integration/scripts:/home/testuser/tests:ro" \
		-v "$(PWD)/test/integration/fixtures:/home/testuser/fixtures:ro" \
		-v "$(PWD)/templates:/home/testuser/templates:ro" \
		camp-integration-test:latest \
		/bin/bash /home/testuser/tests/test-packages.sh

test-flakes:
	@echo "Building camp for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "Running flakes test..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	docker run --rm \
		-v "$(PWD)/camp:/home/testuser/bin/camp:ro" \
		-v "$(PWD)/test/integration/scripts:/home/testuser/tests:ro" \
		-v "$(PWD)/test/integration/fixtures:/home/testuser/fixtures:ro" \
		-v "$(PWD)/templates:/home/testuser/templates:ro" \
		camp-integration-test:latest \
		/bin/bash /home/testuser/tests/test-flakes.sh

test-nuke:
	@echo "Building camp for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o camp main.go
	@echo "Running nuke test..."
	cd test/integration/docker && docker build -t camp-integration-test:latest .
	docker run --rm \
		-v "$(PWD)/camp:/home/testuser/bin/camp:ro" \
		-v "$(PWD)/test/integration/scripts:/home/testuser/tests:ro" \
		-v "$(PWD)/test/integration/fixtures:/home/testuser/fixtures:ro" \
		-v "$(PWD)/templates:/home/testuser/templates:ro" \
		camp-integration-test:latest \
		/bin/bash /home/testuser/tests/test-nuke.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f camp
	rm -f coverage.txt
	rm -rf test-results/
	docker rmi -f camp-integration-test:latest 2>/dev/null || true
	@echo "✓ Clean complete"

# Install camp to system
install: build
	@echo "Installing camp to /usr/local/bin..."
	sudo cp camp /usr/local/bin/camp
	sudo chmod +x /usr/local/bin/camp
	@echo "✓ Camp installed to /usr/local/bin/camp"

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Run linters (requires golangci-lint)
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install from: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	@echo "✓ Linting complete"

# Quick development workflow
dev: fmt build test-unit
	@echo "✓ Development workflow complete"
