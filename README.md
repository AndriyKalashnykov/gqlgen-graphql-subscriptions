[![CI](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml)
[![Hits](https://hits.sh/github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions.svg?view=today-total&style=plastic)](https://hits.sh/github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/AndriyKalashnykov/gqlgen-graphql-subscriptions)

# gqlgen-graphql-subscriptions

GraphQL Subscriptions example built with Go, [gqlgen](https://github.com/99designs/gqlgen), Echo v5, and Redis pub/sub. Includes a JavaScript frontend client that demonstrates real-time messaging between browser windows via WebSocket subscriptions.

## Quick Start

```bash
make deps              # install required tools
make redis-up          # start Redis (Terminal 1)
make run               # start GraphQL API (Terminal 2)
make run-frontend      # start JS client at http://localhost:3000 (Terminal 3)
```

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.26+ | Language runtime and compiler |
| [GNU Make](https://www.gnu.org/software/make/) | 3.81+ | Build orchestration |
| [Docker](https://www.docker.com/) | latest | Container builds and Redis |
| [Node.js / nvm](https://github.com/nvm-sh/nvm) | LTS | Frontend build toolchain |
| [Yarn](https://yarnpkg.com/) | 1.x | Frontend package manager |
| [curl](https://curl.se/) | latest | HTTP client (optional) |

Install all required dependencies:

```bash
make deps
```

## Available Make Targets

Run `make help` to see all available targets.

### Build & Run

| Target | Description |
|--------|-------------|
| `make build` | Build GraphQL API |
| `make run` | Run GraphQL API |
| `make generate` | Generate GraphQL go source code |
| `make build-frontend` | Build JS client frontend |
| `make run-frontend` | Run JS client frontend |

### Docker

| Target | Description |
|--------|-------------|
| `make image-build` | Build Docker image |
| `make image-frontend` | Build JS client Docker image |

### Infrastructure

| Target | Description |
|--------|-------------|
| `make redis-up` | Start Redis |
| `make redis-down` | Stop Redis |

### Code Quality

| Target | Description |
|--------|-------------|
| `make lint` | Run Go linter |
| `make test` | Run tests |

### CI

| Target | Description |
|--------|-------------|
| `make ci` | Run full local CI pipeline |
| `make ci-run` | Run GitHub Actions workflow locally via [act](https://github.com/nektos/act) |

### Utilities

| Target | Description |
|--------|-------------|
| `make deps` | Install required tools (idempotent) |
| `make deps-act` | Install act for local CI (idempotent) |
| `make get` | Download and install packages |
| `make update` | Update dependencies to latest versions |
| `make clean` | Cleanup |
| `make version` | Print current version (tag) |
| `make release` | Create and push a new tag |
| `make kill-backend` | Kill all backend server processes and free port 8080 |
| `make renovate-validate` | Validate Renovate configuration |

## Run

### Terminal 1

Start Redis:

```shell
make redis-up
```

### Terminal 2

Run GraphQL API:

```shell
make run
```

### Terminal 3

Run JS client frontend. Command below should open a browser at [http://localhost:3000](http://localhost:3000).
Open another window at [http://localhost:3000](http://localhost:3000) post a message and see it appear in the other window.

```shell
make run-frontend
```

## CI/CD

GitHub Actions runs on every push to `main`, tags `v*`, and pull requests.

| Job | Triggers | Steps |
|-----|----------|-------|
| **ci** | push, PR, tags | Generate, Lint, Test, Build |
| **docker** | tags only | Build Docker image, Build JS client Docker image |

[Renovate](https://docs.renovatebot.com/) keeps dependencies up to date with platform automerge enabled.

## References

- https://redis.io/topics/streams-intro
- https://github.com/go-redis/redis
- https://pkg.go.dev/github.com/go-redis/redis/v9
- https://towardsdev.com/scalable-event-streaming-with-redis-streams-and-go-dee5fbe8982c
- https://github.com/gmrdn/redis-streams-go
