#!/usr/bin/env bash
#
# docker-local.sh - Run cursor-sim locally in Docker with automated setup and health checks
#
# Usage:
#   ./tools/docker-local.sh                          # Default configuration (runtime, 90 days, medium)
#   DAYS=30 VELOCITY=high ./tools/docker-local.sh   # Custom configuration
#   DETACH=true ./tools/docker-local.sh             # Run in background (detached mode)
#   REBUILD=true ./tools/docker-local.sh            # Force rebuild image
#
# Environment Variables:
#   IMAGE_NAME      Docker image name (default: cursor-sim)
#   IMAGE_TAG       Docker image tag (default: latest)
#   CONTAINER_NAME  Container name (default: cursor-sim-local)
#   PORT            Host port to bind (default: 8080)
#   MODE            Operation mode: runtime or replay (default: runtime)
#   SEED_PATH       Path to seed file inside container (default: /app/seed.json)
#   DAYS            Days of history to generate (default: 90)
#   VELOCITY        Event rate: low, medium, high (default: medium)
#   DETACH          Run in background: true or false (default: false)
#   REBUILD         Force rebuild image: true or false (default: false)

set -euo pipefail

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration with defaults
IMAGE_NAME=${IMAGE_NAME:-cursor-sim}
IMAGE_TAG=${IMAGE_TAG:-latest}
CONTAINER_NAME=${CONTAINER_NAME:-cursor-sim-local}
PORT=${PORT:-8080}
MODE=${MODE:-runtime}
SEED_PATH=${SEED_PATH:-/app/seed.json}
DAYS=${DAYS:-90}
VELOCITY=${VELOCITY:-medium}
DETACH=${DETACH:-false}
REBUILD=${REBUILD:-false}

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Cleanup function to stop and remove container
cleanup() {
    if [ "$DETACH" != "true" ]; then
        log_info "Stopping container ${CONTAINER_NAME}..."
        docker stop "${CONTAINER_NAME}" 2>/dev/null || true
    fi
}

# Trap cleanup on exit signals (only in foreground mode)
if [ "$DETACH" != "true" ]; then
    trap cleanup EXIT INT TERM
fi

# Change to repository root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$( cd "${SCRIPT_DIR}/.." && pwd )"
cd "${REPO_ROOT}"

# Display configuration
log_info "Starting cursor-sim Docker container"
log_debug "Configuration:"
log_debug "  Image: ${IMAGE_NAME}:${IMAGE_TAG}"
log_debug "  Container: ${CONTAINER_NAME}"
log_debug "  Port: ${PORT}"
log_debug "  Mode: ${MODE}"
log_debug "  Seed: ${SEED_PATH}"
log_debug "  Days: ${DAYS}"
log_debug "  Velocity: ${VELOCITY}"
log_debug "  Detach: ${DETACH}"

# Check if image exists, build if not (or if REBUILD=true)
if [ "$REBUILD" = "true" ] || ! docker image inspect "${IMAGE_NAME}:${IMAGE_TAG}" &>/dev/null; then
    if [ "$REBUILD" = "true" ]; then
        log_info "Rebuilding Docker image ${IMAGE_NAME}:${IMAGE_TAG}..."
    else
        log_info "Image not found. Building Docker image ${IMAGE_NAME}:${IMAGE_TAG}..."
    fi

    BUILD_START=$(date +%s)
    docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" services/cursor-sim
    BUILD_END=$(date +%s)
    BUILD_TIME=$((BUILD_END - BUILD_START))

    log_info "Build completed in ${BUILD_TIME} seconds"

    # Show image size
    IMAGE_SIZE=$(docker images "${IMAGE_NAME}:${IMAGE_TAG}" --format "{{.Size}}")
    log_info "Image size: ${IMAGE_SIZE}"
else
    log_info "Using existing image ${IMAGE_NAME}:${IMAGE_TAG}"
fi

# Stop and remove existing container if running
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    log_warn "Container ${CONTAINER_NAME} already exists. Removing..."
    docker stop "${CONTAINER_NAME}" 2>/dev/null || true
    docker rm "${CONTAINER_NAME}" 2>/dev/null || true
fi

# Run container
log_info "Starting container ${CONTAINER_NAME}..."

DOCKER_RUN_ARGS=(
    --name "${CONTAINER_NAME}"
    -p "${PORT}:8080"
    -e "CURSOR_SIM_MODE=${MODE}"
    -e "CURSOR_SIM_SEED=${SEED_PATH}"
    -e "CURSOR_SIM_DAYS=${DAYS}"
    -e "CURSOR_SIM_VELOCITY=${VELOCITY}"
    -e "CURSOR_SIM_PORT=8080"
)

# Add --rm flag only for foreground mode
if [ "$DETACH" = "true" ]; then
    DOCKER_RUN_ARGS+=(-d)
else
    DOCKER_RUN_ARGS+=(--rm)
fi

docker run "${DOCKER_RUN_ARGS[@]}" "${IMAGE_NAME}:${IMAGE_TAG}"

# If detached, verify health
if [ "$DETACH" = "true" ]; then
    log_info "Container started in background. Waiting for service to be ready..."

    # Wait for container to start
    sleep 2

    # Check if container is still running
    if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_error "Container failed to start. Logs:"
        docker logs "${CONTAINER_NAME}"
        exit 1
    fi

    # Try health check (retry 3 times)
    HEALTH_OK=false
    for i in {1..3}; do
        log_debug "Health check attempt ${i}/3..."
        if curl -sf http://localhost:${PORT}/health >/dev/null 2>&1; then
            HEALTH_OK=true
            break
        fi
        sleep 2
    done

    if [ "$HEALTH_OK" = "true" ]; then
        log_info "Container is healthy!"
        echo ""
        log_info "Service running at: http://localhost:${PORT}"
        log_info "Health check: curl http://localhost:${PORT}/health"
        log_info "Example API call: curl -u cursor-sim-dev-key: http://localhost:${PORT}/teams/members"
        log_info "View logs: docker logs -f ${CONTAINER_NAME}"
        log_info "Stop container: docker stop ${CONTAINER_NAME}"
    else
        log_error "Health check failed after 3 attempts"
        log_error "Container logs:"
        docker logs "${CONTAINER_NAME}"
        exit 1
    fi
else
    # Foreground mode - logs will stream
    log_info "Container stopped."
fi
