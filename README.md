# gqlgen-graphql-subscriptions

GraphQL Subscriptions with Go and gqlgen

### Requirements

- [go 1.18](https://go.dev/doc/install)
- [gqlgen](github.com/99designs/gqlgen)
- [docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [curl](https://help.ubidots.com/en/articles/2165289-learn-how-to-install-run-curl-on-windows-macosx-linux)
- [jq](https://github.com/stedolan/jq/wiki/Installation)

### Help

```text
$ make
Usage: make COMMAND
Commands :
help           - List available tasks
clean          - Cleanup
generate       - Generate GraphQL go source code
test           - Run tests
build          - Build GraphQL API
run            - Run GraphQL API
image          - Build Docker image
get            - Download and install packages
deps           - Download and install dependencies
release        - Create and push a new tag. Modify `Version` field in `server.go` as it's used as an actual tag name
update         - Update dependencies to latest versions
version        - Print current version(tag)
```
