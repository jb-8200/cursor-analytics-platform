# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o cursor-sim ./cmd/simulator

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/cursor-sim .

# Expose port
EXPOSE 8080

# Run
CMD ["./cursor-sim"]
