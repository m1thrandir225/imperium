# Go + WebUI Dev Makefile for apps/client

SHELL := /bin/zsh

# Paths
WEBUI_DIR := ./webui
BIN_DIR := ./bin
TMP_DIR := ./tmp
WEBUI_BUILD_DIR := ./internal/webui

# Binaries
AIR := $(HOME)/go/bin/air
GO := go
PNPM := pnpm

.DEFAULT_GOAL := help

.PHONY: help deps dev dev-go dev-web build build-go build-web clean run-prod

deps:
	@if ! command -v $(PNPM) >/dev/null 2>&1; then \
		echo "pnpm not found. Install with: corepack enable && corepack prepare pnpm@latest --activate"; \
		exit 1; \
	fi
	@echo "Installing Air (Go hot reload)..."
	@$(GO) install github.com/air-verse/air@latest
	@echo "Installing webui deps..."
	@cd $(WEBUI_DIR) && $(PNPM) install

dev: ## Run Go and WebUI dev servers in parallel
	@$(MAKE) -j 2 dev-go dev-web

dev-go: ## Run Go app with Air
	@if [ ! -x "$(AIR)" ]; then \
		echo "Air is not installed at $(AIR). Run: make deps"; \
		exit 1; \
	fi
	@$(AIR)

dev-web: ## Run React WebUI (Vite)
	@cd $(WEBUI_DIR) && $(PNPM) dev

build: build-web build-go ## Build webui first, then Go (webui is embedded in Go)

build-go: build-web ## Build Go binary (depends on webui build)
	@echo "Building Go application..."
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/client ./cmd

build-web: ## Build React WebUI
	@echo "Building WebUI..."
	@cd $(WEBUI_DIR) && $(PNPM) install --frozen-lockfile
	@cd $(WEBUI_DIR) && $(PNPM) build
	@echo "WebUI built to $(WEBUI_BUILD_DIR)"

# Production build with optimizations
build-prod: build-web build-go-prod ## Production build with optimizations

build-go-prod: build-web ## Build Go binary with production optimizations
	@echo "Building Go application for production..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GIN_MODE=release $(GO) build -ldflags="-w -s" -o $(BIN_DIR)/client ./cmd

# Check if webui build exists
check-webui:
	@if [ ! -d "$(WEBUI_BUILD_DIR)" ] || [ -z "$$(ls -A $(WEBUI_BUILD_DIR) 2>/dev/null)" ]; then \
		echo "WebUI build not found or empty. Building..."; \
		$(MAKE) build-web; \
	else \
		echo "WebUI build exists at $(WEBUI_BUILD_DIR)"; \
	fi

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) $(TMP_DIR)
	@rm -rf $(WEBUI_BUILD_DIR)
	@echo "Clean complete"

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

run-prod: build-prod ## Build for production and run the application
	@echo "Starting application (production build)..."
	@$(BIN_DIR)/client

mock:
	#Services
	mockgen -package mockservices -destination internal/services/mocks/auth_service.go github.com/m1thrandir225/imperium/apps/client/internal/services/ AuthService
	mockgen -package mockservices -destination internal/services/mocks/client_service.go github.com/m1thrandir225/imperium/apps/client/internal/services/ ClientService
	mockgen -package mockservices -destination internal/services/mocks/host_service.go github.com/m1thrandir225/imperium/apps/client/internal/services/ HostService
	mockgen -package mockservices -destination internal/services/mocks/session_service.go github.com/m1thrandir225/imperium/apps/client/internal/services/ SessionService

	# Server Handlers
	mockgen -package mockserver -destination internal/server/mocks/auth_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ AuthHandler
	mockgen -package mockserver -destination internal/server/mocks/client_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ ClientHandler
	mockgen -package mockserver -destination internal/server/mocks/host_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ HostHandler
	mockgen -package mockserver -destination internal/server/mocks/session_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ SessionHandler
	mockgen -package mockserver -destination internal/server/mocks/status_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ StatusHandler
	mockgen -package mockserver -destination internal/server/mocks/config_handler.go github.com/m1thrandir225/imperium/apps/client/internal/server/ ConfigHandler
