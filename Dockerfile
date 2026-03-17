FROM golang:1.26.1@sha256:c42e4d75186af6a44eb4159dcfac758ef1c05a7011a0052fe8a8df016d8e8fb9 AS builder

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
