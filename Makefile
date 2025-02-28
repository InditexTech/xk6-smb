XK6_VERSION := latest
XK6_BINARY := $(shell which xk6 || echo "")

# Targets
.PHONY: all build run test

all: test run

deps:
	@if [ -z "$(XK6_BINARY)" ]; then \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION); \
	else \
		echo "xk6 is already installed."; \
	fi

smb-server:
	@echo "Starting sftp server..."
	@docker-compose -f docker/docker-compose.yaml up -d

build: deps
	@echo "Building xk6 extension..."
	@xk6 build --with github.com/inditex/xk6-sfp=.

run: deps
	@echo "Running example..."
	@xk6 run --vus 10 --duration 1m ./examples/main.js

test: deps
	@echo "Running unit tests..."
	@go clean -testcache && go test ./...