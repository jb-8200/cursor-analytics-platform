#!/bin/bash
#
# ETL Pipeline Orchestrator
#
# Executes the full data pipeline from cursor-sim API to DuckDB analytics tables.
#
# Usage:
#   ./tools/run_pipeline.sh
#
# Environment Variables:
#   CURSOR_SIM_URL  - Base URL for cursor-sim API (default: http://localhost:8080)
#   DATA_DIR        - Base directory for data files (default: ./data)
#   API_KEY         - API key for cursor-sim (default: cursor-sim-dev-key)
#   START_DATE      - Start date for commits filter (default: 90d)
#   CONTINUE_ON_ERROR - Continue on step failures (default: false)

set -e  # Exit on error (unless CONTINUE_ON_ERROR is set)

# Configuration with defaults
CURSOR_SIM_URL=${CURSOR_SIM_URL:-"http://localhost:8080"}
DATA_DIR=${DATA_DIR:-"./data"}
API_KEY=${API_KEY:-"cursor-sim-dev-key"}
START_DATE=${START_DATE:-"90d"}
CONTINUE_ON_ERROR=${CONTINUE_ON_ERROR:-false}

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Logging helpers
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Error handler
handle_error() {
    log_error "Pipeline failed at step: $1"
    if [ "$CONTINUE_ON_ERROR" = "true" ]; then
        log_warn "Continuing despite error (CONTINUE_ON_ERROR=true)"
        return 0
    else
        exit 1
    fi
}

# Pipeline start
log_info "Starting ETL Pipeline"
log_info "Configuration:"
log_info "  cursor-sim URL: $CURSOR_SIM_URL"
log_info "  Data directory: $DATA_DIR"
log_info "  Start date: $START_DATE"
log_info ""

# Step 1: Extract data from cursor-sim API
log_info "=== Step 1/3: Extract from cursor-sim API ==="
if [ "$CONTINUE_ON_ERROR" = "true" ]; then
    python tools/api-loader/loader.py \
        --url "$CURSOR_SIM_URL" \
        --output "$DATA_DIR/raw" \
        --api-key "$API_KEY" \
        --start-date "$START_DATE" \
        --continue-on-error || handle_error "Extract"
else
    python tools/api-loader/loader.py \
        --url "$CURSOR_SIM_URL" \
        --output "$DATA_DIR/raw" \
        --api-key "$API_KEY" \
        --start-date "$START_DATE" || handle_error "Extract"
fi

log_info "Extraction complete"
log_info ""

# Step 2: Load Parquet files to DuckDB
log_info "=== Step 2/3: Load to DuckDB ==="
python tools/api-loader/duckdb_loader.py \
    --parquet-dir "$DATA_DIR/raw" \
    --db-path "$DATA_DIR/analytics.duckdb" || handle_error "Load"

log_info "DuckDB loading complete"
log_info ""

# Step 3: Run dbt transformations
log_info "=== Step 3/3: Run dbt transformations ==="
cd dbt || handle_error "dbt directory not found"
dbt deps || handle_error "dbt deps"
dbt build --target dev || handle_error "dbt build"
cd ..

log_info "dbt transformations complete"
log_info ""

# Pipeline complete
log_info "================================================"
log_info "Pipeline complete!"
log_info ""
log_info "Output files:"
log_info "  Raw Parquet: $DATA_DIR/raw/*.parquet"
log_info "  DuckDB: $DATA_DIR/analytics.duckdb"
log_info "  dbt artifacts: dbt/target/"
log_info ""
log_info "Query analytics tables:"
log_info "  duckdb $DATA_DIR/analytics.duckdb"
log_info "  > SELECT * FROM mart.velocity LIMIT 10;"
log_info "================================================"
