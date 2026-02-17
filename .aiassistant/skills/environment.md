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
