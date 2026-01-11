# P3-F05: OpenAPI Web UI API Catalog and Documentation

## User Story

**As a** developer or API consumer
**I want** an interactive web-based API documentation interface in the cursor-sim Docker container
**So that** I can explore, understand, and test the API endpoints without external tools

## Acceptance Criteria

### AC1: Authentication
**Given** a user visits `/docs`
**When** they are not authenticated
**Then** they are redirected to `/docs/login`

**Given** a user is on the login page
**When** they enter username `dox` and password `dox-a3`
**Then** they are authenticated and redirected to `/docs`

**Given** a user has an active session
**When** they visit `/docs/logout`
**Then** their session is cleared and they are redirected to login

### AC2: API Documentation Display
**Given** an authenticated user visits `/docs`
**When** the page loads
**Then** they see an interactive API reference powered by Scalar

**Given** an authenticated user is on the docs page
**When** they select a spec from the dropdown
**Then** they can view either Cursor API or GitHub Simulation API documentation

### AC3: Interactive Testing
**Given** an authenticated user is viewing an endpoint
**When** they use the "Try It" feature with a valid API key
**Then** they can execute requests and see responses

### AC4: Docker Integration
**Given** the cursor-sim Docker image is built
**When** the container starts
**Then** the `/docs` endpoint is available on port 8080

## Out of Scope

- User registration or multiple user accounts
- Persistent session storage (in-memory is sufficient)
- Custom branding or theming beyond Scalar defaults
- API key management through the docs UI

## Dependencies

- OpenAPI specs: `specs/openapi/cursor-api.yaml`, `specs/openapi/github-sim-api.yaml`
- Scalar CDN for documentation UI
- Existing cursor-sim HTTP server infrastructure

## References

- [Harvey AI API Docs](https://developers.harvey.ai/api-reference/usage/get-usage-history)
- [GitHub REST API](https://docs.github.com/en/rest/pulls/pulls)
- [Microsoft Graph API](https://learn.microsoft.com/en-us/graph/api/)
- [Scalar Documentation](https://github.com/scalar/scalar)
