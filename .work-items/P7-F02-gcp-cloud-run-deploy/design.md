# Technical Design: GCP Cloud Run Deployment

**Feature ID**: P7-F02-gcp-cloud-run-deploy
**Created**: January 3, 2026
**Status**: Planning

---

## Overview

Deploy cursor-sim to Google Cloud Platform Cloud Run using the Docker image from P7-F01, providing a publicly accessible API endpoint for cursor-analytics-core (P5) integration.

## Architecture

### Deployment Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                     GCP Project                               │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Artifact Registry (us-central1)                        │  │
│  │                                                        │  │
│  │  Repository: cursor-sim                               │  │
│  │  Image: cursor-sim:v2.0.0                             │  │
│  │  Size: ~22MB                                          │  │
│  │  Layers: cached and deduplicated                      │  │
│  └─────────────────────┬──────────────────────────────────┘  │
│                        │                                      │
│                        │ gcloud run deploy                    │
│                        ▼                                      │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Cloud Run Service: cursor-sim (us-central1)           │  │
│  │                                                        │  │
│  │  ┌──────────────────────────────────────────────────┐ │  │
│  │  │ Revision: cursor-sim-00001-abc                   │ │  │
│  │  │                                                  │ │  │
│  │  │  Container:                                      │ │  │
│  │  │  - Image: us-central1-docker.pkg.dev/.../cursor-sim│ │  │
│  │  │  - Port: 8080                                    │ │  │
│  │  │  - CPU: 0.25 vCPU                                │ │  │
│  │  │  - Memory: 512Mi                                 │ │  │
│  │  │  - User: nonroot (65532)                         │ │  │
│  │  │                                                  │ │  │
│  │  │  Env:                                            │ │  │
│  │  │  - CURSOR_SIM_MODE=runtime                       │ │  │
│  │  │  - CURSOR_SIM_SEED=/app/seed.json                │ │  │
│  │  │  - CURSOR_SIM_DAYS=90                            │ │  │
│  │  │  - CURSOR_SIM_VELOCITY=medium                    │ │  │
│  │  │  - CURSOR_SIM_PORT=8080                          │ │  │
│  │  │                                                  │ │  │
│  │  │  Scaling:                                        │ │  │
│  │  │  - Min: 0 (scale-to-zero)                        │ │  │
│  │  │  - Max: 1                                        │ │  │
│  │  │  - Concurrency: 80                               │ │  │
│  │  └──────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  Traffic: 100% to latest revision                     │  │
│  │  Authentication: Allow unauthenticated                 │  │
│  │  URL: https://cursor-sim-<hash>-uc.a.run.app          │  │
│  └────────────────────────────────────────────────────────┘  │
│                        │                                      │
│                        │ HTTPS                                │
└────────────────────────┼──────────────────────────────────────┘
                         │
                         ▼
            ┌─────────────────────────┐
            │  Public Internet        │
            │                         │
            │  - analytics-core (P5)  │
            │  - Browser/curl         │
            │  - Researchers          │
            └─────────────────────────┘
```

### Request Flow

```
1. Client Request
   │
   └─► https://cursor-sim-<hash>-uc.a.run.app/api/team/members
        │
        ▼
2. Cloud Run Load Balancer (HTTPS termination)
   │
   ├─► If no instances: Cold start (< 10s)
   │   └─► Start new container instance
   │
   └─► If warm: Route to existing instance (< 200ms)
        │
        ▼
3. Container Instance
   │
   ├─► Basic Auth check (Authorization: Basic cursor-sim-dev-key)
   │
   └─► Application handler
        │
        └─► Return JSON response
             │
             ▼
4. Response sent back through Load Balancer
```

## Implementation Details

### 1. GCP Prerequisites

**Required APIs**:
- `run.googleapis.com` - Cloud Run service
- `artifactregistry.googleapis.com` - Docker image registry
- `cloudbuild.googleapis.com` - Optional for CI/CD

**Required IAM Roles** (for deploying user):
- `roles/run.admin` - Deploy and manage Cloud Run services
- `roles/artifactregistry.admin` - Push images to Artifact Registry
- `roles/iam.serviceAccountUser` - Use service accounts

**gcloud CLI Configuration**:
```bash
# Set project
gcloud config set project ${PROJECT_ID}

# Set default region
gcloud config set run/region us-central1

# Authenticate Docker
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### 2. Artifact Registry Setup

**Design Decisions**:

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Location | `us-central1` | Low latency for US users, cost-effective |
| Repository name | `cursor-sim` | Matches service name, clear purpose |
| Format | `docker` | Standard Docker image format |
| Cleanup policy | None (retain all) | Small images, low storage cost |

**Repository URI Format**:
```
${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMAGE}:${TAG}

Example:
us-central1-docker.pkg.dev/my-project/cursor-sim/cursor-sim:v2.0.0
```

**Creation Command**:
```bash
gcloud artifacts repositories create cursor-sim \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for cursor-sim application" \
  --project ${PROJECT_ID}
```

### 3. Image Build and Push Strategy

**Build Strategy**: Use same Dockerfile from P7-F01 (no changes needed)

**Tagging Strategy**:
- **Latest**: `latest` - Always points to most recent build
- **Semantic version**: `v2.0.0`, `v2.1.0` - Specific releases
- **Git SHA**: `abc123f` - Git commit for traceability
- **Environment**: `prod`, `staging`, `dev` - Environment-specific tags

**Push Process**:
```bash
# Tag with multiple tags
docker build -t cursor-sim:local services/cursor-sim

docker tag cursor-sim:local \
  ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:latest

docker tag cursor-sim:local \
  ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:v2.0.0

docker tag cursor-sim:local \
  ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:$(git rev-parse --short HEAD)

# Push all tags
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:latest
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:v2.0.0
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:$(git rev-parse --short HEAD)
```

### 4. Cloud Run Service Configuration

**Resource Configuration**:

| Resource | Value | Rationale |
|----------|-------|-----------|
| **CPU** | 0.25 vCPU | Minimal, sufficient for 2-3 developers |
| **Memory** | 512Mi | Go app is memory-efficient, room for caching |
| **Min instances** | 0 | Scale-to-zero for cost savings |
| **Max instances** | 1 | MVP, single instance sufficient |
| **Concurrency** | 80 | Default, allows 80 concurrent requests per instance |
| **Timeout** | 300s (5 min) | Default, sufficient for API responses |

**Scaling Behavior**:
- **Scale-to-zero**: After 15 minutes of no traffic, instance terminates
- **Cold start**: First request after scale-to-zero takes ~5-10s
- **Warm requests**: Subsequent requests < 200ms latency
- **Autoscaling trigger**: CPU utilization > 60% or concurrency > 80

**Authentication**:
- **Allow unauthenticated**: Public API, no IAM required
- **Application-level auth**: Basic Auth with hardcoded key `cursor-sim-dev-key`
- **Future**: IAM-based auth for production

**Environment Variables** (from P7-F01):
```bash
--set-env-vars \
  CURSOR_SIM_MODE=runtime,\
  CURSOR_SIM_SEED=/app/seed.json,\
  CURSOR_SIM_DAYS=90,\
  CURSOR_SIM_VELOCITY=medium,\
  CURSOR_SIM_PORT=8080
```

### 5. Deployment Script Design

**Location**: `tools/deploy-cursor-sim.sh`

**Purpose**: Automate entire deployment pipeline from build to verification

**Features**:
- Environment variable validation (PROJECT_ID required)
- Idempotent (safe to re-run)
- Error handling with `set -euo pipefail`
- Color-coded output
- Progress indicators
- Health verification
- Output service URL for easy access

**Script Flow**:
```
1. Validate environment variables
   ├─ PROJECT_ID (required)
   ├─ REGION (default: us-central1)
   ├─ TAG (default: git SHA)
   └─ IMAGE_URI (computed)

2. Authenticate Docker
   └─ gcloud auth configure-docker ${REGION}-docker.pkg.dev

3. Build Docker image
   └─ docker build -t ${IMAGE_URI} services/cursor-sim

4. Push to Artifact Registry
   └─ docker push ${IMAGE_URI}

5. Deploy to Cloud Run
   └─ gcloud run deploy cursor-sim ...

6. Get service URL
   └─ gcloud run services describe cursor-sim --format="value(status.url)"

7. Verify health
   └─ curl ${SERVICE_URL}/health

8. Output success message
   └─ Service URL, health check command, example API call
```

**Error Handling**:
- Missing PROJECT_ID → exit with error message
- Docker build failure → exit, show build logs
- Image push failure → exit, suggest authentication check
- Deploy failure → exit, show gcloud error output
- Health check failure → warning (don't exit), show logs

### 6. Update Strategy and Rollback

**Update Process** (environment variables only):
```bash
# Update without rebuild
gcloud run services update cursor-sim \
  --region us-central1 \
  --set-env-vars CURSOR_SIM_DAYS=30,CURSOR_SIM_VELOCITY=high \
  --project ${PROJECT_ID}

# Creates new revision, shifts traffic automatically
```

**Rollback Process**:
```bash
# List revisions
gcloud run revisions list \
  --service cursor-sim \
  --region us-central1 \
  --project ${PROJECT_ID}

# Rollback to specific revision
gcloud run services update-traffic cursor-sim \
  --to-revisions=cursor-sim-00001-abc=100 \
  --region us-central1 \
  --project ${PROJECT_ID}
```

**Traffic Splitting** (for gradual rollout, future):
```bash
# 50/50 split
gcloud run services update-traffic cursor-sim \
  --to-revisions=cursor-sim-00002-xyz=50,cursor-sim-00001-abc=50 \
  --region us-central1 \
  --project ${PROJECT_ID}
```

## Monitoring and Observability

### Cloud Logging

**Log Types**:
- **Container logs**: Application stdout/stderr
- **Request logs**: HTTP method, path, status, latency
- **System logs**: Cold starts, scaling events, errors

**Accessing Logs**:
```bash
# Recent logs
gcloud run services logs read cursor-sim \
  --region us-central1 \
  --limit 50 \
  --project ${PROJECT_ID}

# Tail logs
gcloud run services logs tail cursor-sim \
  --region us-central1 \
  --project ${PROJECT_ID}
```

**Cloud Console**: Cloud Run → cursor-sim → Logs tab

### Cloud Monitoring

**Metrics Available**:
- Request count (total, by status code)
- Request latency (p50, p95, p99)
- Container CPU utilization
- Container memory utilization
- Instance count (current running instances)
- Billable time

**Accessing Metrics**: Cloud Console → Cloud Run → cursor-sim → Metrics tab

### Alerts (future enhancement)**:
- High error rate (> 5% of requests return 5xx)
- High latency (p99 > 1s)
- Cold start frequency (> 10 per hour)
- Cost threshold (> $10/month)

## Cost Estimation

**Pricing Model** (as of 2026-01-03):

| Resource | Price | Usage (1000 req/day) | Monthly Cost |
|----------|-------|----------------------|--------------|
| **Requests** | $0.40 per 1M | 30K requests | $0.01 |
| **CPU** | $0.024 per vCPU-hour | ~1 hour billed | $0.024 |
| **Memory** | $0.0025 per GiB-hour | ~0.5 GiB-hours | $0.001 |
| **Network egress** | $0.12 per GB | ~1 GB | $0.12 |
| **Total** | - | - | **~$0.16** |

**Free Tier** (resets monthly):
- 2M requests
- 360K vCPU-seconds
- 180K GiB-seconds
- 1GB network egress

**Estimated monthly cost**: $0-$2 for MVP usage (well within free tier)

## Security Considerations

1. **HTTPS termination**: Cloud Run provides automatic HTTPS with valid certificates
2. **IAM roles**: Principle of least privilege, only grant necessary roles
3. **Service account**: Use default Compute Engine service account (no additional permissions needed)
4. **Network**: Public internet only, no VPC connector (simpler, sufficient for MVP)
5. **Secrets**: No sensitive data in environment variables (hardcoded API key acceptable for dev)
6. **Image scanning**: Artifact Registry automatically scans images for vulnerabilities
7. **Non-root user**: Container runs as UID 65532 (from Dockerfile)

## Performance Characteristics

| Metric | Target | Expected | Notes |
|--------|--------|----------|-------|
| Cold start time | < 10s | ~5-7s | First request after scale-to-zero |
| Warm request latency | < 200ms | ~45ms | Health check endpoint |
| API request latency | < 500ms | ~150ms | /api/team/members endpoint |
| Image push time | < 5 min | ~2-3 min | First push, layers cached after |
| Deployment time | < 5 min | ~2-3 min | Build + push + deploy |
| Health check response | < 1s | ~50ms | Simple health endpoint |

## Testing Strategy

### Pre-Deployment Tests
1. **Local Docker test**: Verify image works locally (P7-F01)
2. **Environment variable test**: Verify all env vars are set correctly
3. **Seed file test**: Verify seed file is accessible in container

### Post-Deployment Tests
1. **Health check**: `curl ${SERVICE_URL}/health` returns 200
2. **API test**: `curl -u cursor-sim-dev-key: ${SERVICE_URL}/api/team/members` returns JSON
3. **Authentication test**: Request without auth returns 401
4. **Load test** (optional): ab -n 100 -c 10 ${SERVICE_URL}/health
5. **Cold start test**: Wait 15 min, verify first request completes in < 10s

### Integration Tests (P5 Analytics-Core)
1. Configure cursor-analytics-core to point to Cloud Run URL
2. Verify data fetching works
3. Verify pagination works
4. Verify authentication works

## Alternatives Considered

| Alternative | Pros | Cons | Decision |
|-------------|------|------|----------|
| **Cloud Run** | Serverless, scale-to-zero, managed | Cold starts, pay-per-use | ✅ Selected |
| App Engine Standard | Similar to Cloud Run, auto-scaling | More expensive, less flexible | Rejected |
| Compute Engine VM | Full control, no cold starts | Manual management, always-on cost | Rejected - overkill |
| GKE (Kubernetes) | Full orchestration, multi-service | Complex, expensive, overkill for single service | Rejected |
| Cloud Functions | Simpler, HTTP functions | Not designed for long-running services | Rejected |

## Dependencies

- **P7-F01 completed**: Dockerfile must exist and work
- **GCP project with billing enabled**: Free tier sufficient
- **gcloud CLI installed**: Version 400+ recommended
- **Docker installed**: For local builds
- **Project-level IAM permissions**: run.admin, artifactregistry.admin

## Rollout Plan

1. **GCP-01**: Enable GCP APIs and create Artifact Registry
2. **GCP-02**: Build and push Docker image to Artifact Registry
3. **GCP-03**: Deploy to Cloud Run with minimal config
4. **GCP-04**: Verify deployment and test endpoints
5. **GCP-05**: Create deployment automation script
6. **GCP-06**: Document usage and troubleshooting
7. **GCP-07**: Integration test with P5 (analytics-core)

## Success Criteria

- ✅ Service deploys successfully to Cloud Run
- ✅ Public URL is accessible via HTTPS
- ✅ Health check returns 200 OK
- ✅ API endpoints return valid JSON with auth
- ✅ Cold start time < 10 seconds
- ✅ Warm latency < 200ms
- ✅ Deployment script completes without errors
- ✅ Monthly cost < $2 (within free tier)
- ✅ Logs and metrics accessible in Cloud Console
- ✅ P5 (analytics-core) can consume API successfully
