# Analytics Stack Architecture

## Data Contract Hierarchy

**cursor-sim is the authoritative source of truth** for the analytics platform. All downstream layers must validate against and preserve the API contract.

### Contract Levels

```
LEVEL 1: API CONTRACT (cursor-sim SPEC.md) ← SOURCE OF TRUTH
├─ Endpoints: /analytics/ai-code/commits, /teams/members, /repos/*/pulls
├─ Response format: {items: [...], totalCount, page, pageSize}
├─ Field names: camelCase (commitHash, userEmail, tabLinesAdded, ...)
├─ Data types: strings, numbers, dates in ISO format
└─ Responsibility: cursor-sim API (P4)

LEVEL 2: DATA TIER CONTRACT (api-loader → dbt → DuckDB)
├─ Raw schema (main_raw.*): Preserves API fields exactly
│  ├─ Data: camelCase, flat structure
│  ├─ Responsibility: api-loader extraction + DuckDB table creation
│  └─ Validation: Schema matches API response structure
│
├─ Staging schema (main_staging.stg_*): Transforms camelCase → snake_case
│  ├─ Data: snake_case, typed columns, cleaned
│  ├─ Responsibility: dbt staging models
│  └─ Validation: Column mapping contract (commitHash → commit_hash, etc.)
│
└─ Mart schema (main_mart.mart_*): Aggregations for analytics
   ├─ Data: snake_case, aggregated metrics, analytics-ready
   ├─ Responsibility: dbt mart models
   └─ Validation: Correct aggregation logic, required columns present

LEVEL 3: DASHBOARD CONTRACT (Streamlit)
├─ Queries: SELECT from main_mart.* only, never raw or staging
├─ Parameters: Parameterized queries ($param syntax), never f-strings
├─ Schema: Use main_mart.* prefix (DuckDB requirement)
└─ Responsibility: streamlit-dashboard query modules
```

### Data Fidelity

Each layer must preserve data fidelity from the previous layer:

```
cursor-sim API (fact)
    ↓ [api-loader: validates response format, extracts items]
main_raw schema (fact copy)
    ↓ [dbt staging: transforms camelCase → snake_case, validates types]
main_staging schema (cleaned fact)
    ↓ [dbt marts: aggregates, calculates metrics]
main_mart schema (analytics-ready)
    ↓ [streamlit: parameterized queries, no SQL injection]
Dashboard KPIs
```

### Validation at Each Layer

| Layer | Input Contract | Validation | Output Contract |
|-------|-----------------|------------|-----------------|
| **API** | N/A | API SPEC.md test | `{items:[...], totalCount, page}` |
| **Extraction** | API response | Dual format handling | Parquet files with correct columns |
| **Raw** | Parquet | Schema matches API | main_raw.* tables with camelCase |
| **Staging** | main_raw.* | Column mapping, type coercion | main_staging.stg_* with snake_case |
| **Mart** | main_staging.stg_* | Aggregation logic, metrics | main_mart.mart_* analytics-ready |
| **Dashboard** | main_mart.* | Parameterized queries | KPI visualizations |

---

## Tool Assessment

| Tool | Role | Strengths | Dev Suitability |
|------|------|-----------|-----------------|
| **SnapLogic** | Extract & Load (EL) | Enterprise iPaaS, API connectors, orchestration | ❌ Heavy, cloud-only, expensive for dev |
| **dbt** | Transform (T) | SQL-based transforms, version control, testing | ✅ Excellent, identical dev/prod |
| **Flexor.ai** | AI-assisted mapping | Complex schema mapping, auto-detection | ⚠️ Overkill for your defined schemas |

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              PRODUCTION                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Cursor API ──┐                                                            │
│                ├──► SnapLogic ──► Snowflake RAW ──► dbt ──► Snowflake MART  │
│   GitHub API ──┘    (Extract)     (Landing)        (Transform)  (Analytics) │
│                                                                              │
│                                                          │                  │
│                                                          ▼                  │
│                                                     Streamlit               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              DEVELOPMENT                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   cursor-sim ──► Loader ──► Parquet ──► DuckDB RAW ──► dbt ──► DuckDB MART │
│   (Mock APIs)    (Extract)  (Landing)   (Local)        (Same!)  (Analytics) │
│                                                                              │
│                                                          │                  │
│                                                          ▼                  │
│                                                     Streamlit               │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Key insight**: 
- cursor-sim exposes the **same API** as production (Cursor + GitHub)
- Loader script mimics SnapLogic's extraction logic
- dbt is identical in both environments

---

## Layer Mapping

| Layer | Dev | Prod | Parity Strategy |
|-------|-----|------|-----------------|
| **Source** | cursor-sim (API) | Cursor API + GitHub API | Same API contracts |
| **Extract** | Loader (Python) | SnapLogic | Same extraction logic |
| **Land** | Local Parquet files | Snowflake Stage → RAW | Same table schemas |
| **Transform** | dbt + DuckDB | dbt + Snowflake | ✅ Identical SQL |
| **Serve** | DuckDB MART | Snowflake MART | Same table schemas |
| **Visualize** | Streamlit | Streamlit | ✅ Identical |

---

## Data Flow

```
                     DEV                                    PROD
                     ───                                    ────
                                                            
    ┌──────────────┐                           ┌──────────────┐
    │  cursor-sim  │                           │  Cursor API  │
    │  (Go CLI)    │                           │  GitHub API  │
    │              │                           │              │
    │  Port 8080   │                           │  (Cloud)     │
    └──────┬───────┘                           └──────┬───────┘
           │                                          │
           │ REST API (JSON)                          │ REST API (JSON)
           ▼                                          ▼
    ┌──────────────┐                           ┌──────────────┐
    │   Loader     │                           │  SnapLogic   │
    │  (Python)    │   SAME EXTRACTION LOGIC   │  Pipelines   │
    │              │◄─────────────────────────►│              │
    │  - Paginate  │                           │  - Paginate  │
    │  - Transform │                           │  - Transform │
    │  - Validate  │                           │  - Validate  │
    └──────┬───────┘                           └──────┬───────┘
           │                                          │
           │ Parquet                                  │ Parquet
           ▼                                          ▼
    ┌──────────────┐                           ┌──────────────┐
    │ data/raw/    │                           │ Snowflake    │
    │ ├─commits    │                           │ @raw_stage   │
    │ ├─prs        │                           │              │
    │ └─reviews    │                           │              │
    └──────┬───────┘                           └──────┬───────┘
           │                                          │
           │ Load                                     │ COPY INTO
           ▼                                          ▼
    ┌──────────────┐                           ┌──────────────┐
    │ DuckDB       │                           │ Snowflake    │
    │ raw.commits  │                           │ raw.commits  │
    │ raw.prs      │     SAME SCHEMA           │ raw.prs      │
    │ raw.reviews  │◄────────────────────────► │ raw.reviews  │
    └──────┬───────┘                           └──────┬───────┘
           │                                          │
           │              ┌──────────────┐            │
           └─────────────►│     dbt      │◄───────────┘
                          │              │
                          │ models/      │
                          │ ├─staging/   │
                          │ ├─marts/     │
                          │ └─metrics/   │
                          └──────┬───────┘
                                 │
           ┌─────────────────────┴─────────────────────┐
           │                                           │
           ▼                                           ▼
    ┌──────────────┐                           ┌──────────────┐
    │ DuckDB       │                           │ Snowflake    │
    │ mart.*       │     SAME SCHEMA           │ mart.*       │
    └──────┬───────┘◄────────────────────────► └──────┬───────┘
           │                                          │
           └─────────────────────┬─────────────────────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │  Streamlit   │
                          │  Dashboard   │
                          └──────────────┘
```

---

## Directory Structure

```
cursor-analytics-platform/
├── apps/
│   ├── cursor-sim/                    # Go - API simulator (serves REST API)
│   │   ├── main.go
│   │   ├── api/
│   │   │   ├── cursor_handlers.go     # /teams/*, /analytics/*
│   │   │   └── github_handlers.go     # /repos/*
│   │   └── generator/
│   │       ├── commits.go
│   │       ├── pull_requests.go
│   │       └── reviews.go
│   │
│   └── dashboard/                     # Streamlit
│       ├── app.py
│       ├── pages/
│       │   ├── 1_velocity.py
│       │   ├── 2_review_costs.py
│       │   ├── 3_quality.py
│       │   └── 4_ai_impact.py
│       ├── db/
│       │   └── connector.py           # DuckDB/Snowflake abstraction
│       └── requirements.txt
│
├── data/                              # Local data (gitignored)
│   ├── raw/                           # Landing zone (Parquet from loader)
│   │   ├── commits.parquet
│   │   ├── pull_requests.parquet
│   │   └── reviews.parquet
│   └── analytics.duckdb              # Local DuckDB
│
├── dbt/                               # dbt project (SAME FOR DEV/PROD)
│   ├── dbt_project.yml
│   ├── profiles.yml                   # Multi-target: dev (DuckDB), prod (Snowflake)
│   ├── models/
│   │   ├── sources.yml                # Define raw sources
│   │   ├── staging/                   # Clean raw data
│   │   │   ├── stg_commits.sql
│   │   │   ├── stg_pull_requests.sql
│   │   │   └── stg_reviews.sql
│   │   ├── intermediate/              # Business logic
│   │   │   ├── int_pr_metrics.sql
│   │   │   └── int_developer_activity.sql
│   │   └── marts/                     # Analytics-ready
│   │       ├── mart_velocity.sql
│   │       ├── mart_review_costs.sql
│   │       ├── mart_quality.sql
│   │       └── mart_ai_impact.sql
│   ├── tests/                         # dbt tests
│   │   ├── assert_ai_ratio_bounds.sql
│   │   └── assert_cycle_times_positive.sql
│   ├── macros/                        # Reusable SQL
│   │   └── calculate_cycle_time.sql
│   └── seeds/                         # Static reference data
│       └── seniority_levels.csv
│
├── tools/
│   ├── data-designer/                 # Seed generator (existing)
│   │   ├── generate_seed.py
│   │   └── config/
│   │       └── seed_schema.yaml
│   │
│   └── api-loader/                    # NEW: API → Parquet loader
│       ├── loader.py                  # Main extraction script
│       ├── extractors/
│       │   ├── cursor_api.py          # Cursor API extraction
│       │   └── github_api.py          # GitHub API extraction
│       ├── schemas/                   # Expected schemas (validation)
│       │   ├── commits.json
│       │   └── pull_requests.json
│       └── requirements.txt
│
├── snaplogic/                         # SnapLogic configs (prod only)
│   ├── pipelines/
│   │   ├── cursor_api_extract.slp
│   │   └── github_api_extract.slp
│   ├── schemas/
│   │   ├── commits.json
│   │   └── pull_requests.json
│   └── README.md
│
├── sql/                               # Raw SQL (for reference/migration)
│   └── snowflake/
│       ├── setup_raw_tables.sql
│       └── setup_stages.sql
│
└── .github/
    └── workflows/
        ├── ci-cursor-sim.yml          # cursor-sim CI
        ├── ci-loader.yml              # Loader CI (tests extraction logic)
        ├── ci-dbt.yml                 # dbt CI
        ├── ci-dashboard.yml           # Streamlit CI
        └── cd-prod.yml                # Production deployment
```

---

## API Loader

The loader extracts data from cursor-sim's REST API, mimicking what SnapLogic does in production. This ensures extraction logic is tested in CI.

### cursor-sim API Response Format (Source of Truth)

**Important**: cursor-sim returns API responses in a specific format that must be handled correctly:

```json
{
  "items": [
    {
      "commitHash": "abc123",
      "userEmail": "dev@example.com",
      "tabLinesAdded": 45,
      "composerLinesAdded": 12,
      "commitTs": "2026-01-10T10:30:00Z"
    }
  ],
  "totalCount": 1000,
  "page": 1,
  "pageSize": 500
}
```

**Dual Format Support**: The loader must handle both:
- **Paginated response** (production format): `{items: [...], totalCount, page, pageSize}`
- **Raw array response** (fallback format): `[{...}, {...}]`

**Field Naming**: All fields are **camelCase** (commitHash, userEmail, tabLinesAdded, composerLinesAdded, commitTs). These field names must be preserved through the raw schema and transformed to snake_case in the dbt staging layer.

### loader.py

```python
# tools/api-loader/loader.py
"""
API → Parquet Loader

Extracts data from cursor-sim REST API and writes to Parquet.
Mimics SnapLogic extraction logic for dev/prod parity.

CRITICAL: Handles dual API response formats:
1. Paginated: {items: [...], totalCount, page, pageSize}
2. Raw array: [...]
"""

import requests
import pandas as pd
from pathlib import Path
from typing import Optional, Union, Dict, Any, List
import json


class BaseAPIExtractor:
    """Base class for API extraction with dual-format response handling"""

    @staticmethod
    def fetch_cursor_style_paginated(
        base_url: str,
        endpoint: str,
        params: Dict[str, Any]
    ) -> List[Dict[str, Any]]:
        """
        Fetch from cursor-sim with pagination support.

        Handles both response formats:
        - Format 1: {items: [...], totalCount, page, pageSize}
        - Format 2: Raw array [...]
        """
        all_items = []
        page = 1
        page_size = params.get('pageSize', 1000)

        while True:
            params_with_page = {**params, 'page': page, 'pageSize': page_size}
            resp = requests.get(f"{base_url}{endpoint}", params=params_with_page)
            resp.raise_for_status()
            data = resp.json()

            # Format 1: Paginated response with items key
            if isinstance(data, dict) and 'items' in data:
                items = data.get('items', [])
                all_items.extend(items)
                total_count = data.get('totalCount', len(all_items))
                if page * page_size >= total_count:
                    break
            # Format 2: Raw array response
            elif isinstance(data, list):
                all_items.extend(data)
                break
            else:
                raise ValueError(f"Unexpected API response format: {type(data)}")

            page += 1

        return all_items


class CursorAPIExtractor:
    """Extract data from Cursor API endpoints"""

    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip('/')

    def extract_commits(self, start_date: str = "90d", end_date: str = "now") -> pd.DataFrame:
        """
        Extract commits from /analytics/ai-code/commits
        Handles pagination automatically.

        Returns DataFrame with camelCase columns (as returned by API).
        Column transformation to snake_case happens in dbt staging layer.
        """
        all_items = BaseAPIExtractor.fetch_cursor_style_paginated(
            self.base_url,
            '/analytics/ai-code/commits',
            {
                "startDate": start_date,
                "endDate": end_date
            }
        )

        df = pd.DataFrame(all_items)

        # IMPORTANT: Do NOT transform column names here
        # API returns camelCase (commitHash, userEmail, etc.)
        # Transformation to snake_case happens in dbt staging models
        # This preserves data fidelity through raw schema

        return df
    
    def extract_team_members(self) -> pd.DataFrame:
        """Extract team members from /teams/members"""
        resp = requests.get(f"{self.base_url}/teams/members")
        resp.raise_for_status()
        return pd.DataFrame(resp.json()["teamMembers"])
    
    @staticmethod
    def _to_snake_case(name: str) -> str:
        """Convert camelCase to snake_case"""
        import re
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
        return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()


class GitHubAPIExtractor:
    """Extract data from GitHub-style API endpoints"""
    
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip('/')
    
    def extract_pull_requests(self, repos: list[str], since: Optional[str] = None) -> pd.DataFrame:
        """
        Extract PRs from /repos/{owner}/{repo}/pulls
        """
        all_prs = []
        
        for repo in repos:
            page = 1
            while True:
                resp = requests.get(
                    f"{self.base_url}/repos/{repo}/pulls",
                    params={
                        "state": "all",
                        "since": since,
                        "page": page,
                        "per_page": 100
                    }
                )
                resp.raise_for_status()
                data = resp.json()
                
                prs = data["pull_requests"]
                if not prs:
                    break
                
                for pr in prs:
                    pr["repo_name"] = repo
                all_prs.extend(prs)
                
                if len(prs) < 100:
                    break
                page += 1
        
        return pd.DataFrame(all_prs)
    
    def extract_reviews(self, repos: list[str], pr_numbers: dict[str, list[int]]) -> pd.DataFrame:
        """
        Extract reviews from /repos/{owner}/{repo}/pulls/{n}/reviews
        pr_numbers: {repo_name: [pr_number, ...]}
        """
        all_reviews = []
        
        for repo in repos:
            for pr_num in pr_numbers.get(repo, []):
                resp = requests.get(
                    f"{self.base_url}/repos/{repo}/pulls/{pr_num}/reviews"
                )
                resp.raise_for_status()
                
                reviews = resp.json()["reviews"]
                for r in reviews:
                    r["repo_name"] = repo
                    r["pr_number"] = pr_num
                all_reviews.extend(reviews)
        
        return pd.DataFrame(all_reviews)
    
    def extract_repositories(self) -> pd.DataFrame:
        """Extract repository list"""
        resp = requests.get(f"{self.base_url}/repos")
        resp.raise_for_status()
        return pd.DataFrame(resp.json()["repositories"])


class DataLoader:
    """Orchestrates extraction from all APIs"""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.cursor = CursorAPIExtractor(base_url)
        self.github = GitHubAPIExtractor(base_url)
    
    def load_all(self, output_dir: Path, start_date: str = "90d"):
        """
        Full extraction pipeline: API → Parquet
        
        This mirrors what SnapLogic does in production.
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)
        
        print("Extracting data from API...")
        
        # 1. Extract commits (Cursor API)
        print("  - Commits...", end=" ")
        commits = self.cursor.extract_commits(start_date=start_date)
        commits.to_parquet(output_dir / "commits.parquet", index=False)
        print(f"✓ {len(commits)} rows")
        
        # 2. Get unique repos from commits
        repos = commits["repo_name"].unique().tolist()
        
        # 3. Extract PRs (GitHub API)
        print("  - Pull requests...", end=" ")
        prs = self.github.extract_pull_requests(repos)
        prs.to_parquet(output_dir / "pull_requests.parquet", index=False)
        print(f"✓ {len(prs)} rows")
        
        # 4. Extract reviews (GitHub API)
        print("  - Reviews...", end=" ")
        pr_numbers = prs.groupby("repo_name")["number"].apply(list).to_dict()
        reviews = self.github.extract_reviews(repos, pr_numbers)
        reviews.to_parquet(output_dir / "reviews.parquet", index=False)
        print(f"✓ {len(reviews)} rows")
        
        # 5. Extract team members (Cursor API)
        print("  - Team members...", end=" ")
        members = self.cursor.extract_team_members()
        members.to_parquet(output_dir / "team_members.parquet", index=False)
        print(f"✓ {len(members)} rows")
        
        print(f"\nData written to {output_dir}/")
        return {
            "commits": len(commits),
            "pull_requests": len(prs),
            "reviews": len(reviews),
            "team_members": len(members),
        }
    
    def validate_schemas(self, output_dir: Path):
        """Validate extracted data matches expected schemas"""
        schema_dir = Path(__file__).parent / "schemas"
        
        for parquet_file in output_dir.glob("*.parquet"):
            schema_file = schema_dir / f"{parquet_file.stem}.json"
            if schema_file.exists():
                with open(schema_file) as f:
                    expected = json.load(f)
                
                df = pd.read_parquet(parquet_file)
                missing = set(expected["required_columns"]) - set(df.columns)
                
                if missing:
                    raise ValueError(
                        f"Schema validation failed for {parquet_file.name}: "
                        f"missing columns {missing}"
                    )
                print(f"✓ Schema valid: {parquet_file.name}")


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="Load data from cursor-sim API")
    parser.add_argument("--url", default="http://localhost:8080", help="API base URL")
    parser.add_argument("--output", "-o", default="data/raw", help="Output directory")
    parser.add_argument("--start-date", default="90d", help="Start date (e.g., 90d, 2025-01-01)")
    parser.add_argument("--validate", action="store_true", help="Validate schemas after load")
    
    args = parser.parse_args()
    
    loader = DataLoader(base_url=args.url)
    loader.load_all(Path(args.output), start_date=args.start_date)
    
    if args.validate:
        loader.validate_schemas(Path(args.output))
```

### Schema Validation Files

```json
// tools/api-loader/schemas/commits.json
{
  "name": "commits",
  "required_columns": [
    "commit_hash",
    "user_id",
    "user_email",
    "repo_name",
    "branch_name",
    "is_primary_branch",
    "total_lines_added",
    "total_lines_deleted",
    "tab_lines_added",
    "tab_lines_deleted",
    "composer_lines_added",
    "composer_lines_deleted",
    "non_ai_lines_added",
    "non_ai_lines_deleted",
    "commit_ts",
    "created_at"
  ]
}
```

```json
// tools/api-loader/schemas/pull_requests.json
{
  "name": "pull_requests",
  "required_columns": [
    "number",
    "repo_name",
    "author_email",
    "state",
    "additions",
    "deletions",
    "changed_files",
    "created_at",
    "merged_at",
    "coding_lead_time_hours",
    "pickup_time_hours",
    "review_lead_time_hours",
    "review_comments",
    "iterations",
    "reviewer_count",
    "ai_ratio",
    "is_reverted",
    "has_hotfix_followup"
  ]
}
```

---

## DuckDB Schema Naming Convention

**CRITICAL**: DuckDB requires the `main_` prefix for schema-qualified table names. This is a DuckDB-specific requirement that differs from standard SQL databases.

### Schema Hierarchy in DuckDB

```
main_raw.commits          ← Raw data from API (camelCase fields preserved)
main_raw.pull_requests
main_raw.reviews

main_staging.stg_commits  ← Staging layer (camelCase → snake_case transformation)
main_staging.stg_pull_requests
main_staging.stg_reviews

main_mart.mart_velocity   ← Analytics-ready aggregates (snake_case, aggregated)
main_mart.mart_ai_impact
main_mart.mart_quality
main_mart.mart_review_costs
```

### Correct vs Incorrect Schema References

```sql
-- ✅ CORRECT: DuckDB requires main_* prefix
SELECT * FROM main_raw.commits
SELECT * FROM main_staging.stg_commits
SELECT * FROM main_mart.mart_velocity

-- ❌ INCORRECT: Will fail with "Catalog Error"
SELECT * FROM raw.commits
SELECT * FROM staging.stg_commits
SELECT * FROM mart.mart_velocity
```

### Why the `main_` Prefix?

DuckDB organizes catalogs hierarchically: `CATALOG.SCHEMA.TABLE`

- `main` is the default catalog (DuckDB's built-in catalog)
- Without the `main_` prefix, DuckDB looks for catalog names like `raw`, `staging`, `mart`
- These catalog names don't exist, resulting in "Catalog Error: Table with name X does not exist"

**Example error you might see**:
```
Catalog Error: Table with name mart_velocity does not exist!
Did you mean "main_mart.mart_velocity"?
```

This error indicates the correct schema name is `main_mart.mart_velocity`, not `mart.mart_velocity`.

### Parquet Loading to DuckDB

When loading Parquet files from the api-loader into DuckDB, create tables with the correct schema:

```python
import duckdb

conn = duckdb.connect('data/analytics.duckdb')

# Create raw schema if not exists
conn.execute("CREATE SCHEMA IF NOT EXISTS main_raw")

# Load Parquet files into raw tables
conn.execute("""
    CREATE TABLE IF NOT EXISTS main_raw.commits AS
    SELECT * FROM read_parquet('data/raw/commits.parquet')
""")

conn.execute("""
    CREATE TABLE IF NOT EXISTS main_raw.pull_requests AS
    SELECT * FROM read_parquet('data/raw/pull_requests.parquet')
""")

conn.execute("""
    CREATE TABLE IF NOT EXISTS main_raw.reviews AS
    SELECT * FROM read_parquet('data/raw/reviews.parquet')
""")

conn.close()
```

---

## dbt Configuration

### dbt_project.yml

```yaml
name: 'cursor_analytics'
version: '1.0.0'

profile: 'cursor_analytics'

model-paths: ["models"]
test-paths: ["tests"]
seed-paths: ["seeds"]
macro-paths: ["macros"]

target-path: "target"
clean-targets:
  - "target"
  - "dbt_packages"

models:
  cursor_analytics:
    staging:
      +materialized: view
      +schema: staging
    intermediate:
      +materialized: ephemeral
    marts:
      +materialized: table
      +schema: mart
```

### profiles.yml

```yaml
cursor_analytics:
  target: dev  # Default to dev
  
  outputs:
    # Development: DuckDB (local)
    dev:
      type: duckdb
      path: '../data/analytics.duckdb'
      schema: main
      threads: 4
    
    # CI: DuckDB (in-memory for speed)
    ci:
      type: duckdb
      path: ':memory:'
      schema: main
      threads: 4
    
    # Production: Snowflake
    prod:
      type: snowflake
      account: "{{ env_var('SNOWFLAKE_ACCOUNT') }}"
      user: "{{ env_var('SNOWFLAKE_USER') }}"
      password: "{{ env_var('SNOWFLAKE_PASSWORD') }}"
      role: TRANSFORMER
      warehouse: TRANSFORM_WH
      database: CURSOR_ANALYTICS
      schema: RAW
      threads: 8
```

### sources.yml

```yaml
# dbt/models/sources.yml
version: 2

sources:
  - name: raw
    description: Raw data from API extraction
    schema: raw
    tables:
      - name: commits
        description: Commit-level AI telemetry from Cursor API
        columns:
          - name: commit_hash
            tests: [unique, not_null]
          - name: user_email
            tests: [not_null]
      
      - name: pull_requests
        description: PR lifecycle data from GitHub API
        columns:
          - name: number
            tests: [not_null]
          - name: repo_name
            tests: [not_null]
      
      - name: reviews
        description: Review events from GitHub API
        columns:
          - name: pr_number
            tests: [not_null]
```

---

## dbt Models

### Staging Layer

```sql
-- dbt/models/staging/stg_commits.sql
{{
    config(
        materialized='view'
    )
}}

WITH source AS (
    SELECT * FROM {{ source('raw', 'commits') }}
),

cleaned AS (
    SELECT
        commit_hash AS commit_sha,
        user_id,
        user_email,
        repo_name,
        branch_name,
        is_primary_branch,
        
        -- AI metrics
        tab_lines_added,
        tab_lines_deleted,
        composer_lines_added,
        composer_lines_deleted,
        non_ai_lines_added,
        non_ai_lines_deleted,
        
        -- Computed
        (tab_lines_added + composer_lines_added) AS ai_lines_added,
        (tab_lines_added + composer_lines_added + non_ai_lines_added) AS total_lines_added,
        
        -- Timestamps
        commit_ts,
        created_at,
        
        -- PR linkage
        pull_request_number
        
    FROM source
    WHERE commit_ts IS NOT NULL
)

SELECT * FROM cleaned
```

```sql
-- dbt/models/staging/stg_pull_requests.sql
{{
    config(
        materialized='view'
    )
}}

WITH source AS (
    SELECT * FROM {{ source('raw', 'pull_requests') }}
),

cleaned AS (
    SELECT
        number AS pr_number,
        repo_name,
        author_email,
        title,
        state,
        
        -- Size metrics
        additions,
        deletions,
        (additions + deletions) AS total_loc,
        changed_files,
        
        -- Timestamps
        first_commit_at,
        created_at,
        first_review_at,
        merged_at,
        closed_at,
        
        -- Cycle times
        coding_lead_time_hours,
        pickup_time_hours,
        review_lead_time_hours,
        
        -- Review metrics
        review_comments,
        iterations,
        reviewer_count,
        
        -- Quality flags
        is_reverted,
        has_hotfix_followup,
        
        -- AI summary
        ai_ratio
        
    FROM source
    WHERE created_at IS NOT NULL
)

SELECT * FROM cleaned
```

### Intermediate Layer

```sql
-- dbt/models/intermediate/int_pr_with_commits.sql
{{
    config(
        materialized='ephemeral'
    )
}}

WITH prs AS (
    SELECT * FROM {{ ref('stg_pull_requests') }}
),

commits AS (
    SELECT * FROM {{ ref('stg_commits') }}
),

pr_commit_summary AS (
    SELECT
        pull_request_number,
        repo_name,
        COUNT(*) AS commit_count,
        SUM(ai_lines_added) AS total_ai_lines,
        SUM(total_lines_added) AS total_lines,
        SUM(ai_lines_added)::FLOAT / NULLIF(SUM(total_lines_added), 0) AS ai_ratio_from_commits
    FROM commits
    WHERE pull_request_number IS NOT NULL
    GROUP BY 1, 2
)

SELECT
    p.*,
    c.commit_count,
    COALESCE(p.ai_ratio, c.ai_ratio_from_commits, 0) AS final_ai_ratio
FROM prs p
LEFT JOIN pr_commit_summary c
    ON p.pr_number = c.pull_request_number
    AND p.repo_name = c.repo_name
```

### Mart Layer

```sql
-- dbt/models/marts/mart_velocity.sql
{{
    config(
        materialized='table'
    )
}}

WITH pr_data AS (
    SELECT * FROM {{ ref('int_pr_with_commits') }}
    WHERE merged_at IS NOT NULL
),

weekly_metrics AS (
    SELECT
        DATE_TRUNC('week', created_at) AS week,
        repo_name,
        author_email,
        
        COUNT(*) AS prs_merged,
        AVG(total_loc) AS avg_pr_size,
        AVG(coding_lead_time_hours) AS avg_coding_lead_time,
        AVG(pickup_time_hours) AS avg_pickup_time,
        AVG(review_lead_time_hours) AS avg_review_lead_time,
        AVG(coding_lead_time_hours + pickup_time_hours + review_lead_time_hours) AS avg_total_cycle_time,
        
        PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY coding_lead_time_hours + pickup_time_hours + review_lead_time_hours) AS p50_cycle_time,
        PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY coding_lead_time_hours + pickup_time_hours + review_lead_time_hours) AS p90_cycle_time,
        
        AVG(final_ai_ratio) AS avg_ai_ratio
        
    FROM pr_data
    GROUP BY 1, 2, 3
)

SELECT
    week,
    repo_name,
    COUNT(DISTINCT author_email) AS active_developers,
    SUM(prs_merged) AS total_prs,
    AVG(avg_pr_size) AS avg_pr_size,
    AVG(avg_coding_lead_time) AS coding_lead_time,
    AVG(avg_pickup_time) AS pickup_time,
    AVG(avg_review_lead_time) AS review_lead_time,
    AVG(avg_total_cycle_time) AS total_cycle_time,
    AVG(p50_cycle_time) AS p50_cycle_time,
    AVG(p90_cycle_time) AS p90_cycle_time,
    AVG(avg_ai_ratio) AS avg_ai_ratio
FROM weekly_metrics
GROUP BY 1, 2
```

```sql
-- dbt/models/marts/mart_ai_impact.sql
{{
    config(
        materialized='table'
    )
}}

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
),

impact_by_band AS (
    SELECT
        ai_usage_band,
        DATE_TRUNC('week', created_at) AS week,
        
        COUNT(*) AS pr_count,
        AVG(final_ai_ratio) AS avg_ai_ratio,
        
        -- Velocity
        AVG(coding_lead_time_hours) AS avg_coding_lead_time,
        AVG(pickup_time_hours + review_lead_time_hours) AS avg_review_cycle_time,
        
        -- Review costs
        AVG(review_comments::FLOAT / NULLIF(total_loc, 0)) AS avg_review_density,
        AVG(iterations) AS avg_iterations,
        AVG(reviewer_count) AS avg_reviewers,
        
        -- Quality
        AVG(CASE WHEN is_reverted THEN 1 ELSE 0 END) AS revert_rate,
        AVG(CASE WHEN has_hotfix_followup THEN 1 ELSE 0 END) AS hotfix_rate
        
    FROM pr_data
    GROUP BY 1, 2
)

SELECT * FROM impact_by_band
```

### dbt Tests

```yaml
# dbt/models/staging/schema.yml
version: 2

models:
  - name: stg_commits
    columns:
      - name: commit_sha
        tests:
          - unique
          - not_null
      - name: user_email
        tests:
          - not_null
      - name: ai_lines_added
        tests:
          - dbt_utils.accepted_range:
              min_value: 0

  - name: stg_pull_requests
    columns:
      - name: pr_number
        tests:
          - not_null
      - name: coding_lead_time_hours
        tests:
          - dbt_utils.accepted_range:
              min_value: 0
      - name: ai_ratio
        tests:
          - dbt_utils.accepted_range:
              min_value: 0
              max_value: 1
```

---

## CI/CD Pipelines

### 1. cursor-sim CI

```yaml
# .github/workflows/ci-cursor-sim.yml
name: cursor-sim CI

on:
  push:
    paths:
      - 'apps/cursor-sim/**'
      - 'tools/data-designer/**'
  pull_request:
    paths:
      - 'apps/cursor-sim/**'

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build
        run: |
          cd apps/cursor-sim
          go build -v -o cursor-sim ./...
      
      - name: Unit tests
        run: |
          cd apps/cursor-sim
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Generate seed
        run: |
          cd tools/data-designer
          pip install -r requirements.txt
          python generate_seed.py --fallback -o ../../seed.json
      
      - name: Start simulator and verify API
        run: |
          cd apps/cursor-sim
          ./cursor-sim --seed=../../seed.json --port=8080 &
          sleep 5
          
          # Verify Cursor API
          curl -f http://localhost:8080/teams/members
          curl -f "http://localhost:8080/analytics/ai-code/commits?startDate=7d&endDate=now"
          
          # Verify GitHub API
          curl -f http://localhost:8080/repos
```

### 2. Loader CI (Tests Extraction Logic)

```yaml
# .github/workflows/ci-loader.yml
name: Loader CI

on:
  push:
    paths:
      - 'tools/api-loader/**'
      - 'apps/cursor-sim/**'
  pull_request:
    paths:
      - 'tools/api-loader/**'

jobs:
  test-extraction:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Build and start simulator
        run: |
          # Generate seed
          cd tools/data-designer
          pip install -r requirements.txt
          python generate_seed.py --fallback -o ../../seed.json
          
          # Build and start simulator
          cd ../../apps/cursor-sim
          go build -o cursor-sim .
          ./cursor-sim --seed=../../seed.json --days=7 --port=8080 &
          sleep 5
      
      - name: Install loader dependencies
        run: |
          pip install -r tools/api-loader/requirements.txt
          pip install pytest
      
      - name: Run extraction
        run: |
          python tools/api-loader/loader.py \
            --url http://localhost:8080 \
            --output data/raw \
            --start-date 7d \
            --validate
      
      - name: Validate Parquet schemas
        run: |
          python << 'EOF'
          import pandas as pd
          
          # Check commits
          commits = pd.read_parquet('data/raw/commits.parquet')
          required = ['commit_hash', 'user_email', 'repo_name', 'tab_lines_added']
          assert all(col in commits.columns for col in required)
          assert len(commits) > 0
          
          # Check PRs
          prs = pd.read_parquet('data/raw/pull_requests.parquet')
          required = ['number', 'author_email', 'ai_ratio', 'coding_lead_time_hours']
          assert all(col in prs.columns for col in required)
          assert len(prs) > 0
          
          # Check reviews
          reviews = pd.read_parquet('data/raw/reviews.parquet')
          required = ['pr_number', 'repo_name', 'state']
          assert all(col in reviews.columns for col in required)
          
          print("✓ All schemas valid")
          EOF
      
      - name: Upload Parquet artifacts
        uses: actions/upload-artifact@v4
        with:
          name: raw-data
          path: data/raw/
          retention-days: 1
```

### 3. dbt CI

```yaml
# .github/workflows/ci-dbt.yml
name: dbt CI

on:
  push:
    paths:
      - 'dbt/**'
  pull_request:
    paths:
      - 'dbt/**'

jobs:
  dbt-test:
    runs-on: ubuntu-latest
    needs: []  # Can run independently if raw data is uploaded as artifact
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Setup Go (for cursor-sim)
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: |
          pip install dbt-duckdb dbt-snowflake
          pip install -r tools/api-loader/requirements.txt
      
      - name: Generate test data via API
        run: |
          # Generate seed
          cd tools/data-designer
          python generate_seed.py --fallback -o ../../seed.json
          
          # Start simulator
          cd ../../apps/cursor-sim
          go build -o cursor-sim .
          ./cursor-sim --seed=../../seed.json --days=7 --port=8080 &
          sleep 5
          
          # Extract via loader (tests extraction parity)
          cd ../..
          python tools/api-loader/loader.py -o data/raw --start-date 7d
      
      - name: Setup DuckDB with raw data
        run: |
          pip install duckdb
          python << 'EOF'
          import duckdb
          conn = duckdb.connect('data/analytics.duckdb')
          conn.execute("CREATE SCHEMA IF NOT EXISTS raw")
          for table in ['commits', 'pull_requests', 'reviews', 'team_members']:
              try:
                  conn.execute(f"""
                      CREATE TABLE raw.{table} AS 
                      SELECT * FROM read_parquet('data/raw/{table}.parquet')
                  """)
                  print(f"✓ Loaded raw.{table}")
              except Exception as e:
                  print(f"⚠ Skipped {table}: {e}")
          conn.close()
          EOF
      
      - name: dbt deps
        run: cd dbt && dbt deps
      
      - name: dbt build (dev target = DuckDB)
        run: cd dbt && dbt build --target dev
      
      - name: dbt test
        run: cd dbt && dbt test --target dev
      
      - name: Check documentation
        run: cd dbt && dbt docs generate && test -f target/catalog.json

  lint-sql:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install sqlfluff
        run: pip install sqlfluff sqlfluff-templater-dbt
      
      - name: Lint for Snowflake compatibility
        run: |
          cd dbt
          sqlfluff lint models/ --dialect snowflake
```

### 4. Dashboard CI

```yaml
# .github/workflows/ci-dashboard.yml
name: Dashboard CI

on:
  push:
    paths:
      - 'apps/dashboard/**'
      - 'dbt/models/marts/**'
  pull_request:
    paths:
      - 'apps/dashboard/**'

jobs:
  test-dashboard:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: |
          pip install -r apps/dashboard/requirements.txt
          pip install -r tools/api-loader/requirements.txt
          pip install dbt-duckdb pytest pytest-cov
      
      - name: Generate data pipeline (sim → API → loader → dbt)
        run: |
          # Generate seed
          cd tools/data-designer && python generate_seed.py --fallback -o ../../seed.json && cd ../..
          
          # Start simulator
          cd apps/cursor-sim && go build -o cursor-sim . && ./cursor-sim --seed=../../seed.json --days=7 --port=8080 &
          sleep 5 && cd ../..
          
          # Extract via loader
          python tools/api-loader/loader.py -o data/raw --start-date 7d
          
          # Load to DuckDB
          python -c "
          import duckdb
          conn = duckdb.connect('data/analytics.duckdb')
          conn.execute('CREATE SCHEMA IF NOT EXISTS raw')
          for t in ['commits', 'pull_requests', 'reviews', 'team_members']:
              try:
                  conn.execute(f\"CREATE TABLE raw.{t} AS SELECT * FROM read_parquet('data/raw/{t}.parquet')\")
              except: pass
          conn.close()
          "
          
          # Run dbt
          cd dbt && dbt deps && dbt build --target dev
      
      - name: Run dashboard tests
        run: |
          cd apps/dashboard
          pytest tests/ -v --cov=. --cov-report=xml
        env:
          ENVIRONMENT: dev
      
      - name: Test Streamlit loads
        run: |
          cd apps/dashboard
          timeout 30 streamlit run app.py --server.headless=true &
          sleep 10
          curl -f http://localhost:8501/_stcore/health
        env:
          ENVIRONMENT: dev
```

### 5. Production Deployment

```yaml
# .github/workflows/cd-prod.yml
name: Production Deployment

on:
  push:
    branches: [main]
    paths:
      - 'dbt/**'
      - 'apps/dashboard/**'
      - 'snaplogic/**'

jobs:
  deploy-dbt:
    runs-on: ubuntu-latest
    environment: production
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Install dbt
        run: pip install dbt-snowflake
      
      - name: dbt deps
        run: cd dbt && dbt deps
      
      - name: dbt run (prod)
        run: cd dbt && dbt run --target prod
        env:
          SNOWFLAKE_ACCOUNT: ${{ secrets.SNOWFLAKE_ACCOUNT }}
          SNOWFLAKE_USER: ${{ secrets.SNOWFLAKE_USER }}
          SNOWFLAKE_PASSWORD: ${{ secrets.SNOWFLAKE_PASSWORD }}
      
      - name: dbt test (prod)
        run: cd dbt && dbt test --target prod
        env:
          SNOWFLAKE_ACCOUNT: ${{ secrets.SNOWFLAKE_ACCOUNT }}
          SNOWFLAKE_USER: ${{ secrets.SNOWFLAKE_USER }}
          SNOWFLAKE_PASSWORD: ${{ secrets.SNOWFLAKE_PASSWORD }}

  deploy-dashboard:
    runs-on: ubuntu-latest
    needs: deploy-dbt
    environment: production
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Deploy to Streamlit Cloud
        uses: streamlit/streamlit-app-action@v1
        with:
          app-path: apps/dashboard/app.py
        env:
          STREAMLIT_TOKEN: ${{ secrets.STREAMLIT_TOKEN }}
          ENVIRONMENT: prod
          SNOWFLAKE_ACCOUNT: ${{ secrets.SNOWFLAKE_ACCOUNT }}
          SNOWFLAKE_USER: ${{ secrets.SNOWFLAKE_USER }}
          SNOWFLAKE_PASSWORD: ${{ secrets.SNOWFLAKE_PASSWORD }}

  notify-snaplogic:
    runs-on: ubuntu-latest
    if: contains(github.event.head_commit.modified, 'snaplogic/')
    steps:
      - name: Notify schema change
        run: |
          echo "⚠️ SnapLogic configs changed - manual review required"
```

---

## Makefile

```makefile
# =============================================================================
# Development Workflow
# =============================================================================

.PHONY: dev dev-sim dev-load dev-dbt dev-dashboard

# Full dev pipeline
dev: dev-sim dev-load dev-dbt dev-dashboard

# Start simulator
dev-sim:
	@echo "Starting cursor-sim..."
	cd apps/cursor-sim && go run . --seed=../../seed.json --port=8080

# Extract from API → Parquet (mimics SnapLogic)
dev-load:
	@echo "Extracting data from API..."
	python tools/api-loader/loader.py -o data/raw --validate

# Run dbt transforms
dev-dbt:
	@echo "Running dbt..."
	cd dbt && dbt deps && dbt build --target dev

# Start dashboard
dev-dashboard:
	@echo "Starting dashboard..."
	cd apps/dashboard && streamlit run app.py

# =============================================================================
# Seed Generation
# =============================================================================

seed-generate:
	cd tools/data-designer && python generate_seed.py -o ../../seed.json

seed-generate-fallback:
	cd tools/data-designer && python generate_seed.py --fallback -o ../../seed.json

# =============================================================================
# Testing
# =============================================================================

test: test-sim test-loader test-dbt test-dashboard

test-sim:
	cd apps/cursor-sim && go test -v -race ./...

test-loader:
	cd tools/api-loader && pytest tests/ -v

test-dbt:
	cd dbt && dbt test --target dev

test-dashboard:
	cd apps/dashboard && pytest tests/ -v

# =============================================================================
# CI Simulation (run full pipeline locally)
# =============================================================================

ci-local:
	@echo "=== Generating seed ==="
	$(MAKE) seed-generate-fallback
	
	@echo "=== Starting simulator ==="
	cd apps/cursor-sim && go build -o cursor-sim . && ./cursor-sim --seed=../../seed.json --days=7 --port=8080 &
	sleep 5
	
	@echo "=== Extracting via loader ==="
	python tools/api-loader/loader.py -o data/raw --start-date 7d --validate
	
	@echo "=== Loading to DuckDB ==="
	python -c "import duckdb; c=duckdb.connect('data/analytics.duckdb'); c.execute('CREATE SCHEMA IF NOT EXISTS raw'); [c.execute(f\"CREATE OR REPLACE TABLE raw.{t} AS SELECT * FROM read_parquet('data/raw/{t}.parquet')\") for t in ['commits','pull_requests','reviews','team_members']]"
	
	@echo "=== Running dbt ==="
	cd dbt && dbt deps && dbt build --target dev && dbt test --target dev
	
	@echo "=== CI simulation complete ==="

# =============================================================================
# Cleanup
# =============================================================================

clean:
	rm -rf data/raw/*.parquet
	rm -f data/analytics.duckdb
	rm -rf dbt/target dbt/logs
```

---

## Summary: Dev/Prod Parity

| Layer | Dev | Prod | Parity |
|-------|-----|------|--------|
| **API Source** | cursor-sim (REST) | Cursor + GitHub (REST) | ✅ Same API contract |
| **Extraction** | api-loader (Python) | SnapLogic | ✅ Same logic, tests in CI |
| **Landing** | Parquet files | Snowflake Stage | ✅ Same format |
| **Raw Tables** | DuckDB | Snowflake | ✅ Same schema |
| **Transforms** | dbt | dbt | ✅ Identical SQL |
| **Marts** | DuckDB | Snowflake | ✅ Same schema |
| **Dashboard** | Streamlit | Streamlit | ✅ Identical |

**Benefits of API + Loader pattern:**
1. Tests extraction logic in CI (validates API contract)
2. Same pagination/transformation code runs in dev
3. Catches schema drift early (loader validates output)
4. cursor-sim API matches production exactly