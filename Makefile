#!/bin/bash

.PHONY: generate_api
generate_api:
	docker run -v ${PWD}/api:/defs namely/protoc-all -f api.proto -l go

.PHONY: build
build:
	docker build -t server -f deployments/server/Dockerfile . && \
	docker build -t client -f deployments/client/Dockerfile .

.PHONY: start
start:
	docker-compose up -d

.PHONY: start-server
start-server:
	docker compose up server

.PHONY: start-server
stop-server:
	docker compose stop server

.PHONY: start-client
start-client:
	docker compose up client

.PHONY: lint
lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint golangci-lint run

.PHONY: test
test:
	go test ./... -count=1
	