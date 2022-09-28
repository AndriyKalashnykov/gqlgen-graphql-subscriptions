.DEFAULT_GOAL := help

VERSION := $(shell cat server.go | grep "const Version ="| cut -d"\"" -f2)
GOFLAGS=-mod=mod

#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-14s\033[0m - %s\n", $$1, $$2}'

#clean: @ Cleanup
clean:
	@rm -rf ./.bin/
	@sudo rm -rf vendor/
	@mkdir ./.bin/

#generate: @ Generate GraphQL go source code
generate: clean
	@sudo rm -rf graph/model
	@sudo rm -rf graph/generated
	@export GOFLAGS=$(GOFLAGS); go run github.com/99designs/gqlgen generate

#test: @ Run tests
test: generate
	@export GOFLAGS=$(GOFLAGS); go test -v ./...

#build: @ Build Threeport GraphQL API
build:
	@export GOFLAGS=$(GOFLAGS); go build -o ./.bin/server server.go

#run: @ Run Threeport GraphQL API
run:
	@export GOFLAGS=$(GOFLAGS); go run server.go -graphqlPort="8080" -restApiAddr="http://172.28.1.11:1323" -restApiToken=""

#image: @ Build Docker image
image: generate
	docker build  -t gqlgen-graphql-subscriptions  .

#get: @ Download and install packages
get: clean
	@export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#deps: @ Download and install dependencies
deps:
	@export GOFLAGS=$(GOFLAGS); go install github.com/99designs/gqlgen@latest

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

#create-user: @ Execute createUser mutation
create-user:
	curl -s -H "Content-Type: application/json" -d '{ "query": "mutation {createUser(input:{Email:\"ben@enduser.com\",Password:\"secret\",FirstName:\"Ben\",LastName:\"Smith\",DateOfBirth:\"1985-01-30T00:00:00Z\",CountryOfResidence:\"United States\",Nationality:\"United States\"}){ID}}" }' http://localhost:8080/query | jq .

#get-user: @ Execute getUser query
get-user:
	curl -s -H "Content-Type: application/json" -d '{ "query": "{getUser(Email:\"ben@enduser.com\"){ID FirstName}}" }' http://localhost:8080/query | jq .

#get-user-by-id: @ Execute getUserById query
get-user-by-id:
	curl -s -H "Content-Type: application/json" -d '{ "query": "{getUserById(Id:8){ID Email FirstName LastName DateOfBirth CompanyID CountryOfResidence Nationality OtpSecret HasMfaConfigured PasswordResetToken}}" }' http://localhost:8080/query | jq .
