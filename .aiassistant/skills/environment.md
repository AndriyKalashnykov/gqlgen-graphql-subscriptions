---
description: Development environment setup and tool locations
---

# Development Environment

## Go Compiler (via gvm)

### Active Version
- **Go Version**: Check `go.mod` file for required version (currently 1.26.0)
- **Managed by**: gvm (Go Version Manager)
- **Important**: Always reference `go.mod` for the active Go version requirement

### Paths
```bash
GOROOT=/home/andriy/.gvm/gos/go1.26.0
GOPATH=/home/andriy/.gvm/pkgsets/go1.26.0/global
GVM_ROOT=/home/andriy/.gvm
```

### Binaries
- **go**: `/home/andriy/.gvm/gos/go1.26.0/bin/go`
- **gofmt**: `/home/andriy/.gvm/gos/go1.26.0/bin/gofmt`

### Environment Activation
The Go environment is automatically loaded via `~/.zshrc`:
```bash
[[ -s "/home/andriy/.gvm/scripts/gvm" ]] && source "/home/andriy/.gvm/scripts/gvm"
gvm use go1.26.0 --default
```

### Available Go Versions
- go1.26.0 (active/default)

## Usage Notes

When running Go commands in a new shell session, the environment should already be active. If not, source gvm:
```bash
source ~/.gvm/scripts/gvm
gvm use go1.26.0
```

For this project, always check and use the Go version specified in `go.mod`.

## Build Flags
This project uses:
```bash
GOFLAGS=-mod=mod
```
Set in Makefile for all Go operations.

## Common Commands

### Backend (Go GraphQL Server)
```bash
make build          # Build the server binary to .bin/server
make run            # Build and run the server (port 8080)
make test           # Run tests with coverage
make generate       # Regenerate GraphQL code from schema
make kill-backend   # Kill all server processes and free port 8080
```

### Frontend (React App)
```bash
make build-frontend # Install deps and build frontend
make run-frontend   # Run frontend dev server (port 3000)
```

### Redis
```bash
make redis-up       # Start Redis in Docker (port 6379)
make redis-down     # Stop Redis and clean up
```

### Development Workflow
1. Start Redis: `make redis-up` (in separate terminal)
2. Start backend: `make run` (in separate terminal)
3. Start frontend: `make run-frontend` (in separate terminal)
4. Access app: http://localhost:3000
5. Access GraphQL playground: http://localhost:8080/playground

### Troubleshooting
- Backend port busy: `make kill-backend`
- Redis not responding: `make redis-down && make redis-up`
- Frontend not updating: Check if dev server reloaded (should auto-reload)
