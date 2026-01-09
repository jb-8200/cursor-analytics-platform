"""
Pytest configuration and fixtures for Streamlit dashboard tests.

Feature: P9-F01 Streamlit Dashboard
"""

import pytest
import os
import tempfile
from pathlib import Path


@pytest.fixture
def test_db_path(tmp_path):
    """Create a temporary DuckDB database path."""
    db_path = tmp_path / "test.duckdb"
    return str(db_path)


@pytest.fixture
def setup_duckdb_env(monkeypatch, test_db_path):
    """Setup DuckDB environment variables for testing."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", test_db_path)
    yield test_db_path


@pytest.fixture
def sample_data_db(setup_duckdb_env):
    """Create a DuckDB with sample data for testing."""
    import duckdb

    conn = duckdb.connect(setup_duckdb_env)

    # Create sample mart tables
    conn.execute("""
        CREATE TABLE IF NOT EXISTS mart.velocity (
            week DATE,
            repo_name VARCHAR,
            active_developers INTEGER,
            total_prs INTEGER,
            avg_pr_size DOUBLE,
            coding_lead_time DOUBLE,
            pickup_time DOUBLE,
            review_lead_time DOUBLE,
            total_cycle_time DOUBLE,
            p50_cycle_time DOUBLE,
            p90_cycle_time DOUBLE,
            avg_ai_ratio DOUBLE
        )
    """)

    conn.execute("""
        INSERT INTO mart.velocity VALUES
        ('2026-01-01', 'acme/platform', 10, 50, 150.0, 2.0, 0.5, 1.0, 3.5, 2.5, 5.0, 0.45),
        ('2026-01-08', 'acme/platform', 12, 60, 180.0, 1.8, 0.4, 0.9, 3.1, 2.3, 4.8, 0.50)
    """)

    conn.close()
    yield setup_duckdb_env
