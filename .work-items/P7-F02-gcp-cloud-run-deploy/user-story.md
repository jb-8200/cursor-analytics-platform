# User Story: GCP Cloud Run Deployment

**Feature ID**: P7-F02-gcp-cloud-run-deploy
**Created**: January 3, 2026
**Status**: Planning (Ready to Start)

---

## Story

**As a** researcher or platform operator
**I want** to deploy cursor-sim to Google Cloud Platform Cloud Run with a public URL
**So that** I can access the API from anywhere, share it with collaborators, and integrate it with analytics services without managing servers

## Background

Local Docker deployment (P7-F01) works well for development and testing, but production usage requires:

- **Public accessibility**: Analytics-core (P5) needs to consume cursor-sim API over HTTP
- **Zero server management**: No VMs to patch, scale, or monitor
- **Scale-to-zero cost model**: Pay only when processing requests
- **Automatic HTTPS**: Secure endpoints without certificate management
- **Built-in monitoring**: Logs, metrics, and traces in Google Cloud Console

Cloud Run provides:
- Fully managed serverless container platform
- Automatic scaling (including to zero)
- Built-in load balancing and HTTPS
- Pay-per-use pricing (no idle costs)
- Integration with GCP ecosystem (Artifact Registry, Cloud Build, IAM)

Use case: Deploy cursor-sim as a persistent API endpoint for cursor-analytics-core (P5) to consume, enabling the visualization dashboard (P6) to display real-time analytics.

## Acceptance Criteria

### AC-1: GCP Project and API Enablement

**Given** I have a GCP account and gcloud CLI installed
**When** I run the project setup commands:
```bash
gcloud services enable run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  --project ${PROJECT_ID}
```
**Then** all required APIs are enabled successfully
**And** I receive confirmation messages for each API
**And** the APIs are immediately available (no 5-minute delay)
**And** IAM permissions are correctly configured for my account

Required APIs:
- Cloud Run API (run.googleapis.com)
- Artifact Registry API (artifactregistry.googleapis.com)
- Cloud Build API (cloudbuild.googleapis.com)

### AC-2: Artifact Registry Repository Creation

**Given** Cloud Run API is enabled
**When** I create an Artifact Registry repository:
```bash
gcloud artifacts repositories create cursor-sim \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for cursor-sim" \
  --project ${PROJECT_ID}
```
**Then** the repository is created in region `us-central1`
**And** the repository accepts Docker image pushes
**And** I can authenticate Docker: `gcloud auth configure-docker us-central1-docker.pkg.dev`
**And** repository appears in Cloud Console Artifact Registry page

### AC-3: Docker Image Build and Push to Artifact Registry

**Given** Artifact Registry repository exists and Docker is authenticated
**When** I build and push the cursor-sim image:
```bash
IMAGE_URI=us-central1-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:v2.0.0
docker build -t ${IMAGE_URI} services/cursor-sim
docker push ${IMAGE_URI}
```
**Then** the image builds successfully (same Dockerfile as P7-F01)
**And** the image is pushed to Artifact Registry in < 3 minutes
**And** the image appears in Cloud Console with correct tags and metadata
**And** the image digest is stable and verifiable

Edge cases:
- First push takes longer (all layers uploaded)
- Subsequent pushes leverage layer caching
- Network interruptions retry automatically
- Failed pushes provide actionable error messages

### AC-4: Cloud Run Service Deployment

**Given** Docker image exists in Artifact Registry
**When** I deploy to Cloud Run:
```bash
gcloud run deploy cursor-sim \
  --project ${PROJECT_ID} \
  --region us-central1 \
  --image ${IMAGE_URI} \
  --port 8080 \
  --min-instances 0 \
  --max-instances 1 \
  --cpu 0.25 \
  --memory 512Mi \
  --allow-unauthenticated \
  --set-env-vars CURSOR_SIM_MODE=runtime,CURSOR_SIM_SEED=/app/seed.json,CURSOR_SIM_DAYS=90,CURSOR_SIM_VELOCITY=medium,CURSOR_SIM_PORT=8080
```
**Then** the service deploys successfully in < 2 minutes
**And** I receive a public URL: `https://cursor-sim-<hash>-uc.a.run.app`
**And** the service is in "Ready" status in Cloud Console
**And** first cold start completes in < 10 seconds

Configuration verification:
- CPU: 0.25 vCPU (cost-optimized for low traffic)
- Memory: 512Mi (sufficient for 2-3 developers)
- Min instances: 0 (scale to zero for cost savings)
- Max instances: 1 (single instance for MVP)
- Authentication: Allow unauthenticated (public API)

### AC-5: Service Health Verification

**Given** Cloud Run service is deployed
**When** I verify the deployment:
```bash
SERVICE_URL=$(gcloud run services describe cursor-sim \
  --platform managed \
  --region us-central1 \
  --format="value(status.url)" \
  --project ${PROJECT_ID})

curl -u cursor-sim-dev-key: ${SERVICE_URL}/health
```
**Then** the health endpoint returns 200 OK
**And** the response body contains `{"status":"healthy"}`
**And** subsequent requests have latency < 200ms (warm container)
**And** cold start requests complete in < 10 seconds

### AC-6: Environment Variable Updates Without Rebuild

**Given** a deployed Cloud Run service
**When** I update environment variables:
```bash
gcloud run services update cursor-sim \
  --region us-central1 \
  --set-env-vars CURSOR_SIM_DAYS=30,CURSOR_SIM_VELOCITY=high \
  --project ${PROJECT_ID}
```
**Then** the service updates without rebuilding the Docker image
**And** a new revision is created automatically
**And** traffic shifts to the new revision in < 1 minute
**And** the new configuration takes effect immediately
**And** old revision is preserved for rollback

### AC-7: Deployment Automation Script

**Given** the `tools/deploy-cursor-sim.sh` script
**When** I run:
```bash
PROJECT_ID=my-project REGION=us-central1 TAG=v2.0.0 ./tools/deploy-cursor-sim.sh
```
**Then** the script executes all deployment steps automatically:
1. Authenticate Docker with Artifact Registry
2. Build the Docker image
3. Push image to Artifact Registry
4. Deploy to Cloud Run
5. Output the service URL and health check command
**And** the script validates all required environment variables
**And** errors halt execution with clear messages
**And** the script is idempotent (safe to re-run)

Script features:
- Environment variable defaults (PROJECT_ID required, REGION defaults to us-central1)
- `set -euo pipefail` for strict error handling
- Color-coded output (green for success, red for errors)
- Progress indicators for long-running operations
- Final verification: health check on deployed URL

### AC-8: Cloud Console Monitoring and Logs

**Given** cursor-sim is deployed and receiving traffic
**When** I access Cloud Console → Cloud Run → cursor-sim service
**Then** I can view:
- **Logs**: Real-time container stdout/stderr in Cloud Logging
- **Metrics**: Request count, latency (p50, p95, p99), CPU/memory usage
- **Revisions**: List of all revisions with traffic splits
- **Traffic**: Current revision receiving 100% traffic
- **Environment**: All configured environment variables displayed
**And** logs are searchable and filterable by severity, timestamp
**And** metrics update every 60 seconds

### AC-9: Cost Optimization Verification

**Given** Cloud Run service with min-instances=0
**When** no requests are received for 15 minutes
**Then** the service scales to zero instances
**And** no CPU/memory charges accrue during idle time
**And** first request after scale-to-zero triggers cold start (< 10s)
**And** estimated monthly cost for 1000 requests/day is < $2 USD

Cost breakdown (for verification):
- Request charges: ~$0.40 per 1M requests
- CPU time: ~$0.024 per vCPU-hour
- Memory: ~$0.0025 per GiB-hour
- Free tier: 2M requests, 360K vCPU-seconds, 180K GiB-seconds per month

### AC-10: Rollback Procedure

**Given** a new revision causes issues
**When** I rollback to the previous revision:
```bash
gcloud run services update-traffic cursor-sim \
  --to-revisions=cursor-sim-00001-abc=100 \
  --region us-central1 \
  --project ${PROJECT_ID}
```
**Then** traffic shifts to the previous revision immediately
**And** the service returns to working state
**And** rollback completes in < 30 seconds
**And** no data loss or state corruption occurs

## Out of Scope

- **Custom domain mapping**: Use default Cloud Run URL (*.run.app), no custom DNS
- **GCS-based seed files**: Seed file baked into Docker image, no Cloud Storage integration
- **Multi-region deployment**: Single region (us-central1) only
- **Cloud CDN**: No CDN in front of Cloud Run
- **VPC connector**: Public internet only, no private VPC access
- **IAM-based authentication**: Allow unauthenticated public access (Basic Auth in app)
- **Cloud Armor (WAF)**: No Web Application Firewall
- **Cloud Build CI/CD**: Manual deployment script only, no automated CI/CD pipeline
- **Blue/green deployments**: Simple revision-based rollback only
- **Horizontal scaling beyond 1 instance**: Single instance for MVP

## Dependencies

- **P7-F01 completed**: Dockerfile and .dockerignore must exist and work locally
- **GCP account with billing enabled**: Free tier available but billing must be enabled
- **gcloud CLI installed and authenticated**: Version 400+ recommended
- **Docker installed**: For building images before push
- **Project-level IAM permissions**:
  - `roles/run.admin` - Deploy Cloud Run services
  - `roles/artifactregistry.admin` - Create repositories and push images
  - `roles/iam.serviceAccountUser` - Use default service account
- **docs/cursor-sim-cloud-run.md**: Complete deployment documentation

## Success Metrics

- **Deployment time**: < 5 minutes from script start to service URL ready
- **Cold start time**: < 10 seconds for first request after scale-to-zero
- **Warm latency**: < 200ms for health check and API endpoints
- **Image push time**: < 3 minutes to Artifact Registry
- **Availability**: 99.5% uptime (Cloud Run SLA)
- **Monthly cost**: < $2 for 1000 requests/day with scale-to-zero
- **Rollback time**: < 1 minute to previous revision
- **Documentation completeness**: All gcloud commands and troubleshooting steps documented

## Related Documents

- `services/cursor-sim/SPEC.md` - cursor-sim technical specification
- `docs/cursor-sim-cloud-run.md` - Deployment documentation (GCP sections)
- `.work-items/P7-F01-local-docker-deploy/` - Local Docker deployment (prerequisite)
- `.work-items/P7-F02-gcp-cloud-run-deploy/design.md` - Technical design for GCP deployment
- `.work-items/P7-F02-gcp-cloud-run-deploy/task.md` - Task breakdown and implementation steps
- `.claude/agents/cursor-sim-infra-dev.md` - Infrastructure agent for implementation
