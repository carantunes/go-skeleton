# Helper commands
SERVICE_NAME := go-skeleton
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_DIR_PATH := build/docker

API_PATH = "api"
CLI_PATH = "cli"

lint:
	golangci-lint run
.PHONY:lint

api:
	go run $(API_PATH)/*.go
.PHONY:api

cli:
	go run $(CLI_PATH)/*.go
.PHONY:api

DOCKER_COMPOSE := cd $(DOCKER_DIR_PATH) && docker-compose -f $(DOCKER_COMPOSE_FILE)

copy-env:
	cd $(DOCKER_DIR_PATH) && cp .env.dist .env
.PHONY: copy-env

build-env: copy-env
	$(DOCKER_COMPOSE) up -d --build --remove-orphans $(SERVICE_NAME)
.PHONY: build-env

shell-env:
	$(DOCKER_COMPOSE) run --service-ports $(SERVICE_NAME) /bin/sh
.PHONY: shell-env

env: build-env shell-env
.PHONY: env
