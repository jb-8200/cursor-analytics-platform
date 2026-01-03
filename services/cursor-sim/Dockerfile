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

# Runtime stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy binary from builder with proper ownership
COPY --from=builder --chown=nonroot:nonroot /out/cursor-sim /app/cursor-sim

# Copy default seed file with proper ownership (optional - can be overridden with volume mount)
COPY --chown=nonroot:nonroot testdata/valid_seed.json /app/seed.json

# Expose HTTP port
EXPOSE 8080

# Run as non-root user (defined in distroless base)
USER nonroot:nonroot

# Set entrypoint
ENTRYPOINT ["/app/cursor-sim"]

# Default arguments (can be overridden at runtime)
CMD ["-mode", "runtime", "-seed", "/app/seed.json", "-port", "8080", "-days", "90", "-velocity", "medium"]
