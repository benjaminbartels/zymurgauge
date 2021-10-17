include ./config/dev.env
export $(shell sed 's/=.*//' ./config/dev.env)

DEFAULT_PKG_DIR		=	./cmd/zym
DEFAULT_BINARY_NAME	=	zym

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
build: ## Build the default package and put the output binary in out/bin/
	mkdir -p out/bin
	go build -o out/bin/$(DEFAULT_BINARY_NAME) $(DEFAULT_PKG_DIR)

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out

run: ## Run the default package main
	go run $(DEFAULT_PKG_DIR)

## Test:
test: ## Run the tests of the project
	go test -v ./... -race

## Lint:
lint: ## Use golintci-lint on your project
	golangci-lint run
	
## Codegen:
generate-mocks: ## Generate mocks
	mockery --all --output ./internal/test/mocks 

## Help:
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)