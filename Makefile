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

.PHONY: all
all: delimiter-AUDIT audit delimiter-LINTERS lint delimiter-UNIT-TESTS test delimiter-COMPONENT_TESTS test-component delimiter-FINISH ## Runs multiple targets, audit, lint, test and test-component

.PHONY: audit
audit: ## Audits and finds vulnerable dependencies
	go list -json -m all | nancy sleuth

.PHONY: build 
build: Dockerfile ## Builds ./Dockerfile image name: scrubber
	docker build -t scrubber .

.PHONY: build-bin
build-bin: ## builds bin
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/scrubber

.PHONY: clean
clean: ## Removes /bin folder
	rm -fr ./build
	rm -fr ./vendor
	
.PHONY: debug
debug: ## Runs the api locally in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-nlp-search-scrubber
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-nlp-search-scrubber

.PHONY: delimiter-%
delimiter-%:
	@echo '===================${GREEN} $* ${RESET}==================='

.PHONY: fmt 
fmt: ## Formats the code using go fmt and go vet
	go fmt ./...
	go vet ./...

.PHONY: lint 
lint: ## Automated checking of your source code for programmatic and stylistic errors
	golangci-lint run --timeout=300s ./...

.PHONY: run
run: build ## First builds ./Dockerfile with image name: scrubber and then runs a container, with name: scrubber_container, on port 3002 
	docker run -p 3002:3002 --name scrubber_container -ti --rm scrubber

.PHONY: run-locally 
run-locally: ## Run the app locally
	go run .
 
.PHONY: test
test: ## Runs standard unit test tests
	go test -race -cover ./... 

.PHONY: test-all
test-all: convey test-component	test ## Runs all tests with -race and -cover flags
	go test -race -cover ./...

.PHONY: test-component
test-component: ## Runs component tests
	go test -cover -coverpkg=github.com/ONSdigital/dp-nlp-search-scrubber/... -component

.PHONY: test-convey
test-convey: ## Runs Convey tests
	goconvey ./...

.PHONY: update
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