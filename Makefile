BINPATH ?= build

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all ## runs audit, test and build commands
all: audit test build

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth

.PHONY: build 
build: Dockerfile ## Builds ./Dockerfile image name: scrubber
	docker build -t scrubber .

.PHONY: lint ## Formats the code using go fmt and go vet
lint: 
	golangci-lint run ./...

.PHONY: run
run: build ## First builds ./Dockerfile with image name: scrubber and then runs a container, with name: scrubber_container, on port 3002 
	docker run -p 3002:3002 --name scrubber_container -ti --rm scrubber
 
.PHONY: update
update: ## Go gets all of the dependencies and downloads them
	go get .
	go mod download

.PHONY: debug
debug: ## Runs the api locally in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-nlp-search-scrubber
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-nlp-search-scrubber

.PHONY: test
test:	## Runs all tests with -race and -cover flags
	go test -race -cover ./...

.PHONY: convey
convey: ## Runs Convey tests
	goconvey ./...

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/ONSdigital/dp-nlp-search-scrubber/... -component


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