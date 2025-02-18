BINPATH ?= build

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

MAIN=dp-search-scrubber-api
BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all audit build build-bin clean convey debug delimiter-% fmt lint run run-container test test-all test-component update help

all: delimiter-AUDIT audit delimiter-LINTERS lint delimiter-UNIT-TESTS test delimiter-COMPONENT_TESTS test-component delimiter-FINISH ## Runs multiple targets, audit, lint, test and test-component

audit: ## Audits and finds vulnerable dependencies
	go list -json -m all | nancy sleuth

build: Dockerfile ## Builds ./Dockerfile image name: scrubber
	docker build -t scrubber .

build-bin: ## builds bin
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build $(LDFLAGS) -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN)

clean: ## Removes /bin folder
	rm -fr ./build
	rm -fr ./vendor
	
debug: ## Runs the api locally in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-search-scrubber-api
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-search-scrubber-api

delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

fmt: ## Formats the code using go fmt and go vet
	go fmt ./...
	go vet ./...

validate-specification: ## Quality checking of your OpenAPI spec
	redocly lint swagger.yaml

## Automated checking of your source code for programmatic and stylistic errors
lint: validate-specification
	golangci-lint run ./...

run: ## Run the app locally
	go run . 

run-container: build ## First builds ./Dockerfile with image name: scrubber and then runs a container, with name: scrubber_container, on port :28700 
	docker run -p :28700:28700 --name scrubber_container -ti --rm scrubber
 
test: ## Runs standard unit test tests
	go test -race -cover ./... 

test-all: test-component test ## Runs all tests with -race and -cover flags
	go test -race -cover ./...

test-component: ## Runs component tests
	go test -cover -coverpkg=github.com/ONSdigital/dp-search-scrubber-api/... -component

convey: ## Runs Convey tests
	goconvey ./...

update: ## Go gets all of the dependencies and downloads them
	go get .
	go mod download

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
