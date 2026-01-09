#!/usr/bin/env bash
# Deployment script for cursor-sim to GCP Cloud Run
# Usage: ./tools/deploy-cursor-sim.sh [staging|production]
#
# Optional environment variables:
#   PROJECT_ID      - GCP project ID (default: cursor-sim)
#   REGION          - Cloud Run region (default: us-central1)
#   TAG             - Docker image tag (default: v2.0.1-YYYYMMDD)
#   CUSTOM_DOMAIN   - Custom domain to map (e.g., dox-a3.jishutech.io)
#
# Examples:
#   ./tools/deploy-cursor-sim.sh staging
#   CUSTOM_DOMAIN=dox-a3.jishutech.io ./tools/deploy-cursor-sim.sh staging
#   ./tools/deploy-cursor-sim.sh production
set -euo pipefail

# Configuration
PROJECT_ID=${PROJECT_ID:-cursor-sim}
REGION=${REGION:-us-central1}
ENVIRONMENT=${1:-staging}
TAG=${TAG:-v2.0.1-$(date +%Y%m%d)}
IMAGE_URI=${REGION}-docker.pkg.dev/${PROJECT_ID}/cursor-sim/cursor-sim:${TAG}
CUSTOM_DOMAIN=${CUSTOM_DOMAIN:-}

# Environment-specific settings
if [ "$ENVIRONMENT" = "production" ]; then
  SERVICE_NAME="cursor-sim-prod"
  MIN_INSTANCES=1
  MAX_INSTANCES=3
  CPU=0.5
  MEMORY=1Gi
  DAYS=180
  VELOCITY=high
else
  SERVICE_NAME="cursor-sim"
  MIN_INSTANCES=0
  MAX_INSTANCES=1
  CPU=0.25
  MEMORY=512Mi
  DAYS=90
  VELOCITY=medium
fi

echo "========================================="
echo "cursor-sim Cloud Run Deployment"
echo "========================================="
echo "Project:      ${PROJECT_ID}"
echo "Region:       ${REGION}"
echo "Environment:  ${ENVIRONMENT}"
echo "Service:      ${SERVICE_NAME}"
echo "Image:        ${IMAGE_URI}"
if [ -n "$CUSTOM_DOMAIN" ]; then
  echo "Custom Domain: ${CUSTOM_DOMAIN}"
fi
echo "========================================="
echo

# Authenticate Docker with Artifact Registry
echo "Configuring Docker authentication..."
gcloud auth configure-docker ${REGION}-docker.pkg.dev --quiet

# Build Docker image
echo
echo "Building Docker image..."
docker build -t "${IMAGE_URI}" services/cursor-sim

# Push to Artifact Registry
echo
echo "Pushing image to Artifact Registry..."
docker push "${IMAGE_URI}"

# Deploy to Cloud Run
echo
echo "Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image "${IMAGE_URI}" \
  --port 8080 \
  --min-instances ${MIN_INSTANCES} \
  --max-instances ${MAX_INSTANCES} \
  --cpu ${CPU} \
  --memory ${MEMORY} \
  --allow-unauthenticated \
  --set-env-vars CURSOR_SIM_MODE=runtime,CURSOR_SIM_SEED=/app/seed.json,CURSOR_SIM_DAYS=${DAYS},CURSOR_SIM_VELOCITY=${VELOCITY},CURSOR_SIM_PORT=8080

# Get service URL
echo
echo "Deployment complete!"
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} \
  --platform managed \
  --region "${REGION}" \
  --format="value(status.url)" \
  --project "${PROJECT_ID}")

# Configure custom domain mapping if specified
if [ -n "$CUSTOM_DOMAIN" ]; then
  echo
  echo "Configuring custom domain mapping..."

  # Check if beta components are installed
  if ! gcloud beta --version &> /dev/null; then
    echo "Installing gcloud beta components..."
    gcloud components install beta --quiet
  fi

  # Create domain mapping (or update if exists)
  if gcloud beta run domain-mappings describe --domain="${CUSTOM_DOMAIN}" --region="${REGION}" --project="${PROJECT_ID}" &> /dev/null; then
    echo "Domain mapping already exists for ${CUSTOM_DOMAIN}"
  else
    echo "Creating domain mapping for ${CUSTOM_DOMAIN}..."
    gcloud beta run domain-mappings create \
      --service="${SERVICE_NAME}" \
      --domain="${CUSTOM_DOMAIN}" \
      --region="${REGION}" \
      --project="${PROJECT_ID}"

    echo
    echo "⚠️  IMPORTANT: Ensure your DNS has a CNAME record:"
    echo "   Name:  ${CUSTOM_DOMAIN}"
    echo "   Value: ghs.googlehosted.com"
    echo
    echo "SSL certificate will be provisioned automatically (may take 10-15 minutes)"
  fi
fi

echo
echo "========================================="
echo "Service URL: ${SERVICE_URL}"
if [ -n "$CUSTOM_DOMAIN" ]; then
  echo "Custom URL:  https://${CUSTOM_DOMAIN}"
fi
echo "========================================="
echo
echo "Test endpoints:"
if [ -n "$CUSTOM_DOMAIN" ]; then
  echo "  Health:  curl https://${CUSTOM_DOMAIN}/health"
  echo "  Members: curl -u cursor-sim-dev-key: https://${CUSTOM_DOMAIN}/teams/members"
  echo "  Commits: curl -u cursor-sim-dev-key: https://${CUSTOM_DOMAIN}/analytics/ai-code/commits?page_size=5"
else
  echo "  Health:  curl ${SERVICE_URL}/health"
  echo "  Members: curl -u cursor-sim-dev-key: ${SERVICE_URL}/teams/members"
  echo "  Commits: curl -u cursor-sim-dev-key: ${SERVICE_URL}/analytics/ai-code/commits?page_size=5"
fi
echo
echo "View logs:"
echo "  gcloud run services logs read ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID}"
echo
