"""
Database connector for DuckDB and Snowflake.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-03 Database Connector

This module provides a unified interface for connecting to either:
- DuckDB (local development)
- Snowflake (production)

Configuration via environment variables:
- DB_MODE: "duckdb" or "snowflake" (default: "duckdb")
- DUCKDB_PATH: Path to DuckDB file (default: "/data/analytics.duckdb")
- SNOWFLAKE_*: Snowflake connection credentials
"""

import os
from typing import Optional
import pandas as pd

# Determine database mode
DB_MODE = os.getenv("DB_MODE", "duckdb")


def get_connection():
    """
    Get database connection based on DB_MODE environment variable.

    Returns:
        Database connection object (DuckDB or Snowflake)
    """
    if DB_MODE == "snowflake":
        return _get_snowflake_connection()
    else:
        return _get_duckdb_connection()


def _get_duckdb_connection():
    """
    Create a DuckDB connection.

    Returns:
        DuckDB connection
    """
    import duckdb

    db_path = os.getenv("DUCKDB_PATH", "/data/analytics.duckdb")
    return duckdb.connect(db_path, read_only=False)


def _get_snowflake_connection():
    """
    Create a Snowflake connection.

    Returns:
        Snowflake connection

    Raises:
        Exception: If connection fails or credentials are missing
    """
    import snowflake.connector

    return snowflake.connector.connect(
        account=os.getenv("SNOWFLAKE_ACCOUNT"),
        user=os.getenv("SNOWFLAKE_USER"),
        password=os.getenv("SNOWFLAKE_PASSWORD"),
        database=os.getenv("SNOWFLAKE_DATABASE", "CURSOR_ANALYTICS"),
        schema=os.getenv("SNOWFLAKE_SCHEMA", "MART"),
        warehouse=os.getenv("SNOWFLAKE_WAREHOUSE", "TRANSFORM_WH"),
    )


def query(sql: str, params: Optional[dict] = None) -> pd.DataFrame:
    """
    Execute SQL query and return results as pandas DataFrame.

    Args:
        sql: SQL query string
        params: Optional query parameters (dict)

    Returns:
        pandas DataFrame with query results

    Note:
        In Streamlit context, this would be decorated with @st.cache_data(ttl=300)
        For testing purposes, we keep it undecorated.
    """
    conn = get_connection()

    if DB_MODE == "snowflake":
        cursor = conn.cursor()
        cursor.execute(sql, params or {})
        columns = [desc[0] for desc in cursor.description]
        data = cursor.fetchall()
        cursor.close()
        return pd.DataFrame(data, columns=columns)
    else:
        # DuckDB
        if params:
            # DuckDB parameterized query
            return conn.execute(sql, list(params.values())).df()
        return conn.execute(sql).df()


def refresh_data():
    """
    Trigger ETL pipeline to refresh data (dev mode only).

    In production (Snowflake), data is refreshed via scheduled jobs.
    In development (DuckDB), this triggers the ETL pipeline:
    1. Extract from cursor-sim
    2. Load to DuckDB
    3. Run dbt transformations

    Returns:
        bool: True if refresh succeeded, False otherwise

    Note:
        In Streamlit context, this would clear caches after refresh.
    """
    if DB_MODE == "snowflake":
        print("Refresh not available in production mode.")
        return False

    import subprocess

    cursor_sim_url = os.getenv("CURSOR_SIM_URL", "http://localhost:8080")
    duckdb_path = os.getenv("DUCKDB_PATH", "/data/analytics.duckdb")

    try:
        # Step 1: Extract data from cursor-sim
        print("Extracting data from cursor-sim...")
        result = subprocess.run(
            [
                "python",
                "tools/api-loader/loader.py",
                "--url",
                cursor_sim_url,
                "--output",
                "/data/raw",
            ],
            capture_output=True,
            text=True,
            check=True,
        )

        # Step 2: Load to DuckDB
        print("Loading data to DuckDB...")
        # Import dynamically to avoid import errors in test environment
        try:
            from pipeline.duckdb_loader import load_parquet_to_duckdb
            load_parquet_to_duckdb("/data/raw", duckdb_path)
        except ImportError:
            # Fallback for testing environment
            import sys
            sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))
            from pipeline.duckdb_loader import load_parquet_to_duckdb
            load_parquet_to_duckdb("/data/raw", duckdb_path)

        # Step 3: Run dbt
        print("Running dbt transformations...")
        result = subprocess.run(
            ["dbt", "build", "--target", "dev"],
            cwd="/app/dbt",
            capture_output=True,
            text=True,
            check=True,
        )

        print("Data refresh completed successfully!")
        return True

    except subprocess.CalledProcessError as e:
        print(f"Refresh failed: {e.stderr}")
        return False
    except Exception as e:
        print(f"Refresh failed: {str(e)}")
        return False
