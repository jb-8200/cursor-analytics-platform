"""
DuckDB loader for parquet files.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-10 Refresh Pipeline

This module provides utilities to load parquet files into DuckDB.
"""

import duckdb
from pathlib import Path
from typing import List


def load_parquet_to_duckdb(raw_dir: str, duckdb_path: str) -> None:
    """
    Load parquet files from raw directory into DuckDB.

    Args:
        raw_dir: Path to directory containing parquet files
        duckdb_path: Path to DuckDB database file

    Raises:
        Exception: If loading fails
    """
    raw_path = Path(raw_dir)
    if not raw_path.exists():
        raise FileNotFoundError(f"Raw directory not found: {raw_dir}")

    # Connect to DuckDB
    conn = duckdb.connect(duckdb_path, read_only=False)

    try:
        # Create raw schema if it doesn't exist
        conn.execute("CREATE SCHEMA IF NOT EXISTS raw")

        # Find all parquet files
        parquet_files = list(raw_path.glob("*.parquet"))

        if not parquet_files:
            print(f"No parquet files found in {raw_dir}")
            return

        # Load each parquet file into a table
        for parquet_file in parquet_files:
            table_name = parquet_file.stem  # filename without extension
            full_table_name = f"raw.{table_name}"

            print(f"Loading {parquet_file.name} -> {full_table_name}")

            # Drop table if exists and recreate from parquet
            conn.execute(f"DROP TABLE IF EXISTS {full_table_name}")
            conn.execute(f"""
                CREATE TABLE {full_table_name} AS
                SELECT * FROM read_parquet('{parquet_file}')
            """)

            # Verify data loaded
            count = conn.execute(f"SELECT COUNT(*) FROM {full_table_name}").fetchone()[0]
            print(f"  Loaded {count:,} rows")

    finally:
        conn.close()


def get_parquet_files(raw_dir: str) -> List[str]:
    """
    Get list of parquet files in raw directory.

    Args:
        raw_dir: Path to directory containing parquet files

    Returns:
        List of parquet file paths
    """
    raw_path = Path(raw_dir)
    if not raw_path.exists():
        return []

    return [str(f) for f in raw_path.glob("*.parquet")]
