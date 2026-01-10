# syntax=docker/dockerfile:1.7

# Builder stage
FROM golang:1.22-alpine AS builder

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Set build environment for static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOTOOLCHAIN=auto

WORKDIR /src

# Cache dependencies (layer caching optimization)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build static binary with size optimization
RUN go build \
    -ldflags="-s -w -X main.version=2.0.0" \
    -o /out/cursor-sim \
    ./cmd/simulator

# Runtime stage - use alpine for wget health check support
FROM alpine:3.19

# Install wget for health checks and CA certificates
RUN apk add --no-cache wget ca-certificates && \
    adduser -D -u 1000 appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /out/cursor-sim /app/cursor-sim

# Copy default seed file (can be overridden with volume mount)
COPY testdata/valid_seed.json /app/seed.json

# Set ownership
RUN chown -R appuser:appuser /app

# Expose HTTP port
EXPOSE 8080

# Run as non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["/app/cursor-sim"]

# Default arguments (can be overridden at runtime or via environment variables)
# Use environment variables for days/velocity instead of hardcoding in CMD
CMD ["-mode", "runtime", "-seed", "/app/seed.json", "-port", "8080"]
