---
description: Common issues and their solutions for this GraphQL subscriptions project
---

# Troubleshooting Guide

## Backend Issues

### Issue: GraphQL Query Hangs Indefinitely

**Symptoms:**
- Queries never return
- curl requests timeout
- Frontend loading forever

**Root Cause:**
- Redis XREAD using `Block: 0` (infinite block) with query operations
- Using `"$"` (new messages only) instead of `"0"` (from beginning)

**Solution:**
```go
// ❌ WRONG - This will hang forever waiting for new messages
streams, err := redis.XRead(ctx, &redis.XReadArgs{
    Streams: []string{"room", "$"},
    Block:   0,  // BAD: infinite block
}).Result()

// ✅ CORRECT - Returns immediately with existing messages
streams, err := redis.XRead(ctx, &redis.XReadArgs{
    Streams: []string{"room", "0"},  // Read from beginning
    Block:   -1,  // Don't block
    Count:   100, // Limit results
}).Result()
```

**Prevention:**
- Always use `Block: -1` for query operations
- Reserve `Block: 0` ONLY for subscription streaming
- Test queries with `curl` before deploying

---

### Issue: Port 8080 Already in Use

**Symptoms:**
```
listen tcp :8080: bind: address already in use
```

**Solution:**
```bash
make kill-backend
```

Or manually:
```bash
lsof -ti:8080 | xargs kill -9
```

**Prevention:**
- Always use `make kill-backend` before restarting
- Don't run multiple instances

---

### Issue: Redis Connection Failed

**Symptoms:**
```
failed to connect to Redis: dial tcp [::1]:6379: connect: connection refused
```

**Solution:**
1. Check if Redis is running:
   ```bash
   docker ps | grep redis
   ```

2. Start Redis:
   ```bash
   make redis-up
   ```

3. Restart backend after Redis is up

**Prevention:**
- Always start Redis before backend
- Use `make redis-up` consistently

---

## Frontend Issues

### Issue: Messages Appear Then Disappear After Submit

**Symptoms:**
- Type message and click Submit
- Message shows briefly then vanishes
- Input doesn't clear

**Root Cause:**
- Apollo `useQuery` continuously refetching and replacing state
- Query result overwrites optimistically added messages

**Solution:**
```typescript
// Use ref to track initial load
const initialLoadDone = useRef(false)

useEffect(() => {
  // Only load from query once on mount
  if (queryResult.data?.messages && !initialLoadDone.current) {
    setMessages(queryResult.data.messages)
    initialLoadDone.current = true
  }
}, [queryResult.data?.messages])
```

**Prevention:**
- Don't let queries continuously overwrite state after initial load
- Use subscriptions for real-time updates
- Implement optimistic updates in mutations

---

### Issue: Submit Button Doesn't Work

**Symptoms:**
- Click Submit, nothing happens
- No messages appear
- No network requests in browser DevTools

**Diagnosis Steps:**
1. Check if backend is running:
   ```bash
   curl http://localhost:8080/
   # Should return: Welcome!
   ```

2. Check browser console (F12) for errors

3. Check Network tab for failed GraphQL requests

**Common Causes:**
- Backend server not running → Start with `make run`
- CORS issues → Check backend router CORS config
- GraphQL endpoint wrong → Should be http://localhost:8080/query

---

### Issue: WebSocket Connection Failed

**Symptoms:**
- Messages post but don't appear in real-time
- Console shows WebSocket errors
- Subscription not receiving updates

**Solution:**
1. Check WebSocket endpoint in `apolloClient.ts`:
   ```typescript
   uri: `ws://localhost:8080/subscriptions`  // Must match backend
   ```

2. Verify backend WebSocket handler is configured

3. Check if firewall blocking WebSocket

**Test:**
```bash
# Backend should have WebSocket route
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  http://localhost:8080/subscriptions
```

---

## Build Issues

### Issue: Build Fails - Interface Type Mismatch

**Symptoms:**
```
cannot use client (variable of type *redis.Client) as RedisClient value
```

**Solution:**
- Check interface method signatures match Redis client
- Common issue: `XAdd` returns `*redis.StringCmd` not `*redis.StatusCmd`
- Regenerate mocks after interface changes

---

### Issue: Tests Fail After Refactoring

**Symptoms:**
- Tests that passed now fail
- Mock interfaces out of sync

**Solution:**
1. Update all mock implementations to match new interfaces
2. Check method signatures:
   ```go
   // Interface changed?
   type RedisClient interface {
       XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd
   }

   // Update mock
   func (m *mockRedisClient) XAdd(ctx context.Context, args *redis.XAddArgs) *redis.StringCmd {
       // ...
   }
   ```

3. Run tests:
   ```bash
   make test
   ```

---

## Debugging Commands

### Check Running Processes
```bash
# Backend
ps aux | grep -E "(server|:8080)"

# Frontend
ps aux | grep -E "(react-scripts|node.*3000)"

# Redis
docker ps | grep redis
```

### Check Ports
```bash
lsof -i:8080  # Backend
lsof -i:3000  # Frontend
lsof -i:6379  # Redis
```

### View Logs
```bash
# Backend (if running in background)
tail -f /tmp/gql-server.log

# Redis
docker logs $(docker ps -q -f name=redis)
```

### Test GraphQL Manually
```bash
# Test query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { messages { id message } }"}'

# Test mutation
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { createMessage(message: \"test\") { id message } }"}'
```

---

## Prevention Checklist

Before coding:
- [ ] Read `golang.md` for patterns
- [ ] Check `refactoring.md` for guidelines
- [ ] Understand Redis XREAD parameters

Before testing:
- [ ] Redis is running (`make redis-up`)
- [ ] Backend is running (`make run`)
- [ ] Frontend is running (`make run-frontend`)
- [ ] No port conflicts (`make kill-backend` if needed)

Before committing:
- [ ] Tests pass (`make test`)
- [ ] Build succeeds (`make build`)
- [ ] Manual testing done (submit message, see it appear)
- [ ] No console errors in browser
