include ./config/dev.env
export $(shell sed 's/=.*//' ./config/dev.env)

GOCMD=go
GOTEST=$(GOCMD) test
DEFAULT_PKG_DIR		=	./cmd/zym
DEFAULT_BINARY_NAME	=	zym
VERSION := 0.1 # $(shell git rev-parse --short HEAD)
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

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
	$(GOCMD) build -o out/bin/$(DEFAULT_BINARY_NAME) $(DEFAULT_PKG_DIR)

build-linux-arm: ## Build the default package for linux and put the output binary in out/bin/
	mkdir -p out/bin
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -a -ldflags="-w -extldflags '-static'" -o out/bin/$(DEFAULT_BINARY_NAME) $(DEFAULT_PKG_DIR)

clean: ## Remove build and coverage related file
	rm -fr ./out
	rm -f ./junit-report.xml checkstyle-report.xml ./profile.cov ./coverage.xml yamllint-checkstyle.xml

clean-data: ## Remove local database file
	rm -fr ./zymurgaugedb

tidy: ## Add missing and remove unused modules
	$(GOCMD) mod tidy

run: ## Run the default package main
	$(GOCMD) run $(DEFAULT_PKG_DIR)

## Test:
test: ## Run the tests of the project
	$(GOTEST) -v ./... -race

coverage: ## Generate code coverge report
	$(GOTEST) -v -covermode=atomic -coverpkg=./... -coverprofile=profile.cov  ./...
	$(GOCMD) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	go install github.com/AlekSi/gocov-xml@latest
	go install github.com/axw/gocov/gocov@latest
	gocov convert profile.cov | gocov-xml > coverage.xml
endif	

## Lint:
lint: lint-go lint-yaml ## Run all linters

lint-go: ## Lint go files
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

lint-yaml: ## Lint yaml files
ifeq ($(EXPORT_RESULT), true)
	go install github.com/thomaspoignant/yamllint-checkstyle@latest
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	docker run --rm -it -v $(shell pwd):/data cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)

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