[![CI](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/actions/workflows/ci.yml)
[![Hits](https://hits.sh/github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions.svg?view=today-total&style=plastic)](https://hits.sh/github.com/AndriyKalashnykov/gqlgen-graphql-subscriptions/)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://app.renovatebot.com/dashboard#github/AndriyKalashnykov/gqlgen-graphql-subscriptions)

# gqlgen-graphql-subscriptions

GraphQL Subscriptions example built with Go, [gqlgen](https://github.com/99designs/gqlgen), Echo v5, and Redis pub/sub. Includes a TypeScript/React frontend client that demonstrates real-time messaging between browser windows via WebSocket subscriptions.

| Component | Technology |
|-----------|------------|
| Language | Go 1.26 |
| GraphQL | gqlgen v0.17.89 |
| Router | Echo v5 |
| Pub/Sub | Redis (go-redis/v9) |
| Frontend | TypeScript / React 19 / Vite 8 |
| WebSocket | graphql-ws |
| Container | Docker / Docker Compose |

## Quick Start

```bash
make deps              # install required tools
make redis-up          # start Redis (Terminal 1)
make run               # start GraphQL API (Terminal 2)
make frontend-run      # start frontend at http://localhost:3000 (Terminal 3)
```

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| [Git](https://git-scm.com/) | 2.0+ | Version control |
| [Go](https://go.dev/dl/) | 1.26+ | Language runtime and compiler |
| [GNU Make](https://www.gnu.org/software/make/) | 3.81+ | Build orchestration |
| [Docker](https://www.docker.com/) | latest | Container builds and Redis |
| [Node.js / nvm](https://github.com/nvm-sh/nvm) | 24 (see `.nvmrc`) | Frontend build toolchain |
| [pnpm](https://pnpm.io/) | 10+ | Frontend package manager |

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
| `make frontend-build` | Build frontend client |
| `make frontend-run` | Run frontend client |

### Docker

| Target | Description |
|--------|-------------|
| `make image-build` | Build Docker image |
| `make frontend-image` | Build frontend Docker image |

### Infrastructure

| Target | Description |
|--------|-------------|
| `make redis-up` | Start Redis |
| `make redis-down` | Stop Redis |

### Code Quality

| Target | Description |
|--------|-------------|
| `make static-check` | Generate code, run all linters, and scan for vulnerabilities |
| `make lint` | Run golangci-lint (includes gocritic, gosec) and hadolint |
| `make trivy-fs` | Scan filesystem for vulnerabilities, secrets, and misconfigurations |
| `make test` | Run tests |
| `make coverage-check` | Run tests with coverage and verify threshold |

### CI

| Target | Description |
|--------|-------------|
| `make ci` | Run full local CI pipeline |
| `make ci-run` | Run GitHub Actions workflow locally via [act](https://github.com/nektos/act) |

### Utilities

| Target | Description |
|--------|-------------|
| `make help` | List available tasks |
| `make deps` | Install required tools (idempotent) |
| `make deps-act` | Install act for local CI (idempotent) |
| `make deps-hadolint` | Install hadolint for Dockerfile linting |
| `make deps-trivy` | Install Trivy for security scanning |
| `make get` | Download and install packages |
| `make update` | Update dependencies to latest versions |
| `make clean` | Cleanup |
| `make version` | Print current version (tag) |
| `make release` | Create and push a new tag |
| `make kill-backend` | Kill all backend server processes and free port 8080 |
| `make renovate-bootstrap` | Install nvm and Node.js for Renovate |
| `make renovate-validate` | Validate Renovate configuration |

## CI/CD

GitHub Actions runs on every push to `main`, tags `v*`, and pull requests.

| Job | Needs | Steps |
|-----|-------|-------|
| **static-check** | — | Generate, Lint, Trivy scan |
| **build** | static-check | Build server binary |
| **test** | static-check | Test with coverage (parallel with build) |
| **docker** | build (tags only) | Build Docker image, Build frontend Docker image |

Job graph: `static-check` → `build` + `test` (parallel) → `docker` (tags only).

[Renovate](https://docs.renovatebot.com/) keeps dependencies up to date with platform automerge enabled.

## References

- https://redis.io/docs/latest/develop/data-types/streams/
- https://github.com/redis/go-redis
- https://pkg.go.dev/github.com/redis/go-redis/v9
- https://towardsdev.com/scalable-event-streaming-with-redis-streams-and-go-dee5fbe8982c
- https://github.com/gmrdn/redis-streams-go
