# Go + WebUI Dev Makefile for apps/client

SHELL := /bin/zsh

# Paths
WEBUI_DIR := ./webui
BIN_DIR := ./bin
TMP_DIR := ./tmp

# Binaries
AIR := $(HOME)/go/bin/air
GO := go
PNPM := pnpm

.PHONY: deps dev dev-go dev-web build build-go build-web clean

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

build: build-go build-web ## Build both

build-go:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/client ./cmd

build-web:
	@cd $(WEBUI_DIR) && $(PNPM) build

clean:
	@rm -rf $(BIN_DIR) $(TMP_DIR)
