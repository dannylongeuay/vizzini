.DEFAULT_GOAL := help

.PHONY: help
help: ## View help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/vizzini ./src/...

.PHONY: tests-all
tests-all: ## Run all tests
	go test -v -count=1 ./src/...

.PHONY: tests-short
tests-short: ## Run short tests
	go test -v -count=1 -short ./src/...

.PHONY: lint
lint: ## Run linter
	golangci-lint run --fast

.PHONY: format
format: ## Run formatter 
	go fmt ./src/...

.PHONY: cq
cq: format lint tests-short ## Run code quality tools

.PHONY: cq-check
cq-check: lint tests-short ## Run code quality check

.PHONY: run
run: build ## Run binary
	./bin/vizzini

