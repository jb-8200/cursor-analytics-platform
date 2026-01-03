# User Story: Local Docker Deployment

**Feature ID**: P7-F01-local-docker-deploy
**Created**: January 3, 2026
**Status**: Planning (Ready to Start)

---

## Story

**As a** developer or researcher running cursor-sim locally
**I want** to containerize cursor-sim and run it in Docker with simple commands
**So that** I can ensure consistent environments, test deployments locally, and avoid "works on my machine" issues

## Background

Currently, cursor-sim runs directly on the host machine requiring Go installation, correct PATH setup, and manual dependency management. This creates several pain points:

- **Environment inconsistency**: Different Go versions, OS differences, and missing dependencies cause failures
- **Deployment testing**: No way to test production-like deployments before pushing to cloud
- **Onboarding friction**: New team members need to install Go, configure environments, and understand build steps
- **Reproducibility**: Research datasets should be generated in identical environments

Docker solves these by:
- Packaging the application with all dependencies in an immutable image
- Providing identical runtime environments across machines
- Enabling local testing of production configurations
- Simplifying onboarding to `docker run` instead of Go setup

This feature is the foundation for P7-F02 (GCP Cloud Run deployment), as we'll use the same Docker image locally and in production.

## Acceptance Criteria

### AC-1: Multi-Stage Dockerfile Creation

**Given** cursor-sim source code in `services/cursor-sim/`
**When** I run `docker build -t cursor-sim:latest services/cursor-sim`
**Then** the build completes successfully in < 2 minutes
**And** the final image size is < 50MB
**And** the image uses a distroless or minimal base (not full OS)
**And** the image runs as non-root user
**And** the build uses multi-stage process (builder + runtime)

Additional details:
- Builder stage: `golang:1.22-alpine` with `CGO_ENABLED=0`
- Runtime stage: `gcr.io/distroless/static:nonroot` or `cgr.dev/chainguard/static`
- Binary built with `-ldflags="-s -w"` for size optimization
- Layer caching optimized (dependencies before source code)

### AC-2: Local Container Execution (Runtime Mode)

**Given** a built cursor-sim Docker image
**When** I run the container with runtime mode configuration:
```bash
docker run -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  -e CURSOR_SIM_VELOCITY=medium \
  cursor-sim:latest
```
**Then** the container starts successfully within 5 seconds
**And** the health endpoint responds: `GET http://localhost:8080/health` returns 200
**And** API endpoints are accessible with Basic Auth (key: `cursor-sim-dev-key`)
**And** container logs show "Server listening on :8080"
**And** the container can be stopped gracefully with `Ctrl+C`

Edge cases:
- Missing required env vars should fail fast with clear error message
- Invalid seed file path should be caught at startup
- Port already in use should be reported clearly

### AC-3: Environment Variable Configuration

**Given** a running cursor-sim container
**When** I configure different environment variables:
- `CURSOR_SIM_MODE=runtime` or `CURSOR_SIM_MODE=replay`
- `CURSOR_SIM_DAYS=30` to `CURSOR_SIM_DAYS=180`
- `CURSOR_SIM_VELOCITY=low|medium|high`
- `CURSOR_SIM_PORT=8080` or custom port
**Then** the application uses these configurations without rebuilding the image
**And** invalid values are rejected at startup with descriptive errors
**And** missing required variables (MODE, SEED for runtime) cause immediate exit

### AC-4: .dockerignore Configuration

**Given** the `.dockerignore` file in `services/cursor-sim/`
**When** the Docker build copies files into the image
**Then** the following are excluded from the build context:
- `.git/` directory
- `bin/` directory
- `coverage.out`, `coverage.html`
- `*.log` files
- `.DS_Store` files
- `**/*.test` binaries
**And** the build context size is < 5MB
**And** build time is optimized by excluding unnecessary files

### AC-5: Local Testing Script

**Given** the `tools/docker-local.sh` script
**When** I run `./tools/docker-local.sh`
**Then** the script builds the image, runs the container, and verifies health
**And** I can pass environment variables: `DAYS=30 ./tools/docker-local.sh`
**And** the script displays the service URL and example curl commands
**And** logs are tailed automatically (with option to disable)
**And** `Ctrl+C` stops the container gracefully

Script features:
- Default values for all environment variables
- Health check verification before showing "Ready"
- Color-coded output for readability
- Error handling with clear messages

### AC-6: Volume Mounting for Development

**Given** I'm developing locally and want to use custom seed files
**When** I mount a volume with custom seed:
```bash
docker run -p 8080:8080 \
  -v $(pwd)/testdata/custom_seed.json:/app/seed.json \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  cursor-sim:latest
```
**Then** the container uses my custom seed file
**And** changes to the mounted file are reflected on container restart
**And** file permissions work correctly (container can read mounted files)

### AC-7: Container Health and Logging

**Given** a running cursor-sim container
**When** I check container health and logs
**Then** `docker logs <container-id>` shows structured application logs
**And** startup logs include configuration summary (mode, days, velocity)
**And** API request logs include timestamp, method, path, status code
**And** errors are logged with stack traces
**And** logs use JSON format for easy parsing (optional enhancement)

## Out of Scope

- **Docker Compose multi-service setup**: Single container only, no database/dependencies
- **Docker Swarm or Kubernetes**: Only standalone Docker, orchestration in future phases
- **Image registry push (local)**: Only build and run locally, no Docker Hub/registry push
- **Production hardening**: Security scanning, image signing, vulnerability checks (handled in P7-F02)
- **Replay mode corpus files**: Focus on runtime mode, replay corpus handling is future work
- **Performance optimization**: Beyond basic multi-stage build, no custom base image builds

## Dependencies

- **Go 1.22+**: Required in builder stage
- **Docker Desktop or Docker Engine**: Installed and running on developer machine
- **cursor-sim application code**: Services/cursor-sim builds successfully with `go build`
- **Valid seed file**: Either baked into image or mounted as volume
- **docs/cursor-sim-cloud-run.md**: Reference documentation exists

## Success Metrics

- **Image build time**: < 2 minutes on standard developer machine
- **Image size**: < 50MB final runtime image
- **Container startup time**: < 5 seconds from `docker run` to health check passing
- **Developer onboarding**: New developer can run cursor-sim in < 5 minutes (install Docker, run script)
- **Build cache efficiency**: Rebuilds with no code changes complete in < 30 seconds
- **Health check success rate**: 100% for valid configurations
- **Documentation completeness**: All environment variables and commands documented

## Related Documents

- `services/cursor-sim/SPEC.md` - cursor-sim technical specification
- `docs/cursor-sim-cloud-run.md` - Deployment documentation (Containerization Design section)
- `.work-items/P7-F01-local-docker-deploy/design.md` - Technical design for Docker implementation
- `.work-items/P7-F01-local-docker-deploy/task.md` - Task breakdown and implementation steps
- `.claude/agents/cursor-sim-infra-dev.md` - Infrastructure agent for implementation
