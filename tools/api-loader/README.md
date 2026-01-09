# API Loader

Extract data from cursor-sim REST API into Parquet files for downstream processing.

## Installation

```bash
cd tools/api-loader
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
```

## Usage

```bash
# Extract data from cursor-sim
python loader.py --url http://localhost:8080 --output ../../data/raw
```

## Architecture

- `extractors/cursor_api.py` - Cursor API endpoints (/analytics/*)
- `extractors/github_api.py` - GitHub-style endpoints (/repos/*)
- `schemas/` - JSON schema validation files
- `loader.py` - Main orchestration script
- `duckdb_loader.py` - Load Parquet to DuckDB

## Testing

```bash
pytest tests/
```
