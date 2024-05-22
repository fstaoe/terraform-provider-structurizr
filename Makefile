#!/bin/sh

STRUCTURIZR_CLI_VER=v2024.03.03
STRUCTURIZR_CLI_DIR=internal/client/cli/tools/structurizr-cli

default: help

.PHONY: test
test: ## Run unit tests
	go test -race -cover ./... -v $(TESTARGS) -timeout 10m

.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test -race -cover ./internal/provider -v $(TESTARGS) -timeout 10m

.PHONY: deps
deps: ## Downloads all required dependencies
	@rm -rf ${STRUCTURIZR_CLI_DIR}
	@./scripts/download_structurizr_cli.sh "${STRUCTURIZR_CLI_VER}" "${STRUCTURIZR_CLI_DIR}"

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
