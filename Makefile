.DEFAULT_GOAL := help

VERSION := $(shell cat server.go | grep "const Version ="| cut -d"\"" -f2)
GOFLAGS=-mod=mod

#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-16s\033[0m - %s\n", $$1, $$2}'

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
run: build
	@export GOFLAGS=$(GOFLAGS); go run server.go

#image: @ Build Docker image
image: generate
	@docker build  -t gqlgen-graphql-subscriptions  .

#build-frontend: @ Build JS client frontend
build-frontend:
	@rm -Rf ./frontend/node_modules && rm -f ./frontend/yarn.lock && cd ./frontend && yarn install && yarn upgrade --latest && yarn build

#run-frontend: @ Run JS client frontend
run-frontend:
	@cd ./frontend && yarn start

#image-frontend: @ Build JS client Docker image
image-frontend: build-frontend
	@cd ./frontend  && docker build -t gqlgen-graphql-frontend  .

#get: @ Download and install packages
get: clean
	@export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#deps: @ Download and install dependencies
deps:
	@export GOFLAGS=$(GOFLAGS); go install github.com/99designs/gqlgen@latest
	@pnpm install -g yarn

#release: @ Create and push a new tag. Modify `Version` field in `server.go` as it's used as an actual tag name
release:
	@echo -n "Are you sure to create and push ${VERSION} tag? [y/N] " && read ans && [ $${ans:-N} = y ]
	@git commit -s -m "Cut ${VERSION} release"
	@git tag ${VERSION}
	@git push origin ${VERSION}
	@git push
	@echo "Done."

#update: @ Update dependencies to latest versions
update: clean
	@export GOFLAGS=$(GOFLAGS); go get -u; go mod tidy

#version: @ Print current version(tag)
version:
	@echo ${VERSION}

#redis-up: @ Start Redis
redis-up: redis-down
	docker-compose up

#redis-down: @ Stop Redis
redis-down:
	docker-compose down -v --remove-orphans
