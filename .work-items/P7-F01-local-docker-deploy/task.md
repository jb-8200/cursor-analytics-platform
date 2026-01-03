# Task Breakdown: Local Docker Deployment

**Feature ID**: P7-F01-local-docker-deploy
**Created**: January 3, 2026
**Epic**: P7 - Deployment Infrastructure
**Estimated Time**: 4.0 hours

---

## Task List

| Task ID | Task Name | Status | Est. Time | Actual Time | Assignee |
|---------|-----------|--------|-----------|-------------|----------|
| INFRA-01 | Create Dockerfile and .dockerignore | DONE | 1.0h | 0.5h | cursor-sim-infra-dev |
| INFRA-02 | Build and verify Docker image locally | DONE | 0.5h | 0.5h | cursor-sim-infra-dev |
| INFRA-03 | Create local testing script (docker-local.sh) | DONE | 1.0h | 0.5h | cursor-sim-infra-dev |
| INFRA-04 | Test environment variable configurations | DONE | 0.5h | 0.25h | cursor-sim-infra-dev |
| INFRA-05 | Document Docker usage and troubleshooting | DONE | 0.5h | 0.5h | cursor-sim-infra-dev |
| INFRA-06 | Integration test and commit | DONE | 0.5h | 0.25h | cursor-sim-infra-dev |

**Total Estimated**: 4.0 hours
**Total Actual**: 2.5 hours

---

## Progress Tracker

```
Feature: Local Docker Deployment (P7-F01)
[████████████████████] 100% (6/6 tasks) ✅ COMPLETE

Tasks:
[x] INFRA-01: Create Dockerfile and .dockerignore
[x] INFRA-02: Build and verify Docker image locally
[x] INFRA-03: Create local testing script (docker-local.sh)
[x] INFRA-04: Test environment variable configurations
[x] INFRA-05: Document Docker usage and troubleshooting
[x] INFRA-06: Integration test and commit
```

---

## Detailed Task Breakdown

### INFRA-01: Create Dockerfile and .dockerignore (1.0h)

**Status**: TODO
**Prerequisites**: None
**Agent**: cursor-sim-infra-dev

**Objective**: Create production-ready multi-stage Dockerfile and .dockerignore for cursor-sim

**Acceptance Criteria**:
- ✅ Dockerfile uses multi-stage build (golang:1.22-alpine + distroless/static:nonroot)
- ✅ Builder stage sets CGO_ENABLED=0 and GOTOOLCHAIN=auto
- ✅ Binary built with -ldflags="-s -w" for size optimization
- ✅ Runtime stage uses distroless non-root user
- ✅ EXPOSE 8080 and WORKDIR /app configured
- ✅ ENTRYPOINT points to /app/cursor-sim
- ✅ .dockerignore excludes .git, bin/, coverage files, logs, IDE files
- ✅ Build context size < 5MB

**Steps**:
1. Read `docs/cursor-sim-cloud-run.md` Containerization Design section
2. Review `internal/config/config.go` for environment variable names
3. Create `services/cursor-sim/Dockerfile` with multi-stage build
4. Create `services/cursor-sim/.dockerignore` with exclusions
5. Verify syntax with `docker build --dry-run` (if available)

**TDD Approach**:
- **RED**: Attempt build without Dockerfile → fails
- **GREEN**: Create minimal Dockerfile → builds successfully
- **REFACTOR**: Add optimizations (layer caching, size reduction)

**Files Created**:
- `services/cursor-sim/Dockerfile`
- `services/cursor-sim/.dockerignore`

**Validation**:
```bash
cd services/cursor-sim
docker build -t cursor-sim:test .
# Should complete without errors
```

---

### INFRA-02: Build and Verify Docker Image Locally (0.5h)

**Status**: TODO
**Prerequisites**: INFRA-01
**Agent**: cursor-sim-infra-dev

**Objective**: Build Docker image and verify all acceptance criteria

**Acceptance Criteria**:
- ✅ Image builds in < 2 minutes (cold, no cache)
- ✅ Image builds in < 30 seconds (warm, cached deps)
- ✅ Final image size < 50MB
- ✅ Image runs as non-root user (UID 65532)
- ✅ Image contains only necessary files (binary, distroless base)
- ✅ No vulnerabilities detected (optional scout scan)

**Steps**:
1. Build image: `docker build -t cursor-sim:latest services/cursor-sim`
2. Check image size: `docker images cursor-sim:latest --format "{{.Size}}"`
3. Inspect layers: `docker history cursor-sim:latest`
4. Verify user: `docker run --rm cursor-sim:latest id` (should show nonroot)
5. Optional: Run security scan: `docker scout cves cursor-sim:latest`

**TDD Approach**:
- **RED**: Image doesn't exist → docker run fails
- **GREEN**: Build image → docker run starts successfully
- **REFACTOR**: Optimize image size, verify security

**Validation**:
```bash
# Build
docker build -t cursor-sim:latest services/cursor-sim

# Verify size
SIZE=$(docker images cursor-sim:latest --format "{{.Size}}")
echo "Image size: $SIZE (should be < 50MB)"

# Verify multi-stage worked
docker history cursor-sim:latest | grep golang
# Should NOT show golang in final layers
```

---

### INFRA-03: Create Local Testing Script (docker-local.sh) (1.0h)

**Status**: TODO
**Prerequisites**: INFRA-02
**Agent**: cursor-sim-infra-dev

**Objective**: Create automated script for local Docker testing with health verification

**Acceptance Criteria**:
- ✅ Script builds image if not exists
- ✅ Script runs container with sensible defaults
- ✅ Environment variables can be overridden (DAYS, VELOCITY, MODE, etc.)
- ✅ Health check verification before showing "Ready"
- ✅ Color-coded output (green=success, red=error, yellow=warning)
- ✅ Container cleanup on exit (trap INT/TERM)
- ✅ Detached mode option (DETACH=true)
- ✅ Usage examples in script header comments

**Steps**:
1. Create `tools/docker-local.sh` with bash shebang
2. Add `set -euo pipefail` for strict error handling
3. Define environment variable defaults (MODE=runtime, DAYS=90, etc.)
4. Implement build logic (check if image exists, build if not)
5. Implement run logic with `docker run --rm`
6. Add health check verification loop (retry 3 times, 1s apart)
7. Add cleanup trap for graceful shutdown
8. Test with various environment variable combinations

**TDD Approach**:
- **RED**: No script exists → manual docker commands required
- **GREEN**: Create basic script → container starts
- **REFACTOR**: Add health checks, error handling, colors

**Files Created**:
- `tools/docker-local.sh`

**Validation**:
```bash
# Default run
./tools/docker-local.sh
# Should show: "Container is healthy!" and service URL

# Custom config
DAYS=30 VELOCITY=high ./tools/docker-local.sh

# Detached mode
DETACH=true ./tools/docker-local.sh
curl http://localhost:8080/health
```

---

### INFRA-04: Test Environment Variable Configurations (0.5h)

**Status**: TODO
**Prerequisites**: INFRA-03
**Agent**: cursor-sim-infra-dev

**Objective**: Verify all environment variable combinations work correctly

**Acceptance Criteria**:
- ✅ MODE=runtime with valid SEED starts successfully
- ✅ MODE=replay fails gracefully (not yet implemented)
- ✅ Missing MODE fails with clear error message
- ✅ Missing SEED (in runtime mode) fails with clear error
- ✅ DAYS values (30, 90, 180) all work correctly
- ✅ VELOCITY values (low, medium, high) all work correctly
- ✅ Custom PORT value changes server port
- ✅ Invalid values are rejected at startup

**Steps**:
1. Test default configuration (runtime, 90 days, medium)
2. Test each DAYS value: 30, 60, 90, 180
3. Test each VELOCITY value: low, medium, high
4. Test missing MODE (should fail)
5. Test missing SEED in runtime mode (should fail)
6. Test invalid VELOCITY (should fail or use default)
7. Test custom PORT (verify server listens on new port)
8. Document any edge cases or unexpected behavior

**TDD Approach**:
- **RED**: Invalid config starts container → no validation
- **GREEN**: Add validation in config.go → rejects invalid inputs
- **REFACTOR**: Improve error messages for clarity

**Validation**:
```bash
# Valid configurations
MODE=runtime SEED=/app/seed.json DAYS=30 ./tools/docker-local.sh
MODE=runtime SEED=/app/seed.json VELOCITY=high ./tools/docker-local.sh

# Invalid configurations (should fail gracefully)
MODE=invalid ./tools/docker-local.sh  # Should error
MODE=runtime SEED=/nonexistent ./tools/docker-local.sh  # Should error
```

---

### INFRA-05: Document Docker Usage and Troubleshooting (0.5h)

**Status**: TODO
**Prerequisites**: INFRA-04
**Agent**: cursor-sim-infra-dev

**Objective**: Update docs/cursor-sim-cloud-run.md with local Docker sections

**Acceptance Criteria**:
- ✅ "Manual Docker Deploy (local)" section updated with current commands
- ✅ Environment variable reference table accurate
- ✅ Common issues and troubleshooting section added
- ✅ Example commands for volume mounting documented
- ✅ Health check verification steps documented
- ✅ Image size and build time benchmarks documented

**Steps**:
1. Read existing `docs/cursor-sim-cloud-run.md`
2. Update "Manual Docker Deploy" section with tested commands
3. Add troubleshooting section:
   - "Container starts but health check fails" → check env vars
   - "Build fails" → check .dockerignore, verify go.mod
   - "Image too large" → verify multi-stage build
4. Document volume mounting for custom seed files
5. Add performance benchmarks (image size, build time, startup time)

**Files Modified**:
- `docs/cursor-sim-cloud-run.md`

**Validation**:
- Documentation includes all commands from INFRA-03 testing
- New developer can follow docs and run container successfully

---

### INFRA-06: Integration Test and Commit (0.5h)

**Status**: TODO
**Prerequisites**: INFRA-05
**Agent**: cursor-sim-infra-dev

**Objective**: Run full integration test and commit all changes following SDD workflow

**Acceptance Criteria**:
- ✅ Clean build from scratch completes successfully
- ✅ Health check passes
- ✅ API endpoints return valid JSON
- ✅ Container stops gracefully
- ✅ All files committed with descriptive message
- ✅ DEVELOPMENT.md updated with P7-F01 progress
- ✅ Dependency reflection check performed
- ✅ SPEC sync check performed (no SPEC.md changes needed)

**Steps**:
1. Clean Docker environment: `docker system prune -a`
2. Run full build: `docker build -t cursor-sim:latest services/cursor-sim`
3. Run container: `./tools/docker-local.sh`
4. Verify health: `curl http://localhost:8080/health`
5. Test API endpoint: `curl -u cursor-sim-dev-key: http://localhost:8080/api/team/members`
6. Stop container: `docker stop cursor-sim-local`
7. Run dependency-reflection check
8. Run spec-sync-check (should show no SPEC.md update needed)
9. Stage all files: `git add services/cursor-sim/Dockerfile services/cursor-sim/.dockerignore tools/docker-local.sh docs/cursor-sim-cloud-run.md`
10. Commit with message following format
11. Update `.claude/DEVELOPMENT.md` with P7-F01 completion

**Files Committed**:
- `services/cursor-sim/Dockerfile`
- `services/cursor-sim/.dockerignore`
- `tools/docker-local.sh`
- `docs/cursor-sim-cloud-run.md` (updated)
- `.claude/DEVELOPMENT.md` (updated)

**Commit Message**:
```
feat(infra): implement local Docker deployment (P7-F01)

Add Docker containerization for cursor-sim:
- Multi-stage Dockerfile (golang:1.22-alpine + distroless)
- .dockerignore for optimal build context
- tools/docker-local.sh for automated testing
- Documentation updates

Image size: ~22MB, startup time: ~2s, health check passing.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

**Validation**:
```bash
# Full integration test
git status  # Should show clean working tree
docker images cursor-sim:latest  # Should exist
docker ps -a | grep cursor-sim  # Should show no running containers
```

---

## Dependencies

### External Dependencies
- Docker Desktop or Docker Engine (20+)
- Go 1.22+ (for local development, not required for Docker)
- curl (for health checks)
- bash (for scripts)

### Internal Dependencies
- cursor-sim application builds successfully (`go build ./cmd/simulator`)
- Valid seed file exists (e.g., `testdata/valid_seed.json`)
- `internal/config/config.go` reads environment variables correctly

### Blocking Dependencies
- None (first feature in P7 epic)

---

## Testing Strategy

### Unit Tests (Pre-Docker)
```bash
cd services/cursor-sim
go test ./...
go test -cover ./...
```

### Docker Build Tests
```bash
# Clean build
docker build --no-cache -t cursor-sim:test services/cursor-sim

# Cached build
docker build -t cursor-sim:test services/cursor-sim

# Verify image
docker images cursor-sim:test
docker history cursor-sim:test
docker inspect cursor-sim:test
```

### Docker Runtime Tests
```bash
# Start container
docker run -d --name test-cursor-sim -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  -e CURSOR_SIM_DAYS=90 \
  cursor-sim:test

# Health check
curl http://localhost:8080/health

# API test
curl -u cursor-sim-dev-key: http://localhost:8080/api/team/members

# Logs
docker logs test-cursor-sim

# Cleanup
docker stop test-cursor-sim
docker rm test-cursor-sim
```

### Script Tests
```bash
# Default run
./tools/docker-local.sh

# Custom config
DAYS=30 VELOCITY=high ./tools/docker-local.sh

# Detached mode
DETACH=true ./tools/docker-local.sh
```

---

## Rollback Plan

If issues arise:

1. **Dockerfile build fails**:
   - Revert Dockerfile changes
   - Use `go build` directly as fallback

2. **Image too large (> 50MB)**:
   - Check multi-stage build is configured correctly
   - Verify .dockerignore excludes unnecessary files
   - Consider alternative runtime base (alpine instead of distroless)

3. **Container fails to start**:
   - Check environment variables are set correctly
   - Verify seed file path is accessible in container
   - Review container logs: `docker logs <container>`

4. **Health check fails**:
   - Verify port mapping: `-p 8080:8080`
   - Check firewall/network settings
   - Confirm application actually starts (check logs)

---

## Success Metrics

- **Image build time**: < 2 minutes (cold), < 30 seconds (warm)
- **Image size**: < 50MB
- **Container startup time**: < 5 seconds
- **Health check response**: < 200ms
- **Developer onboarding time**: < 5 minutes to run cursor-sim in Docker
- **Build cache hit rate**: > 90% for subsequent builds

---

## Next Steps

After P7-F01 completion:
- **P7-F02**: GCP Cloud Run Deployment (uses same Docker image)
- **P7-F03**: CI/CD Pipeline (automated builds and deployments)
- **P8**: Multi-environment deployments (dev, staging, prod)
