PROJECT_VERSION := 1.0.0

GOPATH := $(shell command go env GOPATH)

XK6_VERSION := v0.13.4
XK6_BINARY := "$(GOPATH)/bin/xk6"

GOLANGCI_VERSION := v1.64.5
GOLANGCI_BINARY := "$(GOPATH)/bin/golangci-lint"

.DEFAULT_GOAL := all

.PHONY: all
all: fmt lint compose-up test run compose-down

.PHONY: deps
deps:
	@if [ ! -f "$(XK6_BINARY)" ]; then \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION); \
	else \
		echo "xk6 is already installed."; \
	fi

	@if [ ! -f "$(GOLANGCI_BINARY)" ]; then \
			echo "Installing golangci-lint..."; \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION); \
	else \
		echo "golangci-lint is already installed."; \
	fi

.PHONY: compose-up
compose-up:
	@echo "Starting smb server..."
	@docker compose -f docker/docker-compose.yaml up -d

.PHONY: compose-down
compose-down:
	@echo "Destrying smb server..."
	@docker compose -f docker/docker-compose.yaml down

.PHONY: build
build: deps
	@echo "Building xk6 extension..."
	@"$(XK6_BINARY)" build --with github.com/InditexTech/xk6-smb=.

.PHONY: run
run: deps
	@echo "Running example..."
	@"$(XK6_BINARY)" run ./examples/main.js

.PHONY: test
test: deps
	@echo "Running integration tests..."
	@go clean -testcache && go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: fmt
fmt:
	@echo "Running go fmt..."
	go fmt ./...

.PHONY: lint
lint: deps
	@echo "Running golangci-lint..."
	@"$(GOLANGCI_BINARY)" run

.PHONY: verify
verify: fmt lint compose-up test run compose-down

.PHONY: reuse-deps
reuse-deps:
	@if [ -z "reuse" ]; then \
		echo "Installing reuse tool..."; \
		pip3 install --user reuse ;\
	else \
		echo "reuse is already installed."; \
	fi

.PHONY: add-copyright-headers
reuse-annotate: reuse-deps
	@echo "Adding copyright headers..."
	@reuse annotate --copyright "Industria de Diseño Textil S.A. INDITEX" --license "Apache-2.0" --year "$$(date +%Y)" --merge-copyrights *.go
	@reuse lint

.PHONY: get-version
get-version:
	@echo $(PROJECT_VERSION)
