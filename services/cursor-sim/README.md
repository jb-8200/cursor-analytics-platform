# Cursor Sim - Docker Deployment Guide

A Go-based HTTP server that simulates Cursor IDE usage data for testing and development.

## Quick Start with Docker

### Automated Script (Recommended)

The easiest way to run cursor-sim locally:

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

## Manual Docker Commands

### Build Image

```bash
docker build -t cursor-sim:latest services/cursor-sim
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

### Container Management

```bash
# View logs
docker logs -f cursor-sim-local

# Stop container
docker stop cursor-sim-local

# Remove container
docker rm cursor-sim-local

# View running containers
docker ps

# View all containers (including stopped)
docker ps -a
```

## Environment Variables

| Variable | Default | Options | Description |
|----------|---------|---------|-------------|
| `CURSOR_SIM_MODE` | `runtime` | `runtime`, `replay` | Simulation mode (runtime generates data, replay uses corpus) |
| `CURSOR_SIM_SEED` | `/app/seed.json` | file path | Path to seed data file with team/developer configuration |
| `CURSOR_SIM_DAYS` | `90` | number | Number of days of historical data to simulate |
| `CURSOR_SIM_VELOCITY` | `medium` | `low`, `medium`, `high` | Commit frequency (low=~5/day, medium=~15/day, high=~30/day) |
| `CURSOR_SIM_PORT` | `8080` | port number | HTTP server port |

## API Endpoints

Once running, test the API at `http://localhost:8080`:

### Health Check (No Authentication)

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "service": "cursor-sim",
  "version": "2.0.0"
}
```

### Team Members (Basic Auth Required)

```bash
curl -u cursor-sim-dev-key: http://localhost:8080/teams/members
```

### Commit History (Basic Auth Required)

```bash
curl -u cursor-sim-dev-key: http://localhost:8080/commits?limit=10
```

### Daily Statistics (Basic Auth Required)

```bash
curl -u cursor-sim-dev-key: http://localhost:8080/stats/daily?days=7
```

**Authentication:** Default API key is `cursor-sim-dev-key` (hardcoded in development mode)

## Using Custom Seed Files

### Option 1: Volume Mount (No Rebuild)

Mount your custom seed file into the container:

```bash
docker run --rm -p 8080:8080 \
  -v $(pwd)/testdata/custom_seed.json:/app/seed.json \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  cursor-sim:latest
```

### Option 2: Rebuild with Custom Seed

Modify the Dockerfile to copy your seed file, then rebuild:

```dockerfile
# In Dockerfile, add:
COPY path/to/your/custom_seed.json /app/seed.json
```

```bash
docker build -t cursor-sim:custom services/cursor-sim
docker run --rm -p 8080:8080 cursor-sim:custom
```

## Docker Image Details

### Architecture

**Multi-stage build:**
- **Builder stage**: `golang:1.22-alpine` - Compiles Go binary with static linking
- **Runtime stage**: `gcr.io/distroless/static:nonroot` - Minimal distroless image

### Security Features

- Runs as non-root user (UID 65532)
- Static binary with no dependencies
- Minimal attack surface (distroless base)
- No shell, package manager, or unnecessary tools
- Image scanned for vulnerabilities

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
```

Expected output:
```
REPOSITORY    TAG       IMAGE ID       CREATED        SIZE
cursor-sim    latest    abc123def456   2 hours ago    8.75MB
```

## Troubleshooting

### Container Starts but Health Check Fails

**Symptom:** Container runs but `curl http://localhost:8080/health` times out

**Solutions:**
```bash
# Check if port is already in use
lsof -i :8080  # macOS/Linux
netstat -an | grep 8080  # Windows

# Use different port
PORT=8081 ./tools/docker-local.sh

# Check container logs
docker logs cursor-sim-local
```

### Build Fails with "go.mod: no such file or directory"

**Symptom:** Docker build fails during `COPY go.mod go.sum ./`

**Solution:**
```bash
# Ensure you're building from project root
cd /Users/jbellish/VSProjects/cursor-analytics-platform
docker build -t cursor-sim:latest services/cursor-sim

# Check .dockerignore isn't excluding go.mod
cat services/cursor-sim/.dockerignore
```

### Permission Denied Reading /app/seed.json

**Symptom:** Logs show "open /app/seed.json: permission denied"

**Cause:** Non-root user cannot read seed file

**Solution:** Current Dockerfile uses `--chown=nonroot:nonroot` when copying files. If using volume mounts, ensure file permissions allow reading:

```bash
# Make seed file readable
chmod 644 testdata/custom_seed.json
```

### Container Exits Immediately

**Symptom:** Container starts but exits with code 1

**Solutions:**
```bash
# Check logs for error details
docker logs cursor-sim-local

# Common errors:
# "validation failed: seed path is required" → Set CURSOR_SIM_SEED
# "failed to load seed data" → Check seed file path and format
# "bind: address already in use" → Use different port with PORT=8081
```

### Image Size Larger than Expected

**Symptom:** `docker images` shows size > 50MB

**Solution:**
```bash
# Verify multi-stage build is working
docker history cursor-sim:latest | grep golang
# Should return nothing (golang only in builder stage)

# Rebuild without cache
REBUILD=true ./tools/docker-local.sh
```

## Development Workflow

### Local Development (No Docker)

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

```bash
# 1. Make code changes
# 2. Rebuild image
REBUILD=true ./tools/docker-local.sh

# 3. Test endpoints
curl http://localhost:8080/health

# 4. Check logs
docker logs -f cursor-sim-local
```

## Production Deployment

For production deployment to GCP Cloud Run, see:
- [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)

Quick deploy to Cloud Run:
```bash
# Build and push to Artifact Registry
PROJECT_ID=your-project
REGION=us-central1
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:v2.0.0

docker build -t ${IMAGE_URI} services/cursor-sim
docker push ${IMAGE_URI}

# Deploy to Cloud Run
gcloud run deploy cursor-sim \
  --image ${IMAGE_URI} \
  --region ${REGION} \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars CURSOR_SIM_MODE=runtime,CURSOR_SIM_SEED=/app/seed.json
```

## Configuration Presets

### Development (Default)
```bash
./tools/docker-local.sh
# 90 days, medium velocity (~15 commits/day)
```

### Quick Testing (Minimal Data)
```bash
DAYS=7 VELOCITY=low ./tools/docker-local.sh
# 7 days, low velocity (~5 commits/day)
```

### Load Testing (High Volume)
```bash
DAYS=180 VELOCITY=high ./tools/docker-local.sh
# 180 days, high velocity (~30 commits/day)
```

### Demo/Presentation
```bash
DAYS=30 VELOCITY=medium PORT=8080 ./tools/docker-local.sh
# 30 days of realistic data
```

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
└── README.md             # This file
```

## Platform Integration

cursor-sim (P4) is part of the Cursor Analytics Platform:

```
cursor-sim (P4) → cursor-analytics-core (P5) → cursor-viz-spa (P6)
  Docker          Docker Compose             Local npm dev
  Port 8080       Port 4000 (GraphQL)        Port 3000
                  Port 5432 (PostgreSQL)
```

### Integration Status (January 4, 2026)

✅ **P4+P5 Integration**: Complete
- cursor-sim provides REST API for historical data
- cursor-analytics-core consumes and transforms to GraphQL
- Both services run in Docker for local development

### Related Documentation

- **Platform Architecture**: [docs/DESIGN.md](../../docs/DESIGN.md)
- **Integration Guide**: [docs/INTEGRATION.md](../../docs/INTEGRATION.md)
- **Data Contract Testing**: [docs/data-contract-testing.md](../../docs/data-contract-testing.md)

## Support

- **Specification:** [services/cursor-sim/SPEC.md](./SPEC.md)
- **Cloud Deployment:** [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)
- **Issues:** Report bugs in project issue tracker

## License

Internal development tool - not for external distribution.
