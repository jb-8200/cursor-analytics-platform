# Design: cursor-sim Phase 2 (GitHub PR Simulation)

**Feature**: cursor-sim v2 Phase 2
**Status**: NOT_STARTED
**Estimated Hours**: 20-25

---

## Overview

Phase 2 extends cursor-sim with GitHub PR simulation, creating a complete SDLC data pipeline from commit to merge. This enables research into AI-assisted coding's impact on code review cycles, merge times, and code quality.

## Architecture Changes

### New Data Flow

```
Seed File ──► Commit Generator ──► PR Generator ──► Review Simulator
                    │                    │                │
                    ▼                    ▼                ▼
              ┌─────────────────────────────────────────────────┐
              │              In-Memory Storage                   │
              │  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
              │  │ Commits  │ │   PRs    │ │ ReviewComments   │ │
              │  └──────────┘ └──────────┘ └──────────────────┘ │
              └─────────────────────────────────────────────────┘
                                    │
                                    ▼
              ┌─────────────────────────────────────────────────┐
              │              HTTP Router                         │
              │  /repos/*   /repos/{o}/{r}/pulls/*              │
              └─────────────────────────────────────────────────┘
```

### New Packages

```
internal/
├── models/
│   ├── pr.go              # PR, ReviewComment types
│   └── repository.go      # Repository type
├── generator/
│   ├── pr_generator.go    # PR creation from commits
│   └── review_generator.go # Review comment generation
├── storage/
│   └── pr_store.go        # PR storage methods
└── api/
    └── github/            # GitHub-compatible handlers
        ├── repos.go
        ├── pulls.go
        └── commits.go
```

## Data Models

### PullRequest

```go
type PullRequest struct {
    Number      int       `json:"number"`
    Title       string    `json:"title"`
    Body        string    `json:"body"`
    State       string    `json:"state"` // open, closed, merged
    AuthorID    string    `json:"author_id"`
    AuthorEmail string    `json:"author_email"`
    RepoName    string    `json:"repo_name"`
    BaseBranch  string    `json:"base_branch"`
    HeadBranch  string    `json:"head_branch"`
    Reviewers   []string  `json:"reviewers"`
    Labels      []string  `json:"labels"`

    // Metrics
    Additions   int       `json:"additions"`
    Deletions   int       `json:"deletions"`
    ChangedFiles int      `json:"changed_files"`
    CommitCount int       `json:"commit_count"`

    // AI metrics (aggregated from commits)
    AIRatio     float64   `json:"ai_ratio"`
    TabLines    int       `json:"tab_lines"`
    ComposerLines int     `json:"composer_lines"`

    // Timestamps
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    MergedAt    *time.Time `json:"merged_at,omitempty"`
    ClosedAt    *time.Time `json:"closed_at,omitempty"`

    // Quality signals
    WasReverted bool      `json:"was_reverted"`
    IsBugFix    bool      `json:"is_bug_fix"`
}
```

### ReviewComment

```go
type ReviewComment struct {
    ID        int       `json:"id"`
    PRNumber  int       `json:"pr_number"`
    AuthorID  string    `json:"author_id"`
    Body      string    `json:"body"`
    Path      string    `json:"path,omitempty"`
    Line      int       `json:"line,omitempty"`
    State     string    `json:"state"` // pending, approved, changes_requested
    CreatedAt time.Time `json:"created_at"`
}
```

## Generation Algorithm

### PR Generation

1. **Cluster Commits**: Group commits by (developer, branch, time window)
2. **Create PR**: For each cluster:
   - Generate PR number (sequential)
   - Title from first commit or template
   - Body from commit messages
   - State based on lifecycle simulation
3. **Link Commits**: Update commit records with pr_number

### Review Simulation

1. **Assign Reviewers**: Select 1-3 from same team (not author)
2. **Generate Timeline**:
   - review_time = lognormal(seed.pr_lifecycle.review_time_hours)
   - iterations = min(poisson(seed.pr_lifecycle.iterations.mean), max)
3. **Generate Comments**:
   - count = poisson(reviewer.review_thoroughness * 3)
   - Body from seed.text_templates.review_comments

### Quality Outcomes

```go
func (g *Generator) assignQualityOutcome(pr *PullRequest) {
    // Revert probability based on AI ratio
    aiRatioCategory := categorizeAIRatio(pr.AIRatio) // high, medium, low
    revertRate := g.seed.Correlations.AIRatioRevertRate[aiRatioCategory]

    if g.rng.Float64() < revertRate {
        pr.WasReverted = true
    }

    // Bug fix detection
    pr.IsBugFix = strings.HasPrefix(pr.Title, "fix:") ||
                  strings.Contains(pr.Title, "bug")
}
```

## API Endpoints

### GitHub-Compatible Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/repos` | List all repositories |
| GET | `/repos/{owner}/{repo}` | Get repository details |
| GET | `/repos/{owner}/{repo}/pulls` | List PRs (supports state filter) |
| GET | `/repos/{owner}/{repo}/pulls/{number}` | Get PR details |
| GET | `/repos/{owner}/{repo}/pulls/{number}/commits` | Get PR commits |
| GET | `/repos/{owner}/{repo}/pulls/{number}/reviews` | Get PR reviews |

### Query Parameters

- `state`: open, closed, all (default: open)
- `sort`: created, updated, popularity (default: created)
- `direction`: asc, desc (default: desc)
- `per_page`: 1-100 (default: 30)
- `page`: Page number

## Seed Extensions

```json
{
  "pr_lifecycle": {
    "review_time_hours": {"mean": 24, "std_dev": 12},
    "iterations": {"mean": 2, "max": 5}
  },
  "text_templates": {
    "pr_titles": ["[{team}] {type}: {description}"],
    "pr_descriptions": ["## Summary\n{commits}\n\n## Testing\n{test_plan}"],
    "review_comments": ["Consider {suggestion}", "This could be simplified"]
  },
  "correlations": {
    "ai_ratio_revert_rate": {"high": 0.05, "medium": 0.08, "low": 0.12}
  }
}
```

## Testing Strategy

### Unit Tests
- PR generation from commit clusters
- Review comment generation
- Quality outcome assignment
- Seed validation for new fields

### Integration Tests
- Full pipeline: seed → commits → PRs → reviews
- API response schema validation against GitHub API

### E2E Tests
- Complete workflow with Insomnia/curl
- Pagination and filtering

## Decision Log

| Decision | Rationale |
|----------|-----------|
| GitHub API compatibility | Standard tooling support, familiar schema |
| In-memory PR storage | Consistent with Phase 1 architecture |
| Probabilistic reverts | Enables AI quality research |
| Separate review generator | Clean separation of concerns |
