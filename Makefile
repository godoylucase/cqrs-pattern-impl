ROOT_DIR := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
DOCKER_FILE := $(ROOT_DIR)/docker-compose.yml

.DEFAULT_TARGET: build

.PHONY: build
build:
	@echo "Building components for CQRS Architecture Example..."
	@docker-compose -f $(DOCKER_FILE) build --no-cache

.PHONY: run
run:
	@echo "Running CQRS Architecture Example ..."
	@docker-compose -f $(DOCKER_FILE) up -d --remove-orphans