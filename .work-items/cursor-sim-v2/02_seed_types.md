# Step 02: Implement Seed Schema Types

## Objective

Define Go types matching the seed.json schema from DataDesigner.

## Estimated Time: 2 hours
## Recommended Model: Haiku

## Prerequisites

- Step 01: Project structure complete

## Acceptance Criteria

- [ ] All seed types defined with JSON tags
- [ ] Types match tools/data-designer/config/seed_schema.yaml exactly
- [ ] Unit tests verify JSON marshaling/unmarshaling
- [ ] Test with sample seed data passes
- [ ] 90%+ test coverage

## Implementation Steps

### 1. Write Failing Tests First (RED)

```go
// internal/seed/types_test.go
package seed

import (
    "encoding/json"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDeveloper_JSONRoundtrip(t *testing.T) {
    dev := Developer{
        UserID:         "user_001",
        Email:          "dev@example.com",
        Name:           "Jane Developer",
        Seniority:      "senior",
        AcceptanceRate: 0.85,
    }

    data, err := json.Marshal(dev)
    require.NoError(t, err)

    var parsed Developer
    err = json.Unmarshal(data, &parsed)
    require.NoError(t, err)

    assert.Equal(t, dev, parsed)
}

func TestSeedData_JSONTags(t *testing.T) {
    jsonStr := `{
        "developers": [{
            "user_id": "user_001",
            "email": "dev@example.com",
            "acceptance_rate": 0.85
        }]
    }`

    var seed SeedData
    err := json.Unmarshal([]byte(jsonStr), &seed)
    require.NoError(t, err)

    assert.Equal(t, "user_001", seed.Developers[0].UserID)
    assert.Equal(t, 0.85, seed.Developers[0].AcceptanceRate)
}
```

### 2. Define Types (GREEN)

```go
// internal/seed/types.go
package seed

type SeedData struct {
    Developers    []Developer    `json:"developers"`
    Repositories  []Repository   `json:"repositories"`
    Correlations  Correlations   `json:"correlations"`
    TextTemplates TextTemplates  `json:"text_templates"`
}

type Developer struct {
    UserID         string     `json:"user_id"`
    Email          string     `json:"email"`
    Name           string     `json:"name"`
    Org            string     `json:"org"`
    Division       string     `json:"division"`
    Team           string     `json:"team"`
    Seniority      string     `json:"seniority"`
    Region         string     `json:"region"`
    AcceptanceRate float64    `json:"acceptance_rate"`
    PRBehavior     PRBehavior `json:"pr_behavior"`
}

type PRBehavior struct {
    PRsPerWeek      float64 `json:"prs_per_week"`
    AvgPRSizeLoc    int     `json:"avg_pr_size_loc"`
    GreenfieldRatio float64 `json:"greenfield_ratio"`
}

type Repository struct {
    RepoName        string   `json:"repo_name"`
    PrimaryLanguage string   `json:"primary_language"`
    AgeDays         int      `json:"age_days"`
    Maturity        string   `json:"maturity"`
    Teams           []string `json:"teams"`
}

type Correlations struct {
    SeniorityAcceptanceRate map[string]float64 `json:"seniority_acceptance_rate"`
    AIRatioRevertRate       map[string]float64 `json:"ai_ratio_revert_rate"`
}

type TextTemplates struct {
    CommitMessages []string `json:"commit_messages"`
    PRTitles       []string `json:"pr_titles"`
    ReviewComments []string `json:"review_comments"`
}
```

## Test Cases

1. `TestDeveloper_JSONRoundtrip` - Marshal/unmarshal preserves data
2. `TestSeedData_JSONTags` - JSON tags map correctly (snake_case)
3. `TestPRBehavior_Defaults` - Zero values handled correctly
4. `TestRepository_TeamsArray` - Array fields work
5. `TestCorrelations_MapFields` - Map fields unmarshal correctly

## Files to Create

- internal/seed/types.go
- internal/seed/types_test.go

## Definition of Done

- All tests pass
- Coverage > 90%
- Types match seed_schema.yaml
