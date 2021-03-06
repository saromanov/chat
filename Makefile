#!/usr/bin/make

SHELL = /bin/sh
LDFLAGS = "-s -w"

DOCKER_BIN = $(shell command -v docker 2> /dev/null)
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)" app
APP_NAME = $(notdir $(CURDIR))
CURRENT_TAG=$(shell git describe --exact-match --tags $(git log -n1 --pretty='%h') | cut -c 2-)
GO_RUN_ARGS ?=

.PHONY : help build fmt lint gotest test cover run shell redis-cli image clean
.DEFAULT_GOAL : help
.SILENT : test shell redis-cli

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build app binary file
	go build -o chat ./cmd/main.go

fmt: ## Run source code formatter tools
	$(DC_BIN) run $(DC_RUN_ARGS) sh -c 'GO111MODULE=off go get golang.org/x/tools/cmd/goimports && $$GOPATH/bin/goimports -d -w .'
	$(DC_BIN) run $(DC_RUN_ARGS) gofmt -s -w -d .

lint: ## Run app linters
	$(DOCKER_BIN) run --rm -t -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v

gotest: ## Run app tests
	$(DC_BIN) run $(DC_RUN_ARGS) go test -v -race ./...

test: lint gotest ## Run app tests and linters
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'All tests passed!';

cover: ## Run app tests with coverage report
	$(DC_BIN) run $(DC_RUN_ARGS) sh -c 'go test -race -covermode=atomic -coverprofile /tmp/cp.out ./... && go tool cover -html=/tmp/cp.out -o ./coverage.html'
	-sensible-browser ./coverage.html && sleep 2 && rm -f ./coverage.html

run: ## Run app without building binary file
	swag init -g pkg/server/server.go
	golangci-lint run ./...
	go run ./cmd/. $(GO_RUN_ARGS)

compose-run:
	./scripts/prod.sh

compose-run-dev:
	./scripts/dev.sh

swarm-run:
	$(DOCKER_BIN) stack deploy --compose-file docker-compose.yml $(APP_NAME)

swarm-rm:
	$(DOCKER_BIN) stack rm $(APP_NAME)

shell: ## Start shell into container with golang
	$(DC_BIN) run $(DC_RUN_ARGS) bash

image-push: ## Build docker image with app
	$(DOCKER_BIN) build -f ./images/chat/Dockerfile -t motorcode/$(APP_NAME):$(CURRENT_TAG) .
	$(DOCKER_BIN) push motorcode/$(APP_NAME):$(CURRENT_TAG)

clean: ## Make clean
	$(DC_BIN) down -v -t 1
	$(DOCKER_BIN) rmi $(APP_NAME) -f