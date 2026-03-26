.DEFAULT_GOAL := help

VERSION := $(shell cat server.go | grep "const Version ="| cut -d"\"" -f2)
GOFLAGS=-mod=mod

# === Tool Versions (pinned) ===
GQLGEN_VERSION := v0.17.86

#help: @ List available tasks
help:
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-20s\033[0m - %s\n", $$1, $$2}'

#clean: @ Cleanup
clean:
	@rm -rf ./.bin/
	@rm -rf vendor/
	@mkdir ./.bin/

#generate: @ Generate GraphQL go source code
generate:
	@rm -rf graph/model
	@rm -rf graph/generated
	@export GOFLAGS=$(GOFLAGS); go run github.com/99designs/gqlgen generate

#test: @ Run tests
test: generate
	@export GOFLAGS=$(GOFLAGS); go test -v ./...

#build: @ Build GraphQL API
build: generate
	@export GOFLAGS=$(GOFLAGS); go build -o ./.bin/server server.go

#run: @ Run GraphQL API
run: build kill-backend
	@export GOFLAGS=$(GOFLAGS); go run server.go

#image: @ Build Docker image
image: generate
	@docker buildx build --load -t gqlgen-graphql-subscriptions .

#build-frontend: @ Build JS client frontend
build-frontend:
	@cd ./frontend && yarn install && yarn build

#run-frontend: @ Run JS client frontend
run-frontend: build-frontend
	@cd ./frontend && yarn start

#image-frontend: @ Build JS client Docker image
image-frontend: build-frontend
	@cd ./frontend && docker build -t gqlgen-graphql-frontend .

#get: @ Download and install packages
get: clean
	@export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#deps: @ Install required tools (idempotent)
deps:
	@command -v go >/dev/null 2>&1 || { echo "Error: Go required. See https://go.dev/dl/"; exit 1; }
	@command -v gqlgen >/dev/null 2>&1 || { echo "Installing gqlgen $(GQLGEN_VERSION)..."; \
		export GOFLAGS=$(GOFLAGS); go install github.com/99designs/gqlgen@$(GQLGEN_VERSION); \
	}
	@command -v yarn >/dev/null 2>&1 || { echo "Installing yarn..."; npm install -g yarn; }

#lint: @ Run Go linter
lint:
	@export GOFLAGS=$(GOFLAGS); go vet ./...

#ci: @ Run full local CI pipeline
ci: deps generate lint test build
	@echo "Local CI pipeline passed."

#release: @ Create and push a new tag
release:
	@bash -c 'read -p "New tag (current: $(VERSION)): " newtag && \
		echo "$$newtag" | grep -qE "^v[0-9]+\.[0-9]+\.[0-9]+$$" || { echo "Error: Tag must match vN.N.N"; exit 1; } && \
		echo -n "Create and push $$newtag? [y/N] " && read ans && [ "$${ans:-N}" = y ] && \
		sed -i "s/const Version = \"$(VERSION)\"/const Version = \"$$newtag\"/" server.go && \
		git add -A && \
		git commit -s -m "Cut $$newtag release" && \
		git tag $$newtag && \
		git push origin $$newtag && \
		git push && \
		echo "Done."'

#update: @ Update dependencies to latest versions
update: clean
	@export GOFLAGS=$(GOFLAGS); go get -u; go mod tidy

#version: @ Print current version(tag)
version:
	@echo ${VERSION}

#redis-up: @ Start Redis
redis-up: redis-down
	@docker compose up

#redis-down: @ Stop Redis
redis-down:
	@docker compose down -v --remove-orphans

#kill-backend: @ Kill all backend server processes and free port 8080
kill-backend:
	@ps aux | grep -E "(\.bin/server|server\.go)" | grep -v grep | awk '{print $$2}' | xargs -r kill -9 && echo "Killed old server processes" || echo "No server processes found"
	@sleep 1
	@lsof -i:8080 2>/dev/null || echo "Port 8080 is now free"

.PHONY: help clean generate test build run image \
	build-frontend run-frontend image-frontend \
	get deps lint ci release update version \
	redis-up redis-down kill-backend
