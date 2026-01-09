# GitHub API Routing Fix (2026-01-09)

## Problem Summary

Multiple GitHub-style REST API endpoints in cursor-sim were returning 404 errors when accessed through the Insomnia collection. While some endpoints worked correctly, many repo-specific endpoints failed.

### Affected Endpoints

| Endpoint | Status Before | Status After |
|----------|--------------|--------------|
| `GET /repos` | 200 (empty []) | 200 (data) |
| `GET /repos/{owner}/{repo}` | 404 | 200 |
| `GET /repos/{owner}/{repo}/commits` | 404 | 200 |
| `GET /repos/{owner}/{repo}/pulls` | 404 | 200 |
| `GET /repos/{owner}/{repo}/pulls/{number}` | 200 | 200 |
| `GET /repos/{owner}/{repo}/pulls/{number}/commits` | 404 | 200 |
| `GET /repos/{owner}/{repo}/pulls/{number}/files` | 404 | 200 |
| `GET /repos/{owner}/{repo}/pulls/{number}/reviews` | 200 (0 reviews) | 200 (data) |
| `GET /repos/{owner}/{repo}/commits/{sha}` | 200 | 200 |

## Root Causes

### 1. Router Path Segment Counting (router.go)

The `countPathSegments` function in `router.go` correctly excluded leading empty strings from path splits, but all the case statements had segment counts off by 1.

**File:** `internal/api/github/router.go`

**Before:**
```go
case !strings.Contains(path, "/pulls") && !strings.Contains(path, "/commits") && countPathSegments(path) == 4:
    return GetRepository(store)
case strings.HasSuffix(path, "/pulls") && countPathSegments(path) == 5:
    return ListPulls(store)
// ... etc
```

**After:**
```go
case !strings.Contains(path, "/pulls") && !strings.Contains(path, "/commits") && !strings.Contains(path, "/analysis") && countPathSegments(path) == 3:
    return GetRepository(store)
case strings.HasSuffix(path, "/pulls") && countPathSegments(path) == 4:
    return ListPulls(store)
// ... etc
```

**Segment Count Corrections:**
- `GET /repos/{owner}/{repo}`: 4 → 3
- `GET /repos/{owner}/{repo}/pulls`: 5 → 4
- `GET /repos/{owner}/{repo}/pulls/{number}`: 6 → 5
- `GET /repos/{owner}/{repo}/pulls/{number}/reviews`: 7 → 6
- `GET /repos/{owner}/{repo}/pulls/{number}/commits`: 7 → 6
- `GET /repos/{owner}/{repo}/pulls/{number}/files`: 7 → 6
- `GET /repos/{owner}/{repo}/commits`: 5 → 4
- `GET /repos/{owner}/{repo}/commits/{sha}`: 6 → 5

### 2. Missing Data Generation (main.go)

The runtime mode startup only generated commits, not PRs, reviews, or issues.

**File:** `cmd/simulator/main.go`

**Added:**
```go
// Generate PRs from commits
prGen := generator.NewPRGeneratorWithSeed(seedData, store, time.Now().UnixNano())
err := prGen.GeneratePRsFromCommits(startDate, endDate)

// Generate reviews for PRs
reviewGen := generator.NewReviewGenerator(seedData, rand.New(rand.NewSource(time.Now().UnixNano())))
for _, repoName := range repos {
    prs := store.GetPRsByRepo(repoName)
    for _, pr := range prs {
        reviews := reviewGen.GenerateReviewsForPR(pr)
        for _, review := range reviews {
            store.StoreReview(review)
        }
    }
}

// Generate issues for PRs
issueGen := generator.NewIssueGeneratorWithStore(seedData, store, time.Now().UnixNano())
for _, repoName := range repos {
    prs := store.GetPRsByRepo(repoName)
    issueGen.GenerateAndStoreIssuesForPRs(prs, repoName)
}
```

### 3. Review PRID Assignment (review_generator.go)

Reviews were being stored with `PRID: pr.Number` but lookups used `pr.ID`, causing a mismatch.

**File:** `internal/generator/review_generator.go`

**Before:**
```go
review := models.Review{
    PRID: pr.Number,  // Wrong: used PR number
    // ...
}
```

**After:**
```go
review := models.Review{
    PRID: pr.ID,  // Correct: uses PR ID
    // ...
}
```

### 4. PR ID Auto-Generation (memory.go)

The `AddPR` function didn't assign unique IDs to PRs, leaving them all at 0.

**File:** `internal/storage/memory.go`

**Added:**
```go
// Added to struct:
nextPRID int // Auto-incrementing PR ID counter

// In NewMemoryStore:
nextPRID: 1, // Start PR IDs at 1

// In AddPR:
if pr.ID == 0 {
    pr.ID = m.nextPRID
    m.nextPRID++
}
```

### 5. ListPullReviews Handler (pulls.go)

The handler was calling `GetReviewComments` (which returns `[]ReviewComment`) instead of `GetReviewsByRepoPR` (which returns `[]Review`).

**File:** `internal/api/github/pulls.go`

**Before:**
```go
func ListPullReviews(store storage.Store) http.Handler {
    // ...
    reviews := store.GetReviewComments(repoName, prNumber)  // Wrong type
}
```

**After:**
```go
func ListPullReviews(store storage.Store) http.Handler {
    // ...
    reviews, err := store.GetReviewsByRepoPR(repoName, prNumber)  // Correct type
}
```

## Testing

### Insomnia Collection

Updated environment variables in `docs/insomnia/cursor-sim_Insomnia_2026-01-04.yaml`:
- `repoOwner`: `acme-corp` (matches seed data)
- `repoName`: `payment-service` (matches seed data)

### E2E Tests

Added comprehensive E2E tests in `test/e2e/insomnia_endpoints_test.go`:
- Tests all endpoints from the Insomnia collection
- Verifies data is returned (not empty)
- Validates response structure matches API contract
- Tests dynamic endpoints with actual repo/PR data

Run tests with:
```bash
cd services/cursor-sim
go test ./test/e2e/... -run TestInsomnia -v
```

## Verification

After applying all fixes, the final API test showed:

| Category | Results |
|----------|---------|
| Health | ok |
| Team Members | 2 members |
| AI Code Commits | 100+ commits |
| PRs | 100+ total |
| Reviews | 100+ total |
| Issues | 30+ total |
| List Repos | 2 repos |
| PR Reviews | Data returned |
| All 45+ endpoints | 200 OK |

## Files Modified

1. `internal/api/github/router.go` - Fixed segment counts
2. `cmd/simulator/main.go` - Added PR/review/issue generation
3. `internal/generator/review_generator.go` - Fixed PRID assignment
4. `internal/storage/memory.go` - Added PR ID auto-generation
5. `internal/api/github/pulls.go` - Fixed ListPullReviews handler
6. `docs/insomnia/cursor-sim_Insomnia_2026-01-04.yaml` - Updated environment variables
7. `test/e2e/insomnia_endpoints_test.go` - Added comprehensive E2E tests (new file)
