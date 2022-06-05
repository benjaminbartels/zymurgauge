GOCMD=go
GOTEST=$(GOCMD) test
MAIN_DIR=cmd/zym
BINARY_NAME=zym
VERSION?=develop
SERVICE_PORT?=8080

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

build-go: ## Build the default go package and put the output binary in out/bin/
	mkdir -p out/bin
	GOOS=linux GOARCH=arm CGO_ENABLED=0 $(GOCMD) build -a \
	-ldflags="-w -s -extldflags '-static' -X 'main.version=$(VERSION)'" -o out/bin/$(BINARY_NAME) ./$(MAIN_DIR)

build-react: ## Build the React UI
	yarn --cwd "ui" add react-scripts@5.0.0 --network-timeout 100000
	yarn --cwd "ui" build --network-timeout 100000

build-docker: ## Use the dockerfile to build the container
	DOCKER_BUILDKIT=1 docker build -t $(BINARY_NAME) -f build/Dockerfile --target production .

tidy: ## Add missing and remove unused modules
	$(GOCMD) mod tidy

init: ## Initialize folders and database
	mkdir -p ui/build && touch ui/build/.gitkeep
	mkdir -p tmp
	ZYM_DBPATH=tmp/zymurgaugedb go run cmd/zym/main.go init --username=admin --password=password \
		--brewfather-user-id=$(BREWFATHER_USERID) --brewfather-key=$(BREWFATHER_KEY)

watch-go: ## Run the go service with Air to have automatic reload on changes
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run -it --rm \
	-p $(SERVICE_PORT):$(SERVICE_PORT) \
	-v $(shell pwd):/$(PACKAGE_NAME) \
	-w /$(PACKAGE_NAME) \
	cosmtrek/air -c ./.air.conf

watch-react: ## Run the React UI with automatic reload on changes
	yarn --cwd "ui" start

clean: ## Remove build and coverage related file
	rm -fr out tmp
	rm -f  profile.cov
	rm -fr ui/build

test: ## Run the tests of the project
	# Ensure that ui/build has somthing in it so tests will work
	mkdir -p ui/build && touch ui/build/.gitkeep
	$(GOTEST) -v ./... -race

coverage: ## Generate code coverge report
	# Ensure that ui/build has somtehing in it so tests will work
	mkdir -p ui/build && touch ui/build/.gitkeep
	$(GOTEST) -v -covermode=atomic -coverpkg=./... -coverprofile=profile.cov  ./...
	$(GOCMD) tool cover -func profile.cov

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