#!/bin/bash
#
# Test script for pipeline orchestration
#
# Tests that the pipeline script exists, is executable, and has correct structure.

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

PASS_COUNT=0
FAIL_COUNT=0

# Test helper
run_test() {
    local test_name="$1"
    local test_command="$2"

    echo -n "Testing: $test_name... "
    if eval "$test_command"; then
        echo -e "${GREEN}PASS${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}FAIL${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

echo "Running pipeline script tests..."
echo ""

# Test 1: Pipeline script exists
run_test "Pipeline script exists" \
    "test -f tools/run_pipeline.sh"

# Test 2: Pipeline script is executable
run_test "Pipeline script is executable" \
    "test -x tools/run_pipeline.sh"

# Test 3: Script has shebang
run_test "Script has bash shebang" \
    "head -n 1 tools/run_pipeline.sh | grep -q '#!/bin/bash'"

# Test 4: Script uses set -e
run_test "Script uses set -e" \
    "grep -q 'set -e' tools/run_pipeline.sh"

# Test 5: Script has all required steps
run_test "Script has extract step" \
    "grep -q 'python tools/api-loader/loader.py' tools/run_pipeline.sh"

run_test "Script has load step" \
    "grep -q 'python tools/api-loader/duckdb_loader.py' tools/run_pipeline.sh"

run_test "Script has dbt step" \
    "grep -q 'dbt build' tools/run_pipeline.sh"

# Test 6: Script uses environment variables
run_test "Script uses CURSOR_SIM_URL" \
    "grep -q 'CURSOR_SIM_URL' tools/run_pipeline.sh"

run_test "Script uses DATA_DIR" \
    "grep -q 'DATA_DIR' tools/run_pipeline.sh"

# Test 7: Script has error handling
run_test "Script has error handler" \
    "grep -q 'handle_error' tools/run_pipeline.sh"

# Test 8: Makefile has pipeline targets
run_test "Makefile has pipeline target" \
    "grep -q '^pipeline:' Makefile"

run_test "Makefile has extract target" \
    "grep -q '^extract:' Makefile"

run_test "Makefile has load target" \
    "grep -q '^load:' Makefile"

run_test "Makefile has dbt-build target" \
    "grep -q '^dbt-build:' Makefile"

run_test "Makefile has ci-local target" \
    "grep -q '^ci-local:' Makefile"

# Test 9: Validate script syntax
run_test "Script has valid bash syntax" \
    "bash -n tools/run_pipeline.sh"

echo ""
echo "========================================="
echo "Test Results:"
echo "  PASSED: $PASS_COUNT"
echo "  FAILED: $FAIL_COUNT"
echo "========================================="

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
