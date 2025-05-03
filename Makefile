SERVICE_NAME := glimmer
PKG := github.com/busybytelab.com/$(SERVICE_NAME)
BUILD_DIR := build

# Get version from git tag, default to dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
# Variable to hold the Go package path for the Version variable (adjust if needed)
# VERSION_PKG = $(PKG)/cmd # Assuming a cmd package, adjust if main or other -- REMOVED
LDFLAGS := -ldflags="-s -w -X $(PKG).Version=$(VERSION)" # Use PKG directly

.PHONY: all format test dep-upgrade build build-all clean help run dev-create-superuser

all: format test build ## Run format, test, and build for current platform

format: ## Tidy modules and format code
	@echo "==> Tidying modules..."
	@go mod tidy
	@echo "==> Formatting code..."
	@go fmt ./...

test: ## Run unit tests (excluding integration tests)
	@echo "==> Running tests..."
	@go test -count=1 -v --race --timeout 30s $(shell go list ./... | grep -v /tests/) | { grep "\\(FAIL\\|panic:\\)" || test $$? = 1; }

dep-upgrade: ## Upgrade dependencies to latest versions
	@echo "==> Upgrading dependencies..."
	@go get -v -u ./...
	@go mod tidy

build: ## Build the application for the current platform
	@echo "==> Building $(SERVICE_NAME) version $(VERSION) for current platform..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)

build-all: ## Build the application for Linux and macOS (amd64, arm64)
	@echo "==> Building $(SERVICE_NAME) version $(VERSION) for all target platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "--> Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)-linux-amd64
	@echo "--> Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)-linux-arm64
	@echo "--> Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)-darwin-amd64
	@echo "--> Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME)-darwin-arm64
	@echo "==> All builds complete."

run: ## Run the application with env vars from .env file if it exists
	@echo "==> Running $(SERVICE_NAME)..."
	@if [ -f .env ]; then \
		LISTEN_ADDRESS=$$(grep LISTEN_ADDRESS .env | cut -d= -f2); \
		ENCRYPTION_KEY=$$(grep ENCRYPTION_KEY .env | cut -d= -f2 || echo ""); \
		go run main.go serve --http="$$LISTEN_ADDRESS" --encryptionEnv=$$ENCRYPTION_KEY; \
	else \
		go run ./main.go serve; \
	fi

# run go run main.go superuser -h for other sub commands
create-development-superuser:  ## Create a superuser for development with env vars from .env file
	@echo "==> Creating superuser..."
	@if [ -f .env ]; then \
		EMAIL=$$(grep EMAIL .env | cut -d= -f2); \
		PASSWORD=$$(grep PASSWORD .env | cut -d= -f2); \
		ENCRYPTION_KEY=$$(grep ENCRYPTION_KEY .env | cut -d= -f2 || echo ""); \
		echo "Creating superuser with email $$EMAIL"; \
		go run main.go superuser create "$$EMAIL" "$$PASSWORD" --encryptionEnv=$$ENCRYPTION_KEY; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

clean: ## Clean build artifacts and test cache
	@echo "==> Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean -testcache

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {help_text[$$1] = $$2} END {for (cmd in help_text) {printf "\033[36m%-20s\033[0m %s\n", cmd, help_text[cmd]}}'
