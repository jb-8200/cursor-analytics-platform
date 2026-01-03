# Cursor Sim Containerization & Cloud Run Deployment

Design notes and runbooks for packaging `services/cursor-sim` as a Docker image and deploying it to GCP Cloud Run (public URL, no custom DNS).

## Scope & Assumptions
- Service: long-running Go HTTP server on port 8080 with `/health`.
- Environment vars supported (per `internal/config/config.go`): `CURSOR_SIM_MODE` (`runtime|replay`), `CURSOR_SIM_SEED` (path; required for runtime), `CURSOR_SIM_CORPUS` (replay, future), `CURSOR_SIM_DAYS`, `CURSOR_SIM_VELOCITY`, `CURSOR_SIM_PORT`.
- Default Basic Auth key is hardcoded to `cursor-sim-dev-key` in `cmd/simulator/main.go`.
- Seed handling: bake into the image (simplest) or fetch from GCS at startup.

## Containerization Design
- Base: multi-stage. Builder on `golang:1.22-alpine` (with `GOTOOLCHAIN=auto`, `CGO_ENABLED=0`) and runtime on `gcr.io/distroless/static:nonroot` (or `cgr.dev/chainguard/static`).
- Ports: expose 8080.
- Entrypoint: the built binary.
- Health: `/health` (no auth).
- `.dockerignore` (recommended): `.git`, `bin/`, `coverage.out`, `coverage.html`, `testdata/` (unless used as seed), `**/*.log`, `.DS_Store`.
- Runtime config is not baked into the image—change behavior via env vars/flags at run/deploy time (no rebuild needed).

### Sample Dockerfile (summary)
```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.22-alpine AS builder
ENV CGO_ENABLED=0 GOTOOLCHAIN=auto
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/cursor-sim ./cmd/simulator

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/cursor-sim /app/cursor-sim
# COPY seed.json /app/seed.json  # uncomment if baking seed
EXPOSE 8080
ENTRYPOINT ["/app/cursor-sim"]
```

## Step 1: Enable Necessary APIs
- Console: API Library → enable Cloud Run, Artifact Registry, Cloud Build, Cloud Storage.
- gcloud:
```
PROJECT_ID=<PROJECT_ID>
gcloud services enable run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  storage.googleapis.com \
  --project ${PROJECT_ID}
```

## Step 2: Create Artifact Registry Repository
- Console: Artifact Registry → Create Repository → name `cursor-sim`, format Docker, mode Standard, region `us-central1`.
- gcloud:
```
PROJECT_ID=<PROJECT_ID>
gcloud artifacts repositories create cursor-sim \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for cursor-sim application" \
  --project ${PROJECT_ID}
```

## Step 3: Build and Push Docker Image
Authenticate Docker, build, push:
```
PROJECT_ID=<PROJECT_ID>
REGION=us-central1
IMAGE=cursor-sim
TAG=v2.0.0
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/${IMAGE}:${TAG}

gcloud auth configure-docker ${REGION}-docker.pkg.dev
docker build -t ${IMAGE_URI} services/cursor-sim
docker push ${IMAGE_URI}
```

## Step 4: (Optional) Create GCS Bucket for seed.json
- Console: Cloud Storage → Create bucket (e.g., `cursor-sim-seed-bucket`), upload `seed.json`.
- CLI:
```
PROJECT_ID=<PROJECT_ID>
BUCKET=<UNIQUE_BUCKET_NAME>   # must be globally unique
gsutil mb -p ${PROJECT_ID} -l US gs://${BUCKET}/
gsutil cp seed.json gs://${BUCKET}/seed.json
```

## GCP Quick Start
Prereqs: `gcloud` CLI, Docker, a GCP project, Artifact Registry API enabled.

1) Auth & project:
```
gcloud auth login
gcloud config set project <PROJECT_ID>
gcloud auth configure-docker <REGION>-docker.pkg.dev
```

2) Build & push to Artifact Registry:
```
PROJECT_ID=<PROJECT_ID>
REGION=us-central1
IMAGE=cursor-sim
TAG=v2.0.0
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/${IMAGE}:${TAG}

docker build -t ${IMAGE_URI} services/cursor-sim
docker push ${IMAGE_URI}
```

3) Deploy to Cloud Run (public, scale-to-zero):
```
SEED_PATH=/app/seed.json   # adjust if you fetch from GCS
DAYS=90
VELOCITY=medium
PROJECT_ID=<PROJECT_ID>
REGION=us-central1
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:${TAG}

gcloud run deploy cursor-sim \
  --project ${PROJECT_ID} \
  --region ${REGION} \
  --image ${IMAGE_URI} \
  --port 8080 \
  --min-instances 0 \
  --max-instances 1 \
  --cpu 0.25 \
  --memory 512Mi \
  --allow-unauthenticated \
  --set-env-vars CURSOR_SIM_MODE=runtime,\
CURSOR_SIM_SEED=${SEED_PATH},\
CURSOR_SIM_DAYS=${DAYS},\
CURSOR_SIM_VELOCITY=${VELOCITY},\
CURSOR_SIM_PORT=8080
```

Result: Cloud Run gives a URL like `https://cursor-sim-<hash>-uc.a.run.app/` you can hit directly (no DNS needed). Health: `GET /health`.

## Step 5: Grant Cloud Run Access to GCS (if seed in GCS)
```
CLOUD_RUN_SERVICE=cursor-sim
CLOUD_RUN_REGION=us-central1
PROJECT_ID=<PROJECT_ID>
BUCKET=<UNIQUE_BUCKET_NAME>
CLOUD_RUN_SA_EMAIL=$(gcloud run services describe ${CLOUD_RUN_SERVICE} \
  --platform managed \
  --region ${CLOUD_RUN_REGION} \
  --format="value(spec.template.spec.serviceAccountName)" \
  --project ${PROJECT_ID})
gsutil iam ch serviceAccount:${CLOUD_RUN_SA_EMAIL}:objectViewer gs://${BUCKET}
```
Add an init step or entrypoint wrapper to `gsutil cp gs://${BUCKET}/seed.json /app/seed.json` before starting.

## Step 6: Confirm Service URL and Monitor
```
gcloud run services describe cursor-sim \
  --platform managed \
  --region us-central1 \
  --format="value(status.url)" \
  --project <PROJECT_ID>
```
Test: `curl -u cursor-sim-dev-key: https://<service-url>/health` (or with an identity token if private).  
Console: Cloud Run → Logs for stdout/stderr; Metrics for latency/CPU/requests. Use Cloud Monitoring/Logging for alerts.

### Change Params Without Rebuilding
- Local Docker: restart with new envs/flags, same image:
  ```
  docker run --rm -p 8080:8080 \
    -e CURSOR_SIM_MODE=replay \
    -e CURSOR_SIM_CORPUS=/app/events.parquet \
    -e CURSOR_SIM_PORT=8080 \
    cursor-sim:latest
  ```
- Cloud Run: update envs on the existing image (creates a new revision):
  ```
  gcloud run services update cursor-sim \
    --region ${REGION} \
    --set-env-vars CURSOR_SIM_MODE=replay,CURSOR_SIM_CORPUS=/app/events.parquet,CURSOR_SIM_PORT=8080
  ```
  or redeploy with `gcloud run deploy ... --set-env-vars ...`. No rebuild required; this restarts with the new params.

## Programmatic CLI Deploy Script (bash)
Save as `tools/deploy-cursor-sim.sh` if desired; run from repo root.
```bash
#!/usr/bin/env bash
set -euo pipefail

PROJECT_ID=${PROJECT_ID:-your-project}
REGION=${REGION:-us-central1}
TAG=${TAG:-v2.0.0}
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:${TAG}
SEED_PATH=${SEED_PATH:-/app/seed.json}
DAYS=${DAYS:-90}
VELOCITY=${VELOCITY:-medium}

gcloud auth configure-docker ${REGION}-docker.pkg.dev
docker build -t "${IMAGE_URI}" services/cursor-sim
docker push "${IMAGE_URI}"

gcloud run deploy cursor-sim \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image "${IMAGE_URI}" \
  --port 8080 \
  --min-instances 0 \
  --max-instances 1 \
  --cpu 0.25 \
  --memory 512Mi \
  --allow-unauthenticated \
  --set-env-vars CURSOR_SIM_MODE=runtime,\
CURSOR_SIM_SEED=${SEED_PATH},\
CURSOR_SIM_DAYS=${DAYS},\
CURSOR_SIM_VELOCITY=${VELOCITY},\
CURSOR_SIM_PORT=8080
```

## Manual Docker Deploy (Local)

### Quick Start with Automated Script

Use the provided script for easy local testing:

```bash
# Default configuration (runtime mode, 90 days, medium velocity)
./tools/docker-local.sh

# Custom configuration
DAYS=30 VELOCITY=high ./tools/docker-local.sh

# Run in background (detached mode)
DETACH=true ./tools/docker-local.sh

# Force rebuild image
REBUILD=true ./tools/docker-local.sh

# Custom port (if 8080 is in use)
PORT=8081 ./tools/docker-local.sh
```

The script automatically:
- Builds the image if it doesn't exist
- Starts the container with health check verification
- Displays service URL and example commands
- Cleans up gracefully on exit (foreground mode)

### Manual Docker Commands

If you prefer running Docker commands directly:

```bash
# Build image
docker build -t cursor-sim:latest services/cursor-sim

# Run container (foreground with logs)
docker run --rm -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  -e CURSOR_SIM_VELOCITY=medium \
  cursor-sim:latest

# Run container (background/detached)
docker run -d --name cursor-sim-local -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  -e CURSOR_SIM_VELOCITY=medium \
  cursor-sim:latest

# View logs
docker logs -f cursor-sim-local

# Stop container
docker stop cursor-sim-local
docker rm cursor-sim-local
```

### Volume Mounting Custom Seed Files

To use a custom seed file without rebuilding the image:

```bash
docker run --rm -p 8080:8080 \
  -v $(pwd)/testdata/custom_seed.json:/app/seed.json \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  cursor-sim:latest
```

### Performance Benchmarks

Measured on Apple Silicon (M-series):

| Metric | Target | Actual |
|--------|--------|--------|
| Image build time (cold) | < 2 min | ~22s |
| Image build time (cached) | < 30s | ~4s |
| Final image size | < 50MB | 8.75MB |
| Container startup time | < 5s | ~2s |
| Health check response | < 200ms | ~50ms |

### Troubleshooting

#### Container starts but health check fails

**Symptom**: Container runs but `curl http://localhost:8080/health` times out or fails.

**Possible causes**:
- Port already in use: Try a different port with `PORT=8081 ./tools/docker-local.sh`
- Environment variables missing: Check logs with `docker logs cursor-sim-local`
- Seed file not accessible: Verify `/app/seed.json` has correct permissions

**Solution**:
```bash
# Check container logs
docker logs cursor-sim-local

# Verify environment variables
docker inspect cursor-sim-local --format='{{.Config.Env}}'

# Check if port is available
lsof -i :8080  # macOS/Linux
netstat -an | grep 8080  # Windows
```

#### Build fails with "go.mod: no such file or directory"

**Symptom**: Docker build fails during `COPY go.mod go.sum ./`

**Possible causes**:
- Running docker build from wrong directory
- .dockerignore excluding too many files

**Solution**:
```bash
# Build from repository root
cd /path/to/cursor-analytics-platform
docker build -t cursor-sim:latest services/cursor-sim

# Verify build context
docker build -t cursor-sim:latest --no-cache services/cursor-sim
```

#### Image size larger than expected (> 50MB)

**Symptom**: `docker images cursor-sim:latest` shows size > 50MB

**Possible causes**:
- Multi-stage build not working
- .dockerignore not excluding build artifacts

**Solution**:
```bash
# Check image layers
docker history cursor-sim:latest

# Verify golang builder is not in final image
docker history cursor-sim:latest | grep golang
# Should return nothing (golang only in builder stage)

# Rebuild with --no-cache
REBUILD=true ./tools/docker-local.sh
```

#### Permission denied reading /app/seed.json

**Symptom**: Logs show "open /app/seed.json: permission denied"

**Cause**: Non-root user cannot read seed file (fixed in current Dockerfile)

**Solution**: Ensure Dockerfile uses `--chown=nonroot:nonroot` when copying files:
```dockerfile
COPY --chown=nonroot:nonroot testdata/valid_seed.json /app/seed.json
```

#### Container exits immediately with code 1

**Symptom**: Container starts but exits immediately

**Possible causes**:
- Missing required environment variables (MODE, SEED)
- Invalid seed file path
- Port already bound

**Solution**:
```bash
# Check why container exited
docker logs cursor-sim-local

# Common errors and fixes:
# "validation failed: seed path is required" → Set CURSOR_SIM_SEED
# "failed to load seed data" → Check seed file path and permissions
# "bind: address already in use" → Use different port
```

## Cloud Console: Access & Monitoring
- Console login: https://console.cloud.google.com/ (select the correct project).
- Cloud Run service page: view revisions, traffic, and the service URL.
- Logs: Cloud Run → your service → Logs (stdout/stderr from the container).
- Metrics: Cloud Run → your service → Metrics (latency, CPU/mem, requests).
- IAM & access: Allow unauthenticated for public; otherwise bind `roles/run.invoker` to specific users/service accounts.

## Variations & Notes
- Seed from GCS: give the Cloud Run service account `roles/storage.objectViewer`; add a tiny init script to download `gs://bucket/seed.json` to `/app/seed.json` before starting.
- Replay mode (future): set `CURSOR_SIM_MODE=replay` and `CURSOR_SIM_CORPUS` to the mounted path; bump memory if loading Parquet into memory.
- Resource tuning: for 2–3 developers, `--cpu 0.25 --memory 512Mi` and `--max-instances 1` is sufficient; keep `--min-instances 0` to minimize cost. 
