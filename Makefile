SHELL=/bin/bash

.PHONY: clean
clean:
	@golangci-lint cache clean

.PHONY: setup
setup:
	@pre-commit install

.PHONY: pre-commit
pre-commit:
	@pre-commit run --all-files

.PHONY: lint
lint:
	@scripts/lint.sh

.PHONY: test
test:
	@go test -coverprofile="coverage.txt" -covermode=atomic ./... -count=1

.PHONY: gen-src
gen-src:
	mockery

.PHONY: update-mod
update-mod:
	@go get -u ./...

dev-docker-up:
	@docker compose -f docker/development/docker-compose.yml up -d

dev-docker-down:
	@docker compose -f docker/development/docker-compose.yml down

test-docker-up:
	@docker compose -f docker/test/docker-compose.yml up -d

test-docker-down:
	@docker compose -f docker/test/docker-compose.yml down
