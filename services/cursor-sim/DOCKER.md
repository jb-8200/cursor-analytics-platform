# cursor-sim - Docker Deployment Guide

Comprehensive guide for running cursor-sim in Docker containers, from local development to production deployment.

## Table of Contents

- [Quick Start](#quick-start)
- [Automated Script (Recommended)](#automated-script-recommended)
- [Manual Docker Commands](#manual-docker-commands)
- [Environment Variables](#environment-variables)
- [Docker Image Details](#docker-image-details)
- [Custom Seed Files](#using-custom-seed-files)
- [Container Management](#container-management)
- [Troubleshooting](#troubleshooting)
- [Development Workflow](#development-workflow)
- [Production Deployment](#production-deployment)
- [Configuration Presets](#configuration-presets)

---

## Quick Start

### Automated Script (Recommended)

The easiest way to run cursor-sim locally with Docker:

```bash
# From project root
./tools/docker-local.sh
```

This automatically:
- Builds the Docker image if needed
- Starts the container with health check verification
- Shows service URL and example commands
- Cleans up gracefully on exit

### Custom Configuration

```bash
# Run with custom simulation parameters
DAYS=30 VELOCITY=high ./tools/docker-local.sh

# Run in background (detached mode)
DETACH=true ./tools/docker-local.sh

# Use different port (if 8080 is busy)
PORT=8081 ./tools/docker-local.sh

# Force rebuild the image
REBUILD=true ./tools/docker-local.sh

# Combine multiple options
PORT=8081 DAYS=30 VELOCITY=high DETACH=true ./tools/docker-local.sh
```

---

## Manual Docker Commands

### Build Image

```bash
# From project root
docker build -t cursor-sim:latest services/cursor-sim

# Build with specific tag
docker build -t cursor-sim:v2.0.0 services/cursor-sim

# Build with no cache
docker build --no-cache -t cursor-sim:latest services/cursor-sim
```

### Run Container

**Foreground mode (with logs):**
```bash
docker run --rm -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  -e CURSOR_SIM_VELOCITY=medium \
  cursor-sim:latest
```

**Background mode (detached):**
```bash
docker run -d --name cursor-sim-local -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  -e CURSOR_SIM_VELOCITY=medium \
  cursor-sim:latest
```

**With custom port:**
```bash
docker run --rm -p 8081:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  cursor-sim:latest
```

---

## Environment Variables

| Variable | Default | Options | Description |
|----------|---------|---------|-------------|
| `CURSOR_SIM_MODE` | `runtime` | `runtime`, `replay` | Simulation mode (runtime generates data, replay uses corpus) |
| `CURSOR_SIM_SEED` | `/app/seed.json` | file path | Path to seed data file with team/developer configuration |
| `CURSOR_SIM_DAYS` | `90` | number | Number of days of historical data to simulate |
| `CURSOR_SIM_VELOCITY` | `medium` | `low`, `medium`, `high` | Commit frequency (low=~5/day, medium=~15/day, high=~30/day) |
| `CURSOR_SIM_PORT` | `8080` | port number | HTTP server port (inside container) |

### Setting Environment Variables

**Via `-e` flag:**
```bash
docker run -e CURSOR_SIM_DAYS=180 -e CURSOR_SIM_VELOCITY=high cursor-sim:latest
```

**Via `.env` file:**
```bash
# Create .env file
cat > cursor-sim.env <<EOF
CURSOR_SIM_MODE=runtime
CURSOR_SIM_SEED=/app/seed.json
CURSOR_SIM_DAYS=180
CURSOR_SIM_VELOCITY=high
EOF

# Use with --env-file
docker run --env-file cursor-sim.env -p 8080:8080 cursor-sim:latest
```

---

## Docker Image Details

### Architecture

**Multi-stage build** for optimal size and security:

```dockerfile
# Stage 1: Builder (golang:1.22-alpine)
- Compiles Go binary with static linking
- Includes build tools and dependencies
- ~300MB (discarded after build)

# Stage 2: Runtime (gcr.io/distroless/static:nonroot)
- Copies only the compiled binary and seed file
- Minimal distroless base image
- ~8.75MB final size
```

### Security Features

- **Non-root user:** Runs as UID 65532 (nonroot)
- **Static binary:** No dynamic library dependencies
- **Distroless base:** No shell, package manager, or unnecessary tools
- **Minimal attack surface:** Only contains app binary and seed data
- **Vulnerability scanning:** Compatible with container security scanners

### Performance Metrics

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Image build time (cold) | < 2 min | ~22s | ✅ 82% faster |
| Image build time (cached) | < 30s | ~4s | ✅ 87% faster |
| Final image size | < 50MB | 8.75MB | ✅ 82% smaller |
| Container startup time | < 5s | ~2s | ✅ 60% faster |
| Health check response | < 200ms | ~50ms | ✅ 75% faster |

### Image Layers

```bash
# Inspect image layers
docker history cursor-sim:latest

# Verify minimal size
docker images cursor-sim:latest

# Check for vulnerabilities (if using Docker Scout)
docker scout cves cursor-sim:latest
```

**Expected output:**
```
REPOSITORY    TAG       IMAGE ID       CREATED        SIZE
cursor-sim    latest    abc123def456   2 hours ago    8.75MB
```

---

## Using Custom Seed Files

### Option 1: Volume Mount (No Rebuild)

Mount your custom seed file into the container:

```bash
docker run --rm -p 8080:8080 \
  -v $(pwd)/testdata/custom_seed.json:/app/seed.json:ro \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  cursor-sim:latest
```

**Note:** The `:ro` flag mounts the file as read-only for added security.

### Option 2: Rebuild with Custom Seed

Modify the Dockerfile to copy your seed file, then rebuild:

```dockerfile
# In Dockerfile, replace default seed copy with:
COPY path/to/your/custom_seed.json /app/seed.json
```

```bash
# Rebuild with custom seed
docker build -t cursor-sim:custom services/cursor-sim

# Run custom image
docker run --rm -p 8080:8080 cursor-sim:custom
```

### Option 3: Build Argument

Pass seed file path at build time:

```dockerfile
# Add ARG in Dockerfile
ARG SEED_FILE=testdata/valid_seed.json
COPY ${SEED_FILE} /app/seed.json
```

```bash
# Build with custom seed
docker build --build-arg SEED_FILE=testdata/large_team_seed.json \
  -t cursor-sim:large-team services/cursor-sim
```

---

## Container Management

### View Logs

```bash
# Follow logs (foreground)
docker logs -f cursor-sim-local

# View last 100 lines
docker logs --tail 100 cursor-sim-local

# View logs with timestamps
docker logs -t cursor-sim-local
```

### Stop Container

```bash
# Graceful stop (sends SIGTERM, waits 10s, then SIGKILL)
docker stop cursor-sim-local

# Force kill (immediate)
docker kill cursor-sim-local
```

### Remove Container

```bash
# Remove stopped container
docker rm cursor-sim-local

# Force remove running container
docker rm -f cursor-sim-local
```

### View Running Containers

```bash
# View running containers
docker ps

# View all containers (including stopped)
docker ps -a

# View with custom format
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### Inspect Container

```bash
# Full container details
docker inspect cursor-sim-local

# Get specific field (e.g., IP address)
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cursor-sim-local

# View environment variables
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' cursor-sim-local
```

### Execute Commands in Container

```bash
# Note: cursor-sim uses distroless base, so no shell is available
# These commands will NOT work:
# docker exec -it cursor-sim-local sh  ❌ No shell

# Alternative: Use multi-stage build debug variant (if needed)
# Build with alpine as final stage for debugging
docker build --target builder -t cursor-sim:debug services/cursor-sim
docker run -it cursor-sim:debug sh
```

---

## Troubleshooting

### Container Starts but Health Check Fails

**Symptom:** Container runs but `curl http://localhost:8080/health` times out

**Diagnosis:**
```bash
# Check if port is already in use
lsof -i :8080  # macOS/Linux
netstat -an | grep 8080  # Windows

# Check container status
docker ps

# Check container logs
docker logs cursor-sim-local
```

**Solutions:**
```bash
# Use different port
PORT=8081 ./tools/docker-local.sh

# Or manually:
docker run --rm -p 8081:8080 cursor-sim:latest
```

### Build Fails with "go.mod: no such file or directory"

**Symptom:** Docker build fails during `COPY go.mod go.sum ./`

**Cause:** Building from wrong directory or .dockerignore excludes go.mod

**Solutions:**
```bash
# Ensure you're building from project root with correct context
cd /Users/jbellish/VSProjects/cursor-analytics-platform
docker build -t cursor-sim:latest services/cursor-sim

# Check .dockerignore isn't excluding go.mod
cat services/cursor-sim/.dockerignore
```

### Permission Denied Reading /app/seed.json

**Symptom:** Logs show "open /app/seed.json: permission denied"

**Cause:** Non-root user cannot read seed file (when using volume mounts)

**Solutions:**
```bash
# Make seed file readable by all users
chmod 644 testdata/custom_seed.json

# Verify permissions
ls -la testdata/custom_seed.json

# Mount with read-only flag
docker run --rm -p 8080:8080 \
  -v $(pwd)/testdata/custom_seed.json:/app/seed.json:ro \
  cursor-sim:latest
```

### Container Exits Immediately

**Symptom:** Container starts but exits with code 1

**Diagnosis:**
```bash
# Check logs for error details
docker logs cursor-sim-local

# Check exit code
docker inspect cursor-sim-local --format='{{.State.ExitCode}}'
```

**Common Errors:**

1. **"validation failed: seed path is required"**
   ```bash
   # Set CURSOR_SIM_SEED environment variable
   docker run -e CURSOR_SIM_SEED=/app/seed.json cursor-sim:latest
   ```

2. **"failed to load seed data"**
   ```bash
   # Check seed file path and format
   docker run -v $(pwd)/testdata/valid_seed.json:/app/seed.json:ro cursor-sim:latest
   ```

3. **"bind: address already in use"**
   ```bash
   # Use different port
   PORT=8081 ./tools/docker-local.sh
   ```

### Image Size Larger than Expected

**Symptom:** `docker images` shows size > 50MB

**Diagnosis:**
```bash
# Verify multi-stage build is working
docker history cursor-sim:latest | grep golang
# Should return nothing (golang only in builder stage)

# Check image layers
docker history cursor-sim:latest
```

**Solutions:**
```bash
# Rebuild without cache
docker build --no-cache -t cursor-sim:latest services/cursor-sim

# Or use automated script
REBUILD=true ./tools/docker-local.sh
```

### Cannot Connect to Container from Host

**Symptom:** `curl http://localhost:8080/health` connection refused

**Diagnosis:**
```bash
# Check container is running
docker ps | grep cursor-sim

# Check port mapping
docker port cursor-sim-local

# Check if container is listening
docker logs cursor-sim-local | grep "listening"
```

**Solutions:**
```bash
# Ensure port is properly mapped
docker run --rm -p 8080:8080 cursor-sim:latest
#                ^^^^^ Host port : Container port

# Try using container IP directly
CONTAINER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cursor-sim-local)
curl http://${CONTAINER_IP}:8080/health
```

---

## Development Workflow

### Local Development (No Docker)

For faster iteration during development:

```bash
cd services/cursor-sim

# Run tests
go test ./...

# Run with coverage
go test ./... -cover

# Build binary
go build -o bin/cursor-sim ./cmd/simulator

# Run locally
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -port 8080
```

### Docker Development Workflow

When you need to test Docker-specific behavior:

```bash
# 1. Make code changes
vim internal/api/handler.go

# 2. Rebuild image
REBUILD=true ./tools/docker-local.sh

# 3. Test endpoints
curl http://localhost:8080/health
curl -u cursor-sim-dev-key: http://localhost:8080/teams/members

# 4. Check logs
docker logs -f cursor-sim-local

# 5. Stop container
docker stop cursor-sim-local
```

### Testing Different Configurations

```bash
# Test with minimal data (fast iteration)
DAYS=7 VELOCITY=low ./tools/docker-local.sh

# Test with realistic data
DAYS=90 VELOCITY=medium ./tools/docker-local.sh

# Test with high volume (load testing)
DAYS=180 VELOCITY=high ./tools/docker-local.sh
```

---

## Production Deployment

### GCP Cloud Run

For production deployment to GCP Cloud Run, see the comprehensive guide:
- [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)

**Quick Deploy:**

```bash
# Set variables
PROJECT_ID=your-gcp-project
REGION=us-central1
SERVICE_NAME=cursor-sim
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:v2.0.0

# Build and push to Artifact Registry
docker build -t ${IMAGE_URI} services/cursor-sim
docker push ${IMAGE_URI}

# Deploy to Cloud Run
gcloud run deploy ${SERVICE_NAME} \
  --image ${IMAGE_URI} \
  --region ${REGION} \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --cpu 1 \
  --set-env-vars CURSOR_SIM_MODE=runtime,CURSOR_SIM_SEED=/app/seed.json,CURSOR_SIM_DAYS=90
```

### Docker Compose

For multi-service deployments (e.g., cursor-sim + analytics-core):

```yaml
# docker-compose.yml
version: '3.8'

services:
  cursor-sim:
    build: ./services/cursor-sim
    ports:
      - "8080:8080"
    environment:
      CURSOR_SIM_MODE: runtime
      CURSOR_SIM_SEED: /app/seed.json
      CURSOR_SIM_DAYS: 90
      CURSOR_SIM_VELOCITY: medium
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 3s
      retries: 3
    restart: unless-stopped

  analytics-core:
    build: ./services/cursor-analytics-core
    ports:
      - "4000:4000"
    environment:
      CURSOR_API_URL: http://cursor-sim:8080
    depends_on:
      cursor-sim:
        condition: service_healthy
```

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f cursor-sim

# Stop all services
docker-compose down
```

### Kubernetes

For Kubernetes deployments:

```yaml
# cursor-sim-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cursor-sim
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cursor-sim
  template:
    metadata:
      labels:
        app: cursor-sim
    spec:
      containers:
      - name: cursor-sim
        image: cursor-sim:v2.0.0
        ports:
        - containerPort: 8080
        env:
        - name: CURSOR_SIM_MODE
          value: "runtime"
        - name: CURSOR_SIM_SEED
          value: "/app/seed.json"
        - name: CURSOR_SIM_DAYS
          value: "90"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: cursor-sim
spec:
  selector:
    app: cursor-sim
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

```bash
# Deploy to Kubernetes
kubectl apply -f cursor-sim-deployment.yaml

# Check status
kubectl get pods -l app=cursor-sim
kubectl get svc cursor-sim

# View logs
kubectl logs -l app=cursor-sim --tail=100 -f
```

---

## Configuration Presets

Pre-configured setups for common scenarios:

### Development (Default)
```bash
./tools/docker-local.sh
# 90 days, medium velocity (~15 commits/day)
# Best for: General development and testing
```

### Quick Testing (Minimal Data)
```bash
DAYS=7 VELOCITY=low ./tools/docker-local.sh
# 7 days, low velocity (~5 commits/day)
# Best for: Fast iteration, unit testing
```

### Load Testing (High Volume)
```bash
DAYS=180 VELOCITY=high ./tools/docker-local.sh
# 180 days, high velocity (~30 commits/day)
# Best for: Performance testing, stress testing
```

### Demo/Presentation
```bash
DAYS=30 VELOCITY=medium PORT=8080 ./tools/docker-local.sh
# 30 days of realistic data
# Best for: Demonstrations, presentations
```

### CI/CD Pipeline
```bash
DETACH=true DAYS=30 VELOCITY=medium ./tools/docker-local.sh
# Background mode for automated testing
# Best for: Integration tests in CI/CD
```

---

## Files & Structure

```
services/cursor-sim/
├── Dockerfile              # Multi-stage Docker build
├── .dockerignore          # Build context exclusions
├── cmd/simulator/         # Main application entry point
├── internal/
│   ├── api/              # HTTP handlers
│   ├── config/           # Configuration management
│   ├── generator/        # Data generation logic
│   └── seed/             # Seed file loading
├── testdata/             # Test fixtures and seed files
├── README.md             # User guide
└── DOCKER.md             # This file
```

### .dockerignore

Optimizes build context by excluding unnecessary files:

```
# .dockerignore
.git
.github
.vscode
*.md
!README.md
bin/
coverage.out
*.test
```

---

## Platform Integration

cursor-sim (P4) is part of the Cursor Analytics Platform:

```
┌─────────────┐       ┌──────────────────┐       ┌─────────────┐
│ cursor-sim  │──────▶│ analytics-core   │──────▶│  viz-spa    │
│   (P4)      │ REST  │     (P5)         │GraphQL│    (P6)     │
│  Docker     │       │  Docker Compose  │       │  Local npm  │
│  Port 8080  │       │  Port 4000       │       │  Port 3000  │
└─────────────┘       └──────────────────┘       └─────────────┘
                             │
                             ▼
                      ┌─────────────┐
                      │ PostgreSQL  │
                      │  Port 5432  │
                      └─────────────┘
```

### Integration Status (January 8, 2026)

✅ **P4+P5 Integration**: Complete
- cursor-sim provides REST API for historical data
- cursor-analytics-core consumes and transforms to GraphQL
- Both services run in Docker for local development

---

## Related Documentation

- **User Guide:** [README.md](./README.md)
- **Technical Specification:** [SPEC.md](./SPEC.md)
- **Cloud Deployment:** [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)
- **Platform Architecture:** [docs/DESIGN.md](../../docs/DESIGN.md)
- **Data Contract Testing:** [docs/data-contract-testing.md](../../docs/data-contract-testing.md)

---

## Support

- **Issues:** Report bugs in project issue tracker
- **Questions:** See [README.md](./README.md) for general usage
- **Cloud Deploy:** See [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)

---

## License

Internal development tool - not for external distribution.
