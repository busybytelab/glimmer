SERVICE_NAME := glimmer
PKG := github.com/busybytelab.com/$(SERVICE_NAME)
BUILD_DIR := build

# Docker registry configuration
DOCKER_REGISTRY ?= ghcr.io
DOCKER_ORG ?= busybytelab
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(SERVICE_NAME)

# Get version from git tag, default to dev
VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
# Variable to hold the Go package path for the Version variable (adjust if needed)
# VERSION_PKG = $(PKG)/cmd # Assuming a cmd package, adjust if main or other -- REMOVED
LDFLAGS := -ldflags="-s -w -X $(PKG).Version=$(VERSION)" # Use PKG directly

# Supported architectures:
# - linux/amd64: Standard x86_64 servers and desktops
# - linux/arm64: Modern ARM devices (Raspberry Pi 4/5, ARM servers, Apple Silicon M1/M2/M3/M4)
# Note: linux/arm/v7 (32-bit ARM) is not included as it's primarily for older devices
# like Raspberry Pi 3 and earlier. We can add it later if needed for legacy support.
# Note: Apple Silicon (M1/M2/M3/M4) uses arm64 architecture, so it's covered by linux/arm64.

.PHONY: all format test dep-upgrade build build-all clean help run dev-create-superuser docker-build docker-push docker-clean docker-multi-arch seed-db

all: format test build ## Run format, test, and build for current platform

format: ## Tidy modules and format code
	@echo "==> Tidying modules..."
	@go mod tidy
	@echo "==> Formatting code..."
	@go fmt ./...

test: ## Run unit tests (excluding integration tests)
	@echo "==> Running tests..."
	@go test -count=1 -v --race --timeout 30s $(shell go list ./... | grep -v /tests/) | { grep "\\(FAIL\\|panic:\\)" || test $$? = 1; }

test-local: ## Run all tests with variables
	@echo "==> Running tests locally..."
	@if [ -f .env ]; then \
	    export `cat .env` && go test -count=1 -v --race --timeout 30s $(shell go list ./... | grep -v /tests/) | { grep "\\(FAIL\\|panic:\\)" || test $$? = 1; } \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

dep-upgrade: ## Upgrade dependencies to latest versions
	@echo "==> Upgrading dependencies..."
	@go get -v -u ./...
	@go mod tidy

build: ## Build the application for the current platform
	@echo "==> Building $(SERVICE_NAME) version $(VERSION) for current platform..."
	@make build-os OS=$$(go env GOOS) ARCH=$$(go env GOARCH)

build-os: ## Build the application for the specific platform
	@echo "==> Building $(SERVICE_NAME) version $(VERSION) for $(OS) $(ARCH)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$$OS GOARCH=$$ARCH go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)-$(OS)-$(ARCH) cmd/glimmer/main.go

build-all: ## Build the application for Linux and macOS (amd64, arm64)
	@echo "==> Building $(SERVICE_NAME) version $(VERSION) for all target platforms..."
	@make build-os OS=linux ARCH=amd64
	@make build-os OS=linux ARCH=arm64
	@make build-os OS=darwin ARCH=amd64
	@make build-os OS=darwin ARCH=arm64
	@echo "==> All builds complete."

run: ## Run the application with env vars from .env file if it exists
	@echo "==> Running $(SERVICE_NAME)..."
	@if [ -f .env ]; then \
		LISTEN_ADDRESS=$$(grep LISTEN_ADDRESS .env | cut -d= -f2); \
		ENCRYPTION_KEY=$$(grep ENCRYPTION_KEY .env | cut -d= -f2 || echo ""); \
		export `cat .env` && go run cmd/glimmer/main.go --encryptionEnv=$$ENCRYPTION_KEY serve --http="$$LISTEN_ADDRESS"; \
	else \
		go run cmd/glimmer/main.go serve; \
	fi

# run go run main.go superuser -h for other sub commands
create-superuser:  ## Create a superuser for dev environment with values from .env file
	@echo "==> Creating superuser..."
	@if [ -f .env ]; then \
		EMAIL=$$(grep EMAIL .env | cut -d= -f2); \
		PASSWORD=$$(grep PASSWORD .env | cut -d= -f2); \
		ENCRYPTION_KEY=$$(grep ENCRYPTION_KEY .env | cut -d= -f2 || echo ""); \
		echo "Creating superuser with email $$EMAIL"; \
		source .env && go run cmd/glimmer/main.go --encryptionEnv=$$ENCRYPTION_KEY superuser create "$$EMAIL" "$$PASSWORD"; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

docker-build: ## Build and tag Docker image with both latest and version tags
	@echo "==> Building Docker image for $(SERVICE_NAME) version $(VERSION)..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

docker-push: docker-build ## Push Docker image to registry (requires docker login)
	@echo "==> Pushing Docker image $(DOCKER_IMAGE) version $(VERSION)..."
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@docker push $(DOCKER_IMAGE):latest

docker-multi-arch: ## Build and push multi-architecture Docker images
	@echo "==> Building multi-architecture Docker image for $(SERVICE_NAME) version $(VERSION)..."
	@docker buildx create --use --name multi-arch-builder --driver-opt network=host || true
	@docker buildx build --platform linux/arm64,linux/amd64 \
		--builder multi-arch-builder \
		--cache-from type=local,src=/tmp/.buildx-cache \
		--cache-to type=local,dest=/tmp/.buildx-cache \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		--push .

docker-clean: ## Clean up Docker images and containers
	@echo "==> Cleaning up Docker resources..."
	@docker system prune -f
	@docker images $(DOCKER_IMAGE) -q | xargs -r docker rmi
	@echo "==> Docker cleanup complete."

clean: ## Clean build artifacts and test cache
	@echo "==> Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean -testcache

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {help_text[$$1] = $$2} END {for (cmd in help_text) {printf "\033[36m%-20s\033[0m %s\n", cmd, help_text[cmd]}}'

seed-db: ## Seed the database with test data
	@echo "==> Seeding database with test data using YAML configuration..."
	@go run cmd/seed/main.go
