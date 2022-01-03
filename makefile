GOCMD=go
GOTEST=$(GOCMD) test
MAIN_DIR=cmd/zym
BINARY_NAME=zym
#TODO: See https://dev.to/eugenebabichenko/generating-pretty-version-strings-including-nightly-with-git-and-makefiles-48p3
VERSION?=0.0.0
SERVICE_PORT?=8080
DELVE_PORT?=2345
# If set it should end with '/'
DOCKER_REGISTRY?=
# Set EXPORT_RESULT to true isn used in CI
EXPORT_RESULT?=false

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
	GOOS=linux GOARCH=arm CGO_ENABLED=0 $(GOCMD) build -a -ldflags="-w -s -extldflags '-static'" -o out/bin/$(BINARY_NAME) ./$(MAIN_DIR)

clean: ## Remove build and coverage related file
	rm -fr out $(MAIN_DIR)/tmp
	rm -f junit-report.xml checkstyle-report.xml profile.cov coverage.xml yamllint-checkstyle.xml
	rm -f zymurgaugedb $(MAIN_DIR)/zymurgaugedb

tidy: ## Add missing and remove unused modules
	$(GOCMD) mod tidy

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod vendor

watch: ## Run the code with Air to have automatic reload on changes
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm \
	-p $(SERVICE_PORT):$(SERVICE_PORT) \
	--env-file=configs/dev.env \
	-w /$(PACKAGE_NAME)/$(MAIN_DIR) \
	-v $(shell pwd):/$(PACKAGE_NAME) \
	--mount source=zym_data,target=/$(PACKAGE_NAME)/$(MAIN_DIR)/data \
	cosmtrek/air

debug: ## Run the code with Delve to debug
	DOCKER_BUILDKIT=1 docker build -t $(BINARY_NAME)-debug --target debugger  .
	docker run -it --rm \
	-p $(SERVICE_PORT):$(SERVICE_PORT) -p $(DELVE_PORT):$(DELVE_PORT) \
	--env-file=configs/dev.env \
	--mount source=zym_data,target=/src/data \
	$(BINARY_NAME)-debug

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
lint: lint-go lint-yaml  ## Run all linters

# DOCKER NOT WORKING
# lint-dockerfile: ## Lint your Dockerfile
# ifeq ($(shell test -e ./Dockerfile && echo yes),yes)
# 	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
# 	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
# 	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
# 	docker run --rm -i $(CONFIG_OPTION) hadolint/hadolint hadolint $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
# endif

lint-go: ## Lint go files
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s $(OUTPUT_OPTIONS)

lint-yaml: ## Lint yaml files
ifeq ($(EXPORT_RESULT),true)
	go install github.com/thomaspoignant/yamllint-checkstyle@latest
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | yamllint-checkstyle > yamllint-checkstyle.xml)
endif
	docker run --rm -it -v $(shell pwd):/data cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml') $(OUTPUT_OPTIONS)

# DOCKER NOT WORKING
# ## Docker:
# docker-build: ## Use the dockerfile to build the container
# 	DOCKER_BUILDKIT=1 docker build -t $(BINARY_NAME) --target production .

# docker-release: ## Release the container with tag latest and version
# 	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
# 	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
# 	# Push the docker images
# 	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
# 	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

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