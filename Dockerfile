FROM golang:1.26.0@sha256:c83e68f3ebb6943a2904fa66348867d108119890a2c6a2e6f07b38d0eb6c25c5 AS builder

WORKDIR /source
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Cache deps before building and copying source so that we don't need to re-download
# as much and so that source changes don't invalidate our downloaded layer
RUN GOCACHE=OFF
RUN go mod download

# Copy source code
COPY gqlgen.yml gqlgen.yml
COPY server.go server.go
COPY graph/ graph/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o server server.go

FROM scratch
COPY --from=builder /source/server /server
CMD ["/server"]
