# Task Breakdown: Data Tier (ETL Pipeline)

**Feature ID**: P8-F01-data-tier
**Created**: January 9, 2026
**Status**: IN_PROGRESS (3/14 tasks)
**Approach**: TDD (Test-Driven Development)

---

## Progress Tracker

| Phase | Tasks | Status | Estimated | Actual |
|-------|-------|--------|-----------|--------|
| **Infrastructure** | 2 | ✅ 2/2 | 2.0h | 1.5h |
| **Extract Layer** | 4 | 🔄 1/4 | 6.0h | 1.0h |
| **Load Layer** | 2 | ⬜ 0/2 | 2.0h | - |
| **Transform Layer (dbt)** | 4 | ⬜ 0/4 | 8.0h | - |
| **Orchestration & Docker** | 2 | ⬜ 0/2 | 3.0h | - |
| **TOTAL** | **14** | **3/14** | **21.0h** | **2.5h** |

---

## Feature Breakdown

### PHASE 1: INFRASTRUCTURE

#### TASK-P8-01: Create Directory Structure and Dependencies

**Goal**: Set up project structure for P8 data tier

**Status**: COMPLETE
**Estimated**: 1.0h
**Actual**: 0.5h
**Commit**: 2d4cfe8

**Implementation Steps**:
1. ✅ Create directory structure:
   - `tools/api-loader/`
   - `dbt/`
   - `jobs/data-loader/`
   - `jobs/dbt-runner/`
   - `data/raw/` (gitignored)
2. ✅ Create Python requirements.txt for loader
3. ✅ Create dbt project scaffold
4. ✅ Update .gitignore for data directory

**Files Created**:
- NEW: `tools/api-loader/requirements.txt`
- NEW: `tools/api-loader/__init__.py`
- NEW: `tools/api-loader/extractors/__init__.py`
- NEW: `tools/api-loader/tests/__init__.py`
- NEW: `tools/api-loader/README.md`
- NEW: `dbt/dbt_project.yml`
- NEW: `dbt/profiles.yml`
- NEW: `dbt/packages.yml`
- NEW: `data/.gitkeep`
- NEW: `data/raw/.gitkeep`
- MODIFY: `.gitignore`

**Acceptance Criteria**:
- [x] Directory structure created
- [x] Python requirements.txt with dependencies (pandas, pyarrow, requests, duckdb, pytest)
- [x] dbt project files (dbt_project.yml, profiles.yml, packages.yml)
- [x] data/ directory is gitignored

**Notes**: Infrastructure was created by P9 agent in commit 2d4cfe8. All directories and scaffold files are in place and ready for implementation.

---

#### TASK-P8-02: Configure dbt Profiles (DuckDB + Snowflake)

**Goal**: Set up multi-target dbt profiles for dev/prod parity

**Status**: COMPLETE
**Estimated**: 1.0h
**Actual**: 1.0h
**Commit**: 91ea7cd

**Implementation Steps**:
1. ✅ dbt profiles.yml configured with dev/ci/prod targets (created in TASK-P8-01)
2. ✅ Created dbt model directory structure with placeholder SQL files
3. ✅ Created staging models (stg_commits, stg_pull_requests, stg_reviews, stg_repos)
4. ✅ Created intermediate model (int_pr_with_commits)
5. ✅ Created mart models (mart_velocity, mart_ai_impact, mart_quality, mart_review_costs)
6. ✅ Created cross-database macros (date_trunc_week, array_length, percentile_cont)
7. ✅ Created sources.yml defining raw data sources
8. ✅ Created model documentation YAMLs (_staging.yml, _intermediate.yml, _marts.yml)

**Files Created**:
- NEW: `dbt/models/sources.yml`
- NEW: `dbt/models/staging/stg_commits.sql`
- NEW: `dbt/models/staging/stg_pull_requests.sql`
- NEW: `dbt/models/staging/stg_reviews.sql`
- NEW: `dbt/models/staging/stg_repos.sql`
- NEW: `dbt/models/staging/_staging.yml`
- NEW: `dbt/models/intermediate/int_pr_with_commits.sql`
- NEW: `dbt/models/intermediate/_intermediate.yml`
- NEW: `dbt/models/marts/mart_velocity.sql`
- NEW: `dbt/models/marts/mart_ai_impact.sql`
- NEW: `dbt/models/marts/mart_quality.sql`
- NEW: `dbt/models/marts/mart_review_costs.sql`
- NEW: `dbt/models/marts/_marts.yml`
- NEW: `dbt/macros/date_trunc_week.sql`
- NEW: `dbt/macros/array_length.sql`
- NEW: `dbt/macros/percentile_cont.sql`

**Acceptance Criteria**:
- [x] dbt profiles.yml configured with dev/ci/prod targets
- [x] prod profile uses environment variables
- [x] No credentials in version control
- [x] Complete dbt model directory structure with SQL files
- [x] Staging models clean and normalize raw data
- [x] Intermediate model joins PRs with commits
- [x] Mart models provide pre-aggregated analytics
- [x] Cross-database macros for DuckDB/Snowflake compatibility
- [x] Source and model documentation complete

**Notes**:
- profiles.yml was created in TASK-P8-01 (commit 2d4cfe8)
- This task extended that work with complete dbt model structure
- All models follow dbt best practices (staging as views, intermediate as ephemeral, marts as tables)
- Models calculate cycle times, AI ratios, and quality metrics
- Ready for raw data loading and dbt build execution

---

### PHASE 2: EXTRACT LAYER (API Loader)

#### TASK-P8-03: Implement Base API Extractor

**Goal**: Base extractor for both GitHub and Cursor Analytics-style endpoints

**Status**: COMPLETE
**Estimated**: 1.5h
**Actual**: 1.0h

**TDD Approach**:
```python
# tools/api-loader/tests/test_cursor_api.py

def test_extract_commits_pagination():
    """Verify loader handles paginated commits response."""
    # Mock cursor-sim response format
    mock_response = {
        "items": [{"commitHash": "abc123", ...}],
        "totalCount": 150,
        "pagination": {"page": 1, "pageSize": 100}
    }

    extractor = CursorAPIExtractor("http://mock:8080")
    # Should extract items array, handle pagination
    commits = extractor.extract_commits()

    assert len(commits) == 150
    assert "commit_hash" in commits.columns  # snake_case

def test_extract_commits_empty():
    """Handle empty response gracefully."""
    extractor = CursorAPIExtractor("http://mock:8080")
    commits = extractor.extract_commits()

    assert len(commits) == 0
    assert isinstance(commits, pd.DataFrame)
```

**Implementation**:
```python
# tools/api-loader/extractors/cursor_api.py
class CursorAPIExtractor:
    def __init__(self, base_url: str, api_key: str = "cursor-sim-dev-key"):
        self.base_url = base_url.rstrip('/')
        self.auth = (api_key, '')

    def extract_commits(self, start_date: str = "90d") -> pd.DataFrame:
        all_items = []
        page = 1
        page_size = 500

        while True:
            resp = requests.get(
                f"{self.base_url}/analytics/ai-code/commits",
                params={"startDate": start_date, "page": page, "pageSize": page_size},
                auth=self.auth
            )
            resp.raise_for_status()
            data = resp.json()

            all_items.extend(data.get("items", []))

            if page * page_size >= data.get("totalCount", 0):
                break
            page += 1

        df = pd.DataFrame(all_items)
        df.columns = [self._to_snake_case(c) for c in df.columns]
        return df
```

**Files**:
- NEW: `tools/api-loader/extractors/base.py` (base extractor with GitHub and Cursor-style support)
- NEW: `tools/api-loader/tests/test_base.py` (comprehensive test suite)

**Implementation Details**:
- Created `BaseAPIExtractor` class supporting two response formats:
  1. **GitHub-style**: Raw arrays (e.g., `/repos`, `/repos/{o}/{r}/pulls`)
  2. **Cursor Analytics-style**: Wrapped objects with `{data: [...], pagination: {...}}`
- Pagination methods:
  - `fetch_github_style_paginated()`: Uses page/per_page params, terminates on empty array
  - `fetch_cursor_style_paginated()`: Uses page/page_size params, terminates on `hasNextPage: false`
- Single-page methods:
  - `fetch_github_style()`: For non-paginated GitHub endpoints
  - `fetch_cursor_style()`: For non-paginated Cursor endpoints
- Utility methods:
  - `write_parquet()`: Writes DataFrame to Parquet file

**Acceptance Criteria**:
- [x] Tests written before implementation (TDD)
- [x] Pagination handled correctly for both styles
- [x] Basic auth included in requests
- [x] Empty responses handled gracefully
- [x] HTTP errors raise appropriate exceptions
- [x] Parquet file writing supported

---

#### TASK-P8-04: Implement GitHub API Extractor

**Goal**: Extract repos, PRs, and reviews from cursor-sim GitHub-style endpoints

**Status**: NOT_STARTED
**Estimated**: 2.0h

**TDD Approach**:
```python
# tools/api-loader/tests/test_github_api.py

def test_extract_repos_raw_array():
    """Verify loader handles raw array (NOT wrapper object)."""
    # cursor-sim returns raw array, not {"repositories": [...]}
    mock_response = [
        {"full_name": "acme/platform", "default_branch": "main"},
        {"full_name": "acme/frontend", "default_branch": "main"}
    ]

    extractor = GitHubAPIExtractor("http://mock:8080")
    repos = extractor.extract_repositories()

    assert len(repos) == 2
    assert repos[0]["full_name"] == "acme/platform"

def test_extract_prs_raw_array():
    """Verify loader handles raw array for PRs."""
    # cursor-sim returns raw array, not {"pull_requests": [...]}
    mock_response = [
        {"number": 1, "state": "merged", ...},
        {"number": 2, "state": "open", ...}
    ]

    extractor = GitHubAPIExtractor("http://mock:8080")
    prs = extractor.extract_pull_requests(["acme/platform"])

    assert len(prs) == 2

def test_extract_reviews_raw_array():
    """Verify loader handles raw array for reviews."""
    # cursor-sim returns raw array
    extractor = GitHubAPIExtractor("http://mock:8080")
    reviews = extractor.extract_reviews("acme/platform", [1, 2])

    assert isinstance(reviews, pd.DataFrame)
```

**Implementation**:
```python
# tools/api-loader/extractors/github_api.py
class GitHubAPIExtractor:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip('/')

    def extract_repositories(self) -> pd.DataFrame:
        """GET /repos returns raw array (not wrapper object)"""
        resp = requests.get(f"{self.base_url}/repos")
        resp.raise_for_status()
        # cursor-sim returns raw array directly
        return pd.DataFrame(resp.json())

    def extract_pull_requests(self, repos: list[str]) -> pd.DataFrame:
        """GET /repos/{o}/{r}/pulls returns raw array"""
        all_prs = []
        for repo in repos:
            page = 1
            while True:
                resp = requests.get(
                    f"{self.base_url}/repos/{repo}/pulls",
                    params={"state": "all", "page": page, "per_page": 100}
                )
                resp.raise_for_status()
                # cursor-sim returns raw array
                prs = resp.json()
                if not prs:
                    break
                for pr in prs:
                    pr["repo_name"] = repo
                all_prs.extend(prs)
                if len(prs) < 100:
                    break
                page += 1
        return pd.DataFrame(all_prs)

    def extract_reviews(self, repo: str, pr_numbers: list[int]) -> pd.DataFrame:
        """GET /repos/{o}/{r}/pulls/{n}/reviews returns raw array"""
        all_reviews = []
        for pr_num in pr_numbers:
            resp = requests.get(
                f"{self.base_url}/repos/{repo}/pulls/{pr_num}/reviews"
            )
            resp.raise_for_status()
            # cursor-sim returns raw array
            reviews = resp.json()
            for r in reviews:
                r["repo_name"] = repo
                r["pr_number"] = pr_num
            all_reviews.extend(reviews)
        return pd.DataFrame(all_reviews)
```

**Files**:
- NEW: `tools/api-loader/extractors/github_api.py`
- NEW: `tools/api-loader/tests/test_github_api.py`

**Acceptance Criteria**:
- [ ] Tests written before implementation
- [ ] Handles raw array responses (NOT wrapper objects)
- [ ] Pagination for PRs handled
- [ ] repo_name and pr_number added to reviews
- [ ] All tests pass

---

#### TASK-P8-05: Implement Main Loader Script

**Goal**: Orchestrate extraction and write Parquet files

**Status**: NOT_STARTED
**Estimated**: 1.5h

**TDD Approach**:
```python
# tools/api-loader/tests/test_loader.py

def test_loader_writes_parquet_files(tmp_path):
    """Verify loader writes all Parquet files."""
    loader = DataLoader("http://mock:8080")
    loader.run(tmp_path)

    assert (tmp_path / "commits.parquet").exists()
    assert (tmp_path / "pull_requests.parquet").exists()
    assert (tmp_path / "reviews.parquet").exists()
    assert (tmp_path / "repos.parquet").exists()

def test_loader_repo_discovery_from_endpoint():
    """Verify repos come from /repos endpoint, not commits."""
    # This ensures repos with PRs but no commits are included
    loader = DataLoader("http://mock:8080")
    repos = loader._get_repo_list()

    # Should call /repos endpoint
    assert "acme/platform" in repos
```

**Implementation**:
```python
# tools/api-loader/loader.py
class DataLoader:
    def __init__(self, base_url: str):
        self.cursor = CursorAPIExtractor(base_url)
        self.github = GitHubAPIExtractor(base_url)

    def run(self, output_dir: Path):
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        # 1. Extract commits
        print("Extracting commits...")
        commits = self.cursor.extract_commits()
        commits.to_parquet(output_dir / "commits.parquet", index=False)

        # 2. Get repos from /repos endpoint (not from commits)
        print("Extracting repositories...")
        repos = self.github.extract_repositories()
        repos.to_parquet(output_dir / "repos.parquet", index=False)
        repo_names = repos["full_name"].tolist()

        # 3. Extract PRs
        print("Extracting pull requests...")
        prs = self.github.extract_pull_requests(repo_names)
        prs.to_parquet(output_dir / "pull_requests.parquet", index=False)

        # 4. Extract reviews
        print("Extracting reviews...")
        pr_numbers = prs.groupby("repo_name")["number"].apply(list).to_dict()
        all_reviews = []
        for repo, numbers in pr_numbers.items():
            reviews = self.github.extract_reviews(repo, numbers)
            all_reviews.append(reviews)
        if all_reviews:
            pd.concat(all_reviews).to_parquet(
                output_dir / "reviews.parquet", index=False
            )

        print(f"Done! Files written to {output_dir}")
```

**Files**:
- NEW: `tools/api-loader/loader.py`
- NEW: `tools/api-loader/tests/test_loader.py`

**Acceptance Criteria**:
- [ ] Tests written before implementation
- [ ] All 4 Parquet files created
- [ ] Repos from /repos endpoint (not commit-derived)
- [ ] CLI interface with --url, --output flags
- [ ] All tests pass

---

#### TASK-P8-06: Add Schema Validation

**Goal**: Validate extracted data matches expected schema

**Status**: NOT_STARTED
**Estimated**: 1.0h

**Implementation**:
```python
# tools/api-loader/schemas/commits.json
{
  "required_columns": [
    "commit_hash", "user_email", "repo_name", "tab_lines_added",
    "composer_lines_added", "non_ai_lines_added", "commit_ts"
  ]
}

# tools/api-loader/schemas/pull_requests.json
{
  "required_columns": [
    "number", "repo_name", "author_email", "state", "additions",
    "deletions", "ai_ratio", "was_reverted", "created_at"
  ]
}
```

**Files**:
- NEW: `tools/api-loader/schemas/commits.json`
- NEW: `tools/api-loader/schemas/pull_requests.json`
- NEW: `tools/api-loader/schemas/reviews.json`
- MODIFY: `tools/api-loader/loader.py` (add validation)

**Acceptance Criteria**:
- [ ] Schema files define required columns
- [ ] Loader validates output against schemas
- [ ] Clear error messages for missing columns
- [ ] Tests for validation logic

---

### PHASE 3: LOAD LAYER

#### TASK-P8-07: Implement DuckDB Loader

**Goal**: Load Parquet files into DuckDB raw schema

**Status**: NOT_STARTED
**Estimated**: 1.0h

**TDD Approach**:
```python
def test_load_parquet_to_duckdb(tmp_path):
    """Verify Parquet files load to DuckDB."""
    # Create test Parquet
    pd.DataFrame({"id": [1, 2]}).to_parquet(tmp_path / "test.parquet")

    # Load to DuckDB
    db_path = tmp_path / "test.duckdb"
    load_parquet_to_duckdb(tmp_path, db_path)

    # Verify table exists
    conn = duckdb.connect(str(db_path))
    result = conn.execute("SELECT COUNT(*) FROM raw.test").fetchone()
    assert result[0] == 2
```

**Implementation**:
```python
# tools/api-loader/duckdb_loader.py
def load_parquet_to_duckdb(parquet_dir: Path, db_path: Path):
    conn = duckdb.connect(str(db_path))
    conn.execute("CREATE SCHEMA IF NOT EXISTS raw")

    for parquet_file in parquet_dir.glob("*.parquet"):
        table = parquet_file.stem
        conn.execute(f"""
            CREATE OR REPLACE TABLE raw.{table} AS
            SELECT * FROM read_parquet('{parquet_file}')
        """)
        print(f"Loaded raw.{table}")

    conn.close()
```

**Files**:
- NEW: `tools/api-loader/duckdb_loader.py`
- NEW: `tools/api-loader/tests/test_duckdb_loader.py`

**Acceptance Criteria**:
- [ ] Tests written before implementation
- [ ] Creates raw schema
- [ ] Loads all Parquet files as tables
- [ ] Idempotent (CREATE OR REPLACE)
- [ ] All tests pass

---

#### TASK-P8-08: Create Snowflake Stage and COPY Scripts

**Goal**: SQL scripts for Snowflake production loading

**Status**: NOT_STARTED
**Estimated**: 1.0h

**Implementation**:
```sql
-- sql/snowflake/setup_stages.sql
CREATE STAGE IF NOT EXISTS cursor_analytics.raw.gcs_stage
  STORAGE_INTEGRATION = gcs_integration
  URL = 'gcs://cursor-analytics-data/raw/';

-- sql/snowflake/copy_raw_tables.sql
COPY INTO cursor_analytics.raw.commits
  FROM @cursor_analytics.raw.gcs_stage/commits/
  FILE_FORMAT = (TYPE = PARQUET)
  MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE;

COPY INTO cursor_analytics.raw.pull_requests
  FROM @cursor_analytics.raw.gcs_stage/pull_requests/
  FILE_FORMAT = (TYPE = PARQUET)
  MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE;
```

**Files**:
- NEW: `sql/snowflake/setup_stages.sql`
- NEW: `sql/snowflake/setup_raw_tables.sql`
- NEW: `sql/snowflake/copy_raw_tables.sql`

**Acceptance Criteria**:
- [ ] Stage creation script
- [ ] Raw table DDL
- [ ] COPY INTO scripts
- [ ] Documentation in README

---

### PHASE 4: TRANSFORM LAYER (dbt)

#### TASK-P8-09: Create dbt Source Definitions

**Goal**: Define raw sources for dbt

**Status**: NOT_STARTED
**Estimated**: 1.0h

**Implementation**:
```yaml
# dbt/models/sources.yml
version: 2

sources:
  - name: raw
    description: Raw data from cursor-sim extraction
    schema: raw
    tables:
      - name: commits
        description: Commit-level AI telemetry
        columns:
          - name: commit_hash
            tests: [unique, not_null]

      - name: pull_requests
        description: PR lifecycle data
        columns:
          - name: number
            tests: [not_null]
          - name: repo_name
            tests: [not_null]

      - name: reviews
        description: Review events
```

**Files**:
- NEW: `dbt/models/sources.yml`

**Acceptance Criteria**:
- [ ] All 4 raw tables defined
- [ ] Basic tests on primary keys
- [ ] Descriptions for tables and key columns

---

#### TASK-P8-10: Create dbt Staging Models

**Goal**: Clean and normalize raw data with calculated fields

**Status**: NOT_STARTED
**Estimated**: 2.5h

**TDD Approach**:
```sql
-- Test: Cycle times are calculated correctly
-- dbt/tests/assert_positive_cycle_times.sql
SELECT *
FROM {{ ref('stg_pull_requests') }}
WHERE coding_lead_time_hours < 0
   OR pickup_time_hours < 0
   OR review_lead_time_hours < 0
```

**Implementation**:
```sql
-- dbt/models/staging/stg_pull_requests.sql
WITH source AS (
    SELECT * FROM {{ source('raw', 'pull_requests') }}
),

calculated AS (
    SELECT
        number AS pr_number,
        repo_name,
        author_email,
        state,
        additions,
        deletions,
        (additions + deletions) AS total_loc,
        changed_files,
        ai_ratio,

        -- Rename cursor-sim field to standard name
        was_reverted AS is_reverted,
        is_bug_fix,

        -- Timestamps
        created_at,
        merged_at,
        first_commit_at,
        first_review_at,

        -- CALCULATE cycle times (not from API)
        EXTRACT(EPOCH FROM (created_at - first_commit_at)) / 3600
            AS coding_lead_time_hours,
        EXTRACT(EPOCH FROM (first_review_at - created_at)) / 3600
            AS pickup_time_hours,
        EXTRACT(EPOCH FROM (merged_at - first_review_at)) / 3600
            AS review_lead_time_hours,

        -- Calculate reviewer count from array
        {{ array_length('reviewers') }} AS reviewer_count

    FROM source
    WHERE created_at IS NOT NULL
)

SELECT * FROM calculated
```

**Files**:
- NEW: `dbt/models/staging/stg_commits.sql`
- NEW: `dbt/models/staging/stg_pull_requests.sql`
- NEW: `dbt/models/staging/stg_reviews.sql`
- NEW: `dbt/models/staging/_staging.yml`

**Acceptance Criteria**:
- [ ] stg_commits with ai_lines_added calculated
- [ ] stg_pull_requests with cycle times calculated
- [ ] stg_reviews with is_approval flag
- [ ] was_reverted renamed to is_reverted
- [ ] All dbt tests pass

---

#### TASK-P8-11: Create dbt Intermediate Models

**Goal**: Join PRs with commit aggregations

**Status**: NOT_STARTED
**Estimated**: 1.5h

**Implementation**:
```sql
-- dbt/models/intermediate/int_pr_with_commits.sql
WITH prs AS (
    SELECT * FROM {{ ref('stg_pull_requests') }}
),

commit_agg AS (
    SELECT
        pull_request_number,
        repo_name,
        COUNT(*) AS commit_count,
        SUM(ai_lines_added) AS total_ai_lines,
        SUM(total_lines_added) AS total_lines
    FROM {{ ref('stg_commits') }}
    WHERE pull_request_number IS NOT NULL
    GROUP BY 1, 2
)

SELECT
    p.*,
    COALESCE(c.commit_count, 0) AS commit_count,
    COALESCE(c.total_ai_lines, 0) AS pr_ai_lines,
    COALESCE(
        p.ai_ratio,
        c.total_ai_lines::FLOAT / NULLIF(c.total_lines, 0),
        0
    ) AS final_ai_ratio
FROM prs p
LEFT JOIN commit_agg c
    ON p.pr_number = c.pull_request_number
    AND p.repo_name = c.repo_name
```

**Files**:
- NEW: `dbt/models/intermediate/int_pr_with_commits.sql`
- NEW: `dbt/models/intermediate/_intermediate.yml`

**Acceptance Criteria**:
- [ ] PRs joined with commit aggregations
- [ ] final_ai_ratio calculated with fallbacks
- [ ] Materialized as ephemeral

---

#### TASK-P8-12: Create dbt Mart Models

**Goal**: Pre-aggregated analytics tables

**Status**: NOT_STARTED
**Estimated**: 3.0h

**Implementation**:
```sql
-- dbt/models/marts/mart_velocity.sql
WITH pr_data AS (
    SELECT * FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
)

SELECT
    {{ date_trunc_week('merged_at') }} AS week,
    repo_name,
    COUNT(DISTINCT author_email) AS active_developers,
    COUNT(*) AS total_prs,
    AVG(total_loc) AS avg_pr_size,
    AVG(coding_lead_time_hours) AS coding_lead_time,
    AVG(pickup_time_hours) AS pickup_time,
    AVG(review_lead_time_hours) AS review_lead_time,
    AVG(coding_lead_time_hours + pickup_time_hours + review_lead_time_hours)
        AS total_cycle_time,
    {{ percentile_cont(0.5, 'coding_lead_time_hours + pickup_time_hours + review_lead_time_hours') }}
        AS p50_cycle_time,
    AVG(final_ai_ratio) AS avg_ai_ratio
FROM pr_data
GROUP BY 1, 2


-- dbt/models/marts/mart_ai_impact.sql
WITH pr_data AS (
    SELECT
        *,
        CASE
            WHEN final_ai_ratio < 0.3 THEN 'low'
            WHEN final_ai_ratio < 0.6 THEN 'medium'
            ELSE 'high'
        END AS ai_usage_band
    FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
)

SELECT
    ai_usage_band,
    {{ date_trunc_week('merged_at') }} AS week,
    COUNT(*) AS pr_count,
    AVG(final_ai_ratio) AS avg_ai_ratio,
    AVG(coding_lead_time_hours) AS avg_coding_lead_time,
    AVG(pickup_time_hours + review_lead_time_hours) AS avg_review_cycle_time,
    AVG(CASE WHEN is_reverted THEN 1 ELSE 0 END) AS revert_rate
FROM pr_data
GROUP BY 1, 2
```

**Files**:
- NEW: `dbt/models/marts/mart_velocity.sql`
- NEW: `dbt/models/marts/mart_review_costs.sql`
- NEW: `dbt/models/marts/mart_quality.sql`
- NEW: `dbt/models/marts/mart_ai_impact.sql`
- NEW: `dbt/models/marts/_marts.yml`

**Acceptance Criteria**:
- [ ] mart_velocity with weekly aggregations
- [ ] mart_ai_impact with AI band grouping
- [ ] mart_quality with revert/bug fix rates
- [ ] All materialized as tables
- [ ] Cross-engine macros used
- [ ] dbt tests pass

---

### PHASE 5: ORCHESTRATION & DOCKER

#### TASK-P8-13: Create Pipeline Script and Makefile

**Goal**: Single command to run full pipeline

**Status**: NOT_STARTED
**Estimated**: 1.5h

**Implementation**:
```bash
# tools/run_pipeline.sh
#!/bin/bash
set -e

CURSOR_SIM_URL=${CURSOR_SIM_URL:-"http://localhost:8080"}
DATA_DIR=${DATA_DIR:-"./data"}

echo "=== Step 1/3: Extract ==="
python tools/api-loader/loader.py --url "$CURSOR_SIM_URL" --output "$DATA_DIR/raw"

echo "=== Step 2/3: Load to DuckDB ==="
python tools/api-loader/duckdb_loader.py --input "$DATA_DIR/raw" --db "$DATA_DIR/analytics.duckdb"

echo "=== Step 3/3: Run dbt ==="
cd dbt && dbt deps && dbt build --target dev

echo "=== Pipeline complete ==="
```

```makefile
# Makefile additions
pipeline:
	./tools/run_pipeline.sh

pipeline-ci:
	CURSOR_SIM_URL=http://localhost:8080 ./tools/run_pipeline.sh
```

**Files**:
- NEW: `tools/run_pipeline.sh`
- MODIFY: `Makefile`

**Acceptance Criteria**:
- [ ] Single command runs full pipeline
- [ ] Environment variable configuration
- [ ] Error handling (set -e)
- [ ] Progress output

---

#### TASK-P8-14: Create Production Job Dockerfiles

**Goal**: Docker containers for Cloud Run Jobs

**Status**: NOT_STARTED
**Estimated**: 1.5h

**Implementation**:
```dockerfile
# jobs/data-loader/Dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY tools/api-loader/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY tools/api-loader/ ./
ENTRYPOINT ["python", "loader.py"]

# jobs/dbt-runner/Dockerfile
FROM python:3.11-slim
WORKDIR /app
RUN pip install --no-cache-dir dbt-snowflake
COPY dbt/ ./dbt/
COPY jobs/dbt-runner/run.sh ./
RUN chmod +x run.sh
ENTRYPOINT ["./run.sh"]
```

**Files**:
- NEW: `jobs/data-loader/Dockerfile`
- NEW: `jobs/data-loader/requirements.txt`
- NEW: `jobs/dbt-runner/Dockerfile`
- NEW: `jobs/dbt-runner/run.sh`

**Acceptance Criteria**:
- [ ] data-loader container builds
- [ ] dbt-runner container builds
- [ ] Non-root users
- [ ] Minimal image sizes

---

## Dependency Graph

```
TASK-P8-01 (Infrastructure)
    │
    ├──► TASK-P8-02 (dbt Profiles)
    │         │
    │         └──► TASK-P8-09 (dbt Sources)
    │                   │
    │                   └──► TASK-P8-10 (Staging)
    │                             │
    │                             └──► TASK-P8-11 (Intermediate)
    │                                       │
    │                                       └──► TASK-P8-12 (Marts)
    │
    └──► TASK-P8-03 (Cursor Extractor)
              │
              └──► TASK-P8-04 (GitHub Extractor)
                        │
                        └──► TASK-P8-05 (Main Loader)
                                  │
                                  ├──► TASK-P8-06 (Schema Validation)
                                  │
                                  └──► TASK-P8-07 (DuckDB Loader)
                                            │
                                            └──► TASK-P8-08 (Snowflake Scripts)

TASK-P8-12 + TASK-P8-07 ──► TASK-P8-13 (Pipeline Script)
                                  │
                                  └──► TASK-P8-14 (Docker Jobs)
```

---

## Testing Strategy

### Unit Tests (Python)

| Component | Target Coverage |
|-----------|-----------------|
| cursor_api.py | 90% |
| github_api.py | 90% |
| loader.py | 85% |
| duckdb_loader.py | 90% |

### dbt Tests

| Model | Tests |
|-------|-------|
| stg_commits | unique, not_null, accepted_range |
| stg_pull_requests | positive cycle times, valid ai_ratio |
| mart_velocity | not_null on aggregations |

### Integration Tests

```bash
# Full pipeline test
make ci-local

# Verify results
duckdb data/analytics.duckdb <<EOF
SELECT COUNT(*) FROM mart.velocity;
SELECT COUNT(*) FROM mart.ai_impact WHERE ai_usage_band = 'high';
EOF
```

---

## Definition of Done (Per Task)

- [ ] Tests written BEFORE implementation (TDD)
- [ ] All tests pass
- [ ] Code review completed
- [ ] Dependency reflections checked
- [ ] Git commit with descriptive message
- [ ] task.md updated with status

---

## Success Criteria (Feature Complete)

- [ ] All 14 tasks completed
- [ ] Loader extracts from cursor-sim correctly
- [ ] DuckDB populated with raw and mart tables
- [ ] dbt models pass all tests
- [ ] Pipeline runs end-to-end in < 5 minutes
- [ ] Production Dockerfiles build successfully
- [ ] Documentation complete

---

**Next Action**: Start with TASK-P8-01 (Infrastructure)
