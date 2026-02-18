---
apply: always
---

# Refactoring Guidelines

## When to Refactor

### Always Refactor When:
- Adding new features that would duplicate existing code
- Code violates project conventions in `golang.md`
- Functions exceed ~50 lines or have deep nesting (>3 levels)
- Business logic is mixed with resolver/handler code
- Tests are difficult to write due to tight coupling
- Same pattern is repeated 3+ times (Rule of Three)

### Consider Refactoring When:
- Performance issues are identified
- Error handling is inconsistent
- Dependencies are tightly coupled
- Code is difficult to understand or maintain

## Refactoring Principles

### 1. Test First
- Ensure tests exist before refactoring
- Run `make test` before and after changes
- Add tests if coverage is missing

### 2. Small Steps
- Make incremental changes
- Commit after each logical refactoring step
- Keep the code working at each step

### 3. Don't Change Behavior
- Refactoring should not alter functionality
- Only improve structure, readability, or performance
- Bug fixes are separate from refactoring

## Common Refactoring Patterns

### Extract Function
When a function does too much:
```go
// Before
func (r *mutationResolver) CreateMessage(ctx context.Context, message string) (*model.Message, error) {
    m := model.Message{Message: message}
    r.RedisClient.XAdd(ctx, &redis.XAddArgs{
        Stream: "room",
        ID:     "*",
        MaxLen: 1,
        Values: map[string]interface{}{"message": m.Message},
    })
    return &m, nil
}

// After - extract Redis logic
func (r *mutationResolver) CreateMessage(ctx context.Context, message string) (*model.Message, error) {
    m := model.Message{Message: message}
    if err := r.publishMessage(ctx, &m); err != nil {
        return nil, err
    }
    return &m, nil
}

func (r *Resolver) publishMessage(ctx context.Context, m *model.Message) error {
    return r.RedisClient.XAdd(ctx, &redis.XAddArgs{
        Stream: "room",
        ID:     "*",
        MaxLen: 1,
        Values: map[string]interface{}{"message": m.Message},
    }).Err()
}
```

### Extract Service Layer
Move business logic from resolvers to service packages:
```go
// internal/service/message.go
type MessageService struct {
    redis *redis.Client
}

func (s *MessageService) Publish(ctx context.Context, msg string) (*model.Message, error) {
    // Business logic here
}

// graph/schema.resolvers.go
func (r *mutationResolver) CreateMessage(ctx context.Context, message string) (*model.Message, error) {
    return r.messageService.Publish(ctx, message)
}
```

### Introduce Constants
Replace magic values:
```go
// Before
r.RedisClient.XAdd(ctx, &redis.XAddArgs{
    Stream: "room",
    MaxLen: 1,
})

// After
const (
    RedisStreamRoom = "room"
    RedisStreamMaxLen = 1
)

r.RedisClient.XAdd(ctx, &redis.XAddArgs{
    Stream: RedisStreamRoom,
    MaxLen: RedisStreamMaxLen,
})
```

### Dependency Injection
Make dependencies explicit and testable:
```go
// Before
func NewResolver(client *redis.Client) *Resolver {
    return &Resolver{RedisClient: client}
}

// After - interface for testing
type RedisClient interface {
    XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StatusCmd
    XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

func NewResolver(client RedisClient) *Resolver {
    return &Resolver{RedisClient: client}
}
```

## Project-Specific Guidelines

### GraphQL Resolvers
- Keep resolvers thin - delegate to service layer
- Resolvers should only handle:
  - Input validation/transformation
  - Calling service methods
  - Error handling/formatting
- Never put database/Redis operations directly in resolvers

### Service Layer Patterns
**Create dedicated service structs:**
```go
type MessageService struct {
    redis datastore.RedisClient
}

func NewMessageService(redis datastore.RedisClient) *MessageService {
    return &MessageService{redis: redis}
}
```

**Separate concerns:**
- `PublishMessage()` - synchronous operations (mutations)
- `ReadMessages()` - query operations (fetch existing data)
- `StreamMessages()` - continuous operations (subscriptions)

**Each method should:**
- Validate inputs
- Handle errors with wrapped context
- Return domain models (not database types)

### Redis Operations
- Extract stream names and configuration to constants
- Use a Redis service wrapper for ALL operations
- Handle connection errors consistently
- **CRITICAL**: Different XREAD params for queries vs subscriptions (see golang.md)
- Always validate XREAD parameters before deploying

### Error Handling
- Use consistent `errors.Is(err, nil)` pattern
- Wrap errors with context: `fmt.Errorf("failed to publish message: %w", err)`
- Don't log errors in libraries/services - return them
- Handle `redis.Nil` separately from other errors
- Validate inputs before calling external services

### Concurrency
- Extract goroutine logic into named functions
- Use sync.Once for initialization
- Document goroutine lifecycle and cleanup
- Use select statements for context cancellation
- Always close channels when done producing

## What NOT to Refactor

### Don't Touch:
- Generated code in `graph/generated/` - regenerate with `make generate` instead
- Working code without tests (write tests first)
- Code you don't understand (study it first)
- External dependencies (upgrade/replace instead)

### Avoid Over-Engineering:
- Don't create abstractions for single use cases
- Don't prematurely optimize
- Keep it simple - prefer clarity over cleverness
- Don't add unnecessary layers

## Refactoring Checklist

Before committing refactored code:
- [ ] All tests pass (`make test`)
- [ ] Code follows `golang.md` conventions
- [ ] No behavior changes (unless intended)
- [ ] Error handling is consistent
- [ ] Functions are focused and small
- [ ] Magic values replaced with constants
- [ ] Dependencies are injected
- [ ] Code is more maintainable than before
- [ ] Generated code is up-to-date (`make generate`)
- [ ] Commit message explains the refactoring

## When in Doubt
- Prefer readability over cleverness
- Keep functions small and focused
- Follow existing patterns in the codebase
- Ask for review on significant changes
