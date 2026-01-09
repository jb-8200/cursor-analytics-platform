"""
DuckDB loader for loading Parquet files into DuckDB database.

Loads Parquet files from a directory into DuckDB raw schema tables.
Supports both full refresh and incremental loading modes.
"""

import argparse
import logging
from pathlib import Path

import duckdb

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def load_parquet_to_duckdb(
    parquet_dir: Path,
    db_path: Path,
    incremental: bool = False
) -> None:
    """
    Load Parquet files from a directory into DuckDB.

    Creates tables in the 'raw' schema with names matching the Parquet file names.

    Args:
        parquet_dir: Directory containing Parquet files
        db_path: Path to DuckDB database file
        incremental: If True, append to existing tables; if False, replace tables
    """
    parquet_dir = Path(parquet_dir)
    db_path = Path(db_path)

    # Create database directory if needed
    db_path.parent.mkdir(parents=True, exist_ok=True)

    # Connect to DuckDB
    conn = duckdb.connect(str(db_path))

    try:
        # Create raw schema if it doesn't exist
        conn.execute("CREATE SCHEMA IF NOT EXISTS raw")
        logger.info("Ensured raw schema exists")

        # Find all Parquet files
        parquet_files = list(parquet_dir.glob("*.parquet"))

        if not parquet_files:
            logger.warning(f"No Parquet files found in {parquet_dir}")
            return

        logger.info(f"Found {len(parquet_files)} Parquet files to load")

        for parquet_file in parquet_files:
            table_name = parquet_file.stem  # File name without extension

            if incremental:
                # Check if table exists
                table_exists = conn.execute(f"""
                    SELECT COUNT(*) FROM information_schema.tables
                    WHERE table_schema = 'raw' AND table_name = '{table_name}'
                """).fetchone()[0] > 0

                if table_exists:
                    # Append to existing table
                    conn.execute(f"""
                        INSERT INTO raw.{table_name}
                        SELECT * FROM read_parquet('{parquet_file}')
                    """)
                    logger.info(f"Appended to raw.{table_name}")
                else:
                    # Create new table
                    conn.execute(f"""
                        CREATE TABLE raw.{table_name} AS
                        SELECT * FROM read_parquet('{parquet_file}')
                    """)
                    logger.info(f"Created raw.{table_name}")
            else:
                # Full refresh: replace table
                conn.execute(f"""
                    CREATE OR REPLACE TABLE raw.{table_name} AS
                    SELECT * FROM read_parquet('{parquet_file}')
                """)
                logger.info(f"Loaded raw.{table_name}")

        logger.info(f"Done! Loaded {len(parquet_files)} tables to {db_path}")

    finally:
        conn.close()


def create_parser() -> argparse.ArgumentParser:
    """
    Create argument parser for CLI.

    Returns:
        Configured ArgumentParser
    """
    parser = argparse.ArgumentParser(
        description="Load Parquet files into DuckDB database"
    )
    parser.add_argument(
        "--parquet-dir",
        default="data/raw",
        help="Directory containing Parquet files (default: data/raw)"
    )
    parser.add_argument(
        "--db-path",
        default="data/analytics.duckdb",
        help="Path to DuckDB database (default: data/analytics.duckdb)"
    )
    parser.add_argument(
        "--incremental",
        action="store_true",
        help="Append to existing tables instead of replacing"
    )
    return parser


def main() -> None:
    """Main entry point for CLI."""
    parser = create_parser()
    args = parser.parse_args()

    load_parquet_to_duckdb(
        Path(args.parquet_dir),
        Path(args.db_path),
        incremental=args.incremental
    )


if __name__ == "__main__":
    main()
