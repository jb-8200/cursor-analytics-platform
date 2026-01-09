# Design: GitHub Simulation for cursor-sim

**Feature ID**: P2-F01-github-simulation
**Status**: Not Started
**Created**: January 8, 2026

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────┐
│                    Generation Flow                       │
│                                                          │
│  ┌─────────────┐      ┌─────────────┐                   │
│  │   Commits   │  →   │     PRs     │  → Merged         │
│  │  (P1 Done)  │      │  Generator  │                   │
│  └─────────────┘      └──────┬──────┘                   │
│                              │                           │
│                              ↓                           │
│                       ┌─────────────┐                    │
│                       │   Reviews   │                    │
│                       │  Generator  │                    │
│                       └─────────────┘                    │
│                              ↓                           │
│                       ┌─────────────┐                    │
│                       │   Issues    │                    │
│                       │  Generator  │ ← Linked to PRs    │
│                       └─────────────┘                    │
└──────────────────────────────────────────────────────────┘
```

---

## Data Models

### PullRequest

```go
type PullRequest struct {
    ID          int       `json:"id"`
    Number      int       `json:"number"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Author      string    `json:"author"`        // Developer email
    Branch      string    `json:"branch"`        // feature/*, bugfix/*
    BaseBranch  string    `json:"base_branch"`   // main, develop
    Status      PRStatus  `json:"status"`        // open, merged, closed
    CreatedAt   time.Time `json:"created_at"`
    MergedAt    *time.Time `json:"merged_at,omitempty"`
    ClosedAt    *time.Time `json:"closed_at,omitempty"`

    // Relationships
    CommitIDs   []string  `json:"commit_ids"`    // Linked commits
    ReviewIDs   []int     `json:"review_ids"`    // Associated reviews
    IssueIDs    []int     `json:"issue_ids"`     // Closes issues

    // Metrics
    FilesChanged int      `json:"files_changed"`
    Additions    int      `json:"additions"`
    Deletions    int      `json:"deletions"`
}

type PRStatus string
const (
    PRStatusOpen   PRStatus = "open"
    PRStatusMerged PRStatus = "merged"
    PRStatusClosed PRStatus = "closed"
)
```

### Review

```go
type Review struct {
    ID         int         `json:"id"`
    PRID       int         `json:"pr_id"`
    Reviewer   string      `json:"reviewer"`      // Developer email
    State      ReviewState `json:"state"`
    SubmittedAt time.Time  `json:"submitted_at"`
    Body       string      `json:"body"`

    // Review comments
    Comments   []ReviewComment `json:"comments"`
}

type ReviewState string
const (
    ReviewStateApproved        ReviewState = "approved"
    ReviewStateChangesRequested ReviewState = "changes_requested"
    ReviewStateCommented       ReviewState = "commented"
)

type ReviewComment struct {
    ID        int       `json:"id"`
    Body      string    `json:"body"`
    Path      string    `json:"path"`      // File path
    Line      int       `json:"line"`      // Line number
    CreatedAt time.Time `json:"created_at"`
}
```

### Issue

```go
type Issue struct {
    ID          int         `json:"id"`
    Number      int         `json:"number"`
    Title       string      `json:"title"`
    Description string      `json:"description"`
    Creator     string      `json:"creator"`
    State       IssueState  `json:"state"`
    Labels      []string    `json:"labels"`
    CreatedAt   time.Time   `json:"created_at"`
    ClosedAt    *time.Time  `json:"closed_at,omitempty"`

    // Relationships
    ClosedByPRID *int       `json:"closed_by_pr_id,omitempty"`
    Assignee     *string    `json:"assignee,omitempty"`
}

type IssueState string
const (
    IssueStateOpen       IssueState = "open"
    IssueStateInProgress IssueState = "in_progress"
    IssueStateClosed     IssueState = "closed"
)
```

---

## Generation Algorithm

### 1. PR Generation (from Commits)

**Strategy**: Group commits into logical PR units

```go
func GeneratePRs(commits []Commit, config PRConfig) []PullRequest {
    // Group commits by author and temporal proximity
    // - Commits by same author within 24-48 hours → Same PR
    // - 3-10 commits per PR (Poisson distribution, λ=5)
    // - 60% of commits go into PRs, 40% direct to main

    prs := []PullRequest{}

    for _, commitGroup := range groupCommits(commits) {
        pr := PullRequest{
            Number:      nextPRNumber(),
            Title:       generatePRTitle(commitGroup),
            Author:      commitGroup[0].Author,
            Branch:      generateBranchName(),
            CreatedAt:   commitGroup[0].Timestamp.Add(-1 * time.Hour),
            CommitIDs:   extractCommitIDs(commitGroup),
            Status:      choosePRStatus(),  // 85% merged, 10% closed, 5% open
        }

        if pr.Status == PRStatusMerged {
            // Merge happens 1-7 days after creation
            pr.MergedAt = randomTime(pr.CreatedAt, pr.CreatedAt.Add(7*24*time.Hour))
        }

        prs = append(prs, pr)
    }

    return prs
}
```

**Distributions**:
- PR size: Poisson(λ=5) commits per PR
- Merge rate: 85% merged, 10% closed without merge, 5% open
- Time to merge: Exponential(λ=3 days)

### 2. Review Generation

**Strategy**: Assign reviewers based on team structure

```go
func GenerateReviews(pr PullRequest, team []Developer) []Review {
    // Reviewer count: 1-3 reviewers per PR
    // - 50% get 1 reviewer
    // - 35% get 2 reviewers
    // - 15% get 3 reviewers

    reviewerCount := sampleReviewerCount()
    reviewers := selectReviewers(team, pr.Author, reviewerCount)

    reviews := []Review{}
    for _, reviewer := range reviewers {
        review := Review{
            PRID:        pr.ID,
            Reviewer:    reviewer.Email,
            State:       chooseReviewState(),  // 70% approved, 20% changes, 10% commented
            SubmittedAt: randomTimeBetween(pr.CreatedAt, pr.MergedAt),
        }

        // Generate review comments (0-5 per review)
        review.Comments = generateReviewComments(pr)

        reviews = append(reviews, review)
    }

    return reviews
}
```

**Constraints**:
- Reviewer ≠ PR author
- Review timestamp: PR creation < review < PR merge
- First review within 24-48 hours (80% of PRs)

### 3. Issue Generation

**Strategy**: Create issues that are resolved by PRs

```go
func GenerateIssues(prs []PullRequest, config IssueConfig) []Issue {
    // 40% of PRs close an issue
    // 10% of issues remain open

    issues := []Issue{}

    for _, pr := range prs {
        if rand.Float64() < 0.4 && pr.Status == PRStatusMerged {
            issue := Issue{
                Number:       nextIssueNumber(),
                Title:        generateIssueTitle(pr),
                Creator:      selectRandomDeveloper(),
                State:        IssueStateClosed,
                Labels:       selectLabels(),  // bug, feature, enhancement
                CreatedAt:    pr.CreatedAt.Add(-randomDuration(24*time.Hour, 7*24*time.Hour)),
                ClosedAt:     &pr.MergedAt,
                ClosedByPRID: &pr.ID,
            }
            issues = append(issues, issue)
        }
    }

    // Generate 10% open issues
    openIssues := generateOpenIssues(config)
    issues = append(issues, openIssues...)

    return issues
}
```

---

## API Endpoints

### 1. `/analytics/github/prs`

**Method**: GET
**Query Parameters**:
- `status`: open, merged, closed
- `author`: Filter by PR author
- `start_date`, `end_date`: Date range
- `page`, `page_size`: Pagination

**Response**:
```json
{
  "data": [
    {
      "id": 1,
      "number": 42,
      "title": "Add user authentication",
      "author": "alice@example.com",
      "status": "merged",
      "created_at": "2024-01-15T10:00:00Z",
      "merged_at": "2024-01-18T14:30:00Z",
      "commits": 5,
      "reviews": 2,
      "files_changed": 8
    }
  ],
  "pagination": {...}
}
```

### 2. `/analytics/github/reviews`

**Response**:
```json
{
  "data": [
    {
      "pr_number": 42,
      "reviewer": "bob@example.com",
      "state": "approved",
      "submitted_at": "2024-01-17T09:00:00Z",
      "comments_count": 3
    }
  ]
}
```

### 3. `/analytics/github/issues`

**Response**:
```json
{
  "data": [
    {
      "number": 10,
      "title": "Fix login bug",
      "state": "closed",
      "labels": ["bug", "priority-high"],
      "closed_by_pr": 42
    }
  ]
}
```

### 4. `/analytics/github/pr-cycle-time`

**Response**:
```json
{
  "data": {
    "avg_time_to_first_review": "8.5h",
    "avg_time_to_merge": "3.2d",
    "median_time_to_merge": "2.5d",
    "percentiles": {
      "p50": "2.5d",
      "p75": "4.0d",
      "p90": "6.5d"
    }
  }
}
```

### 5. `/analytics/github/review-quality`

**Response**:
```json
{
  "data": {
    "approval_rate": 0.85,
    "avg_reviewers_per_pr": 1.8,
    "avg_comments_per_pr": 4.2,
    "prs_with_multiple_reviewers": 0.50
  }
}
```

---

## Storage Integration

### New Store Methods

```go
type Store interface {
    // Existing methods...

    // PR methods
    StorePR(pr PullRequest) error
    GetPRByID(id int) (*PullRequest, error)
    GetPRsByStatus(status PRStatus) ([]PullRequest, error)
    GetPRsByAuthor(author string) ([]PullRequest, error)

    // Review methods
    StoreReview(review Review) error
    GetReviewsByPRID(prID int) ([]Review, error)
    GetReviewsByReviewer(reviewer string) ([]Review, error)

    // Issue methods
    StoreIssue(issue Issue) error
    GetIssueByNumber(number int) (*Issue, error)
    GetIssuesByState(state IssueState) ([]Issue, error)
}
```

---

## Testing Strategy

### Unit Tests
- PR generation from commits
- Review assignment (reviewer ≠ author)
- Temporal ordering (review between PR creation and merge)
- Issue linkage to PRs

### Integration Tests
- Full PR lifecycle (create → review → merge)
- Issue resolution via PR
- Storage persistence

### E2E Tests
- API endpoints return valid data
- Pagination works correctly
- Filters apply properly

---

## Performance Considerations

- Generate PRs in batch after commit generation
- Use concurrent goroutines for review generation
- Index PRs by author, status for fast queries
- Cache PR metrics calculations

---

## Rollout Plan

1. Models (PR, Review, Issue)
2. PR Generator
3. Review Generator
4. Issue Generator
5. Storage methods
6. API handlers
7. E2E tests

---

**Estimated Effort**: 20-25 hours
