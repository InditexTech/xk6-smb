XK6_VERSION := v0.13.4
XK6_BINARY := $(shell command -v xk6 2> /dev/null)

GOLANGCI_VERSION := v1.64.5
GOLANGCI_BINARY := $(shell command -v golangci-lint 2> /dev/null)

.PHONY: all
all: fmt lint compose-up test run compose-down

.PHONY: deps
deps:
	@if [ -z "$(XK6_BINARY)" ]; then \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION); \
	else \
		echo "xk6 is already installed."; \
	fi

.PHONY: compose-up
compose-up:
	@echo "Starting smb server..."
	@docker-compose -f docker/docker-compose.yaml up -d

.PHONY: compose-down
compose-down:
	@echo "Destrying smb server..."
	@docker-compose -f docker/docker-compose.yaml down

.PHONY: build
build: deps
	@echo "Building xk6 extension..."
	@xk6 build --with github.com/inditex/xk6-sfp=.

.PHONY: run
run: deps
	@echo "Running example..."
	@xk6 run ./examples/main.js

.PHONY: test
test: deps
	@echo "Running integration tests..."
	@go clean -testcache && go test ./...

.PHONY: fmt
fmt:
	@echo "Running go fmt..."
	go fmt ./...

.PHONY: lint
lint: deps
	@echo "Running golangci-lint..."
	@golangci-lint run

.PHONY: verify
verify: fmt lint test run
	@echo "Running verify..."
