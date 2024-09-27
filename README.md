[![CI](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml)
[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2FAndriyKalashnykov%2Fgqlgen-graphql-subscriptions&count_bg=%2340C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/AndriyKalashnykov/gqlgen-graphql-subscriptions)
# gqlgen-graphql-subscriptions

GraphQL Subscriptions with Go and gqlgen

### Requirements

* [gvm](https://github.com/moovweb/gvm) Go 1.23
    ```bash
    gvm install go1.23 --prefer-binary --with-build-tools --with-protobuf
    gvm use go1.23 --default
    ```
- [gqlgen](https://github.com/99designs/gqlgen)
- [docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/install/)
- [nvm](https://github.com/nvm-sh/nvm#install--update-script)
  ```shell
  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
  nvm install v22.1.0
  nvm use nvm v22.1.0
  cd ./frontend && pnpm install -g yarn
  ```
- [curl](https://help.ubidots.com/en/articles/2165289-learn-how-to-install-run-curl-on-windows-macosx-linux)
- [jq](https://github.com/stedolan/jq/wiki/Installation)

### Help

```text
$ make
Usage: make COMMAND
Commands :
help             - List available tasks
clean            - Cleanup
generate         - Generate GraphQL go source code
test             - Run tests
build            - Build GraphQL API
run              - Run GraphQL API
image            - Build Docker image
build-frontend   - Build JS client frontend
run-frontend     - Run JS client frontend
image-frontend   - Build JS client Docker image
get              - Download and install packages
deps             - Download and install dependencies
release          - Create and push a new tag. Modify `Version` field in `server.go` as it's used as an actual tag name
update           - Update dependencies to latest versions
version          - Print current version(tag)
redis-up         - Start Redis
redis-down       - Stop Redis
```
### Run

#### Terminal 1

Start Redis
```shell
make redis-up
```

#### Terminal 2

Run GraphQL API
```shell
make run
```

#### Terminal 3

Run JS client frontend. Command below should open a browser at [http://localhost:3000](http://localhost:3000).
Open another window at [http://localhost:3000](http://localhost:3000) post a message and see it appear in the other window.

```shell
make run-frontend
```

### References

* https://redis.io/topics/streams-intro
* https://github.com/go-redis/redis
* https://pkg.go.dev/github.com/go-redis/redis/v9
* https://towardsdev.com/scalable-event-streaming-with-redis-streams-and-go-dee5fbe8982c
* https://github.com/gmrdn/redis-streams-go
