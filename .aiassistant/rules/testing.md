---
apply: always
---

# Testing Guidelines

## Testing Philosophy

- Tests are first-class code - maintain them with the same care as production code
- Write tests before or alongside implementation (TDD preferred)
- Tests should be fast, isolated, and deterministic
- Every bug fix should include a test that would have caught it

## Running Tests

### Standard Workflow
```bash
make test          # Run all tests (includes make generate)
go test ./...      # Run tests directly
go test -v ./...   # Verbose output
go test -run TestName  # Run specific test
```

### Test Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Organization

### File Structure
- Test files: `*_test.go` in same package
- Package naming: `package mypackage` (white-box) or `package mypackage_test` (black-box)
- Prefer white-box tests for unit tests, black-box for integration tests

### Test Function Naming
```go
func TestFunctionName(t *testing.T)           // Basic test
func TestFunctionName_Scenario(t *testing.T)  // Specific scenario
func TestFunctionName_ErrorCase(t *testing.T) // Error cases
```

Examples:
- `TestNewRedisClient`
- `TestNewRedisClient_ConnectionError`
- `TestCreateMessage_PublishSuccess`
- `TestMessageCreated_SubscriptionCleanup`

## Test Patterns

### Table-Driven Tests
Preferred pattern for testing multiple scenarios:

```go
func TestCreateMessage(t *testing.T) {
    tests := []struct {
        name    string
        message string
        want    *model.Message
        wantErr bool
    }{
        {
            name:    "valid message",
            message: "hello",
            want:    &model.Message{Message: "hello"},
            wantErr: false,
        },
        {
            name:    "empty message",
            message: "",
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := createMessage(tt.message)
            if (err != nil) != tt.wantErr {
                t.Errorf("createMessage() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("createMessage() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Subtests
Use `t.Run()` for organizing related tests:

```go
func TestResolver(t *testing.T) {
    t.Run("CreateMessage", func(t *testing.T) {
        // test create message
    })

    t.Run("Messages", func(t *testing.T) {
        // test query messages
    })
}
```

## Mocking and Test Doubles

### Redis Client Mocking
Create interface for Redis operations:

```go
// Internal package
type RedisClient interface {
    XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StatusCmd
    XRead(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

// Test file
type mockRedisClient struct {
    xAddFunc  func(ctx context.Context, args *redis.XAddArgs) *redis.StatusCmd
    xReadFunc func(ctx context.Context, args *redis.XReadArgs) *redis.XStreamSliceCmd
}

func (m *mockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StatusCmd {
    if m.xAddFunc != nil {
        return m.xAddFunc(ctx, args)
    }
    return redis.NewStatusCmd(ctx)
}
```

### Test Fixtures
For complex test data:

```go
// test_helpers.go
func newTestResolver(t *testing.T) *Resolver {
    t.Helper()
    return &Resolver{
        RedisClient: &mockRedisClient{},
        messageChannels: make(map[string]chan *model.Message),
    }
}

func newTestMessage() *model.Message {
    return &model.Message{
        ID:      "test-id",
        Message: "test message",
    }
}
```

## Testing GraphQL Components

### Testing Resolvers
```go
func TestMutationResolver_CreateMessage(t *testing.T) {
    ctx := context.Background()
    mock := &mockRedisClient{
        xAddFunc: func(ctx context.Context, args *redis.XAddArgs) *redis.StatusCmd {
            // Verify args
            if args.Stream != "room" {
                t.Errorf("expected stream 'room', got %s", args.Stream)
            }
            return redis.NewStatusCmd(ctx)
        },
    }

    resolver := &Resolver{RedisClient: mock}
    mr := &mutationResolver{resolver}

    msg, err := mr.CreateMessage(ctx, "hello")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if msg.Message != "hello" {
        t.Errorf("expected message 'hello', got %s", msg.Message)
    }
}
```

### Testing Subscriptions
```go
func TestSubscriptionResolver_MessageCreated(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    resolver := newTestResolver(t)
    sr := &subscriptionResolver{resolver}

    ch, err := sr.MessageCreated(ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Test channel cleanup on context cancellation
    cancel()
    time.Sleep(10 * time.Millisecond)

    if len(resolver.messageChannels) != 0 {
        t.Error("expected message channels to be cleaned up")
    }
}
```

## Testing Best Practices

### Context Handling
```go
// Use context with timeout for tests
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// For subscription tests
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

### Error Handling in Tests
```go
// Use consistent error checking pattern
if !errors.Is(err, nil) {
    t.Fatalf("unexpected error: %v", err)
}

// For expected errors
if err == nil {
    t.Fatal("expected error, got nil")
}

// Check error messages when needed
if err == nil || !strings.Contains(err.Error(), "connection refused") {
    t.Errorf("expected connection error, got: %v", err)
}
```

### Test Cleanup
```go
func TestSomething(t *testing.T) {
    // Setup
    cleanup := setupTest(t)
    defer cleanup()

    // Or use t.Cleanup
    t.Cleanup(func() {
        // cleanup code
    })
}
```

### Parallel Tests
```go
func TestParallel(t *testing.T) {
    t.Parallel() // Mark test as safe to run in parallel

    tests := []struct{
        name string
    }{
        {name: "test1"},
        {name: "test2"},
    }

    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // test code
        })
    }
}
```

## Integration Tests

### Redis Integration Tests
```go
// +build integration

func TestRedisIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    ctx := context.Background()
    client, err := datastore.NewRedisClient(ctx, "localhost:6379")
    if err != nil {
        t.Skipf("Redis not available: %v", err)
    }
    defer client.Close()

    // Test with real Redis
}
```

Run integration tests:
```bash
go test -tags=integration ./...
go test -short ./...  # Skip integration tests
```

## Test Documentation

### Test Comments
```go
// TestCreateMessage_EmptyInput verifies that CreateMessage returns
// an error when given an empty message string, preventing invalid
// data from being published to Redis.
func TestCreateMessage_EmptyInput(t *testing.T) {
    // ...
}
```

### Example Tests
Use Example tests for documentation:
```go
func ExampleResolver_CreateMessage() {
    ctx := context.Background()
    resolver := newResolver()

    msg, err := resolver.Mutation().CreateMessage(ctx, "hello")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(msg.Message)
    // Output: hello
}
```

## What to Test

### Always Test:
- Public API functions and methods
- Error conditions and edge cases
- Concurrent operations (goroutines, channels)
- Context cancellation and timeouts
- Boundary conditions (empty, nil, max values)

### Consider Testing:
- Private functions with complex logic
- Integration points (Redis, external services)
- Performance-critical paths (benchmarks)

### Don't Test:
- Generated code (`graph/generated/`)
- Trivial getters/setters
- Third-party libraries
- Code that's only glue (no logic)

## Test Quality Checklist

- [ ] Test names clearly describe what is being tested
- [ ] Tests are independent (can run in any order)
- [ ] Tests clean up resources (defer, t.Cleanup)
- [ ] Error messages are descriptive
- [ ] Mocks are used appropriately (not over-mocked)
- [ ] Table-driven tests used for multiple scenarios
- [ ] Context handling is correct
- [ ] Tests are fast (< 1 second for unit tests)
- [ ] Tests follow project error handling pattern (`errors.Is(err, nil)`)

## Benchmarking

```go
func BenchmarkCreateMessage(b *testing.B) {
    ctx := context.Background()
    resolver := newTestResolver(b)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := resolver.Mutation().CreateMessage(ctx, "test")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run benchmarks:
```bash
go test -bench=. ./...
go test -bench=BenchmarkCreateMessage -benchmem
```

## Common Pitfalls

### Avoid:
- Tests that depend on external services without fallback
- Tests with sleep/timing dependencies (use channels/sync)
- Sharing state between tests
- Testing implementation details instead of behavior
- Ignoring test failures or flaky tests
- Over-mocking (test becomes meaningless)

### Remember:
- Use `t.Helper()` in helper functions for better error reporting
- Capture range variables in parallel subtests
- Always defer cleanup functions
- Use `testing.Short()` for slow tests
- Mock at boundaries, not everywhere
