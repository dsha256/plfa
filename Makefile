.DEFAULT_GOAL = help

MAIN_API = "cmd/api/main.go"
MAIN_API_BIN = "bin/api"

CURRENT_TIME = $(shell date --iso-8601=seconds)

build: ## Builds the cmd/api application.
	@echo 'Building cmd/api...'
	go build -ldflags='-s -X main.buildTime=${CURRENT_TIME}' -o=$(MAIN_API_BIN) $(MAIN_API)

unit-tests: ## Run unit tests in verbose mode.
	@echo 'Starting tests...'
	go test -v -cover -race ./...

swag-gen: ## Generate swagger files.
	@echo 'Generating swagger files...'
	swag init -q -g $(MAIN_API) -o docs/swagger

help: ## Prints this message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build unit-tests swag-gen help
