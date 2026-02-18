---
apply: always
---

# Go Code Style Guide for gqlgen-graphql-subscriptions

## Project Overview
This is a GraphQL subscriptions service built with:
- **gqlgen** for GraphQL schema-first development
- **Echo v5** web framework
- **Redis Streams** for real-time messaging
- **WebSocket** for GraphQL subscriptions

## General Go Conventions

### Error Handling
- Use `errors.Is(err, nil)` pattern consistently (as seen throughout the codebase)
- Always handle errors explicitly; never ignore them
- Return errors up the call stack; log only at the top level (main/handlers)
- Use `log.Fatalln()` for unrecoverable startup errors

### Code Organization
- Follow standard Go project layout:
  - `internal/` for private application code
  - `graph/` for GraphQL schema, resolvers, and generated code
  - Package names match directory names
- Keep packages focused and cohesive
- Use interfaces for dependencies (e.g., `*redis.Client`)

### Naming Conventions
- Use descriptive names: `NewRedisClient`, `CreateMessage`, `MessageCreated`
- Constructors: prefix with `New` (e.g., `NewResolver`, `NewRouter`)
- Interfaces: typically noun or noun-phrase
- Package names: short, lowercase, no underscores

## gqlgen Specific Guidelines

### Schema Management
- GraphQL schemas go in `graph/*.graphqls`
- Generated code lives in `graph/generated/`
- Models in `graph/model/`
- Keep resolvers in `graph/schema.resolvers.go`

### Resolvers
- Use embedded resolver pattern (mutationResolver, queryResolver, subscriptionResolver)
- Keep resolver methods focused on orchestration
- Delegate business logic to service/datastore layers
- Return proper GraphQL types from `graph/model`

### Code Generation
- Always run `make generate` after schema changes
- Don't manually edit generated files
- Use `gqlgen.yml` for configuration
- Run generation before builds and tests

### Subscriptions
- Use Go channels for subscription streams (`<-chan *model.Message`)
- Implement proper cleanup with context cancellation (`<-ctx.Done()`)
- Use mutexes for concurrent access to shared subscription state
- Generate unique tokens for subscription tracking

## Redis & Datastore

### Redis Client
- Initialize once at startup, pass as dependency
- Use `defer client.Close()` in main
- Use context for all Redis operations
- Handle connection errors at startup (fatal) and runtime (return error)

### Redis Streams
- Use `XAdd` for publishing messages
- Use `XRead` for consuming messages
- Configure `MaxLen` to limit stream size
- Use proper error handling with `errors.Is(err, nil)`

### Redis XREAD Parameters - CRITICAL
**Queries vs Subscriptions:**
- **Queries** (fetch existing messages):
  - Stream ID: `"0"` (read from beginning) NOT `"$"` (new messages only)
  - Block: `-1` (don't block) or omit
  - Count: Set limit (e.g., `100`) to prevent loading too many messages
- **Subscriptions** (real-time updates):
  - Stream ID: `"$"` (only new messages)
  - Block: `0` (block indefinitely waiting for new messages)
  - Count: `1` or small number for real-time delivery

**Common Pitfall:** Using `Block: 0` with queries causes infinite hanging. Always use non-blocking reads for queries.

## Web Framework (Echo v5)

### Router Setup
- Centralize route configuration in `internal/router`
- Use Echo middleware appropriately
- Bind GraphQL handler to `/graphql` endpoint
- Enable WebSocket support for subscriptions

### Server Lifecycle
- Initialize all dependencies before starting server
- Use proper error handling for `e.Start()`
- Graceful shutdown should be implemented for production

## Testing

### Test Conventions
- Run tests with `make test`
- Tests should be isolated and repeatable
- Use table-driven tests where appropriate
- Mock external dependencies (Redis, etc.)

## Dependencies Management

### Module Management
- Use `go.mod` with Go 1.26.0+
- Run `make update` to update dependencies
- Use `GOFLAGS=-mod=mod` for module-aware builds
- Keep dependencies minimal and well-maintained

### Key Dependencies
- `github.com/99designs/gqlgen` - GraphQL server
- `github.com/labstack/echo/v5` - Web framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/gorilla/websocket` - WebSocket support

## Build & Development

### Makefile Targets
- `make generate` - Regenerate GraphQL code
- `make build` - Build the server binary
- `make test` - Run tests
- `make run` - Build and run server
- `make redis-up/redis-down` - Manage Redis container

### Docker
- Use multi-stage builds for production images
- Keep images minimal and secure
- Include necessary runtime dependencies only

## Code Quality

### Formatting
- Use `gofmt` (automatically applied)
- Maintain consistent indentation (tabs)
- Keep lines reasonably short (< 120 chars)

### Best Practices
- Avoid global state
- Use dependency injection
- Keep functions small and focused
- Write self-documenting code
- Add comments for complex logic, not obvious code
- Use constants for magic values (e.g., `const Version`)

### Concurrency
- Use mutexes for shared state protection
- Close channels when done producing
- Use context for cancellation
- Avoid goroutine leaks with proper cleanup

## Version Control
- Semantic versioning in `server.go` (`const Version`)
- Use `make release` to tag releases
- Follow conventional commit messages
- Keep commits atomic and focused
