SHELL = /bin/bash
.SHELLFLAGS = -o pipefail -e -c

GOPRIVATE := github.com/kohofinancial
SERVICE_NAME := trading

tab := $(shell printf '\t')

WORKSPACE_PATH = $(PWD)

.PHONY: dep
dep: ## Install code / documentation dependencies
	go mod tidy

.PHONY: lint
lint: ## Lint go files
	@hasdiff=$$(gofmt -l . 2>&1 | grep -v '^vendor/' || true); \
	if [ -n "$$hasdiff" ]; then \
		echo "$$hasdiff"; \
		exit 1; \
	fi; \
	echo "no lint errors"

.PHONY: platform
platform: ## Start the dependant platform services
		docker compose --profile platform up --detach

.PHONY: start
start: ## Start / re-start service
	docker compose --profile service up --detach


.PHONY: stop
stop: ## Stop service
	docker compose stop

.PHONY: test
test: ## Run tests inside docker environment
	go test -tags=integration test/it/it_test.go
	go test ./...

.PHONY: clean
clean: ## Clean up containers and images. Use "cleanCache=true" to cleanup go build cache.
	docker compose down --remove-orphans --volumes --rmi local
	docker network prune -f

.PHONY: build
build: ## Builds the image
	docker build -t homedepot/backend:latest .

.PHONY: start-local
start-local: build platform start ## Start / re-start  local service


.PHONY: start-no-build
start-no-build: platform start ## Start / re-start  local service