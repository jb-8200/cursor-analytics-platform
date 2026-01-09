"""
Tests for DuckDB loader.

Tests loading Parquet files into DuckDB database with raw schema.
"""

import pytest
import pandas as pd
import duckdb
from pathlib import Path


@pytest.fixture
def sample_parquet_files(tmp_path):
    """Create sample Parquet files for testing."""
    raw_dir = tmp_path / "raw"
    raw_dir.mkdir()

    # Create repos.parquet
    repos_df = pd.DataFrame([
        {"full_name": "acme/platform", "default_branch": "main"},
        {"full_name": "acme/frontend", "default_branch": "main"}
    ])
    repos_df.to_parquet(raw_dir / "repos.parquet", index=False)

    # Create commits.parquet
    commits_df = pd.DataFrame([
        {"commitHash": "abc123", "repoName": "acme/platform", "tabLinesAdded": 50},
        {"commitHash": "def456", "repoName": "acme/frontend", "tabLinesAdded": 30}
    ])
    commits_df.to_parquet(raw_dir / "commits.parquet", index=False)

    # Create pull_requests.parquet
    prs_df = pd.DataFrame([
        {"number": 1, "repo_name": "acme/platform", "state": "merged"},
        {"number": 2, "repo_name": "acme/platform", "state": "open"}
    ])
    prs_df.to_parquet(raw_dir / "pull_requests.parquet", index=False)

    # Create reviews.parquet
    reviews_df = pd.DataFrame([
        {"id": 101, "repo_name": "acme/platform", "pr_number": 1, "state": "APPROVED"},
        {"id": 102, "repo_name": "acme/platform", "pr_number": 1, "state": "CHANGES_REQUESTED"}
    ])
    reviews_df.to_parquet(raw_dir / "reviews.parquet", index=False)

    return raw_dir


class TestLoadParquetToDuckDB:
    """Tests for load_parquet_to_duckdb function."""

    def test_load_parquet_to_duckdb_creates_tables(self, sample_parquet_files, tmp_path):
        """Verify Parquet files are loaded as DuckDB tables."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        # Verify tables were created
        conn = duckdb.connect(str(db_path))

        # Check repos table
        result = conn.execute("SELECT COUNT(*) FROM raw.repos").fetchone()
        assert result[0] == 2

        # Check commits table
        result = conn.execute("SELECT COUNT(*) FROM raw.commits").fetchone()
        assert result[0] == 2

        # Check pull_requests table
        result = conn.execute("SELECT COUNT(*) FROM raw.pull_requests").fetchone()
        assert result[0] == 2

        # Check reviews table
        result = conn.execute("SELECT COUNT(*) FROM raw.reviews").fetchone()
        assert result[0] == 2

        conn.close()

    def test_load_parquet_creates_raw_schema(self, sample_parquet_files, tmp_path):
        """Verify raw schema is created."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        conn = duckdb.connect(str(db_path))

        # Check schema exists
        result = conn.execute("""
            SELECT schema_name FROM information_schema.schemata
            WHERE schema_name = 'raw'
        """).fetchone()
        assert result is not None
        assert result[0] == "raw"

        conn.close()

    def test_load_parquet_preserves_data(self, sample_parquet_files, tmp_path):
        """Verify data is preserved correctly."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        conn = duckdb.connect(str(db_path))

        # Check repos data
        result = conn.execute("SELECT full_name FROM raw.repos ORDER BY full_name").fetchall()
        assert [r[0] for r in result] == ["acme/frontend", "acme/platform"]

        # Check commits data
        result = conn.execute("SELECT commitHash FROM raw.commits ORDER BY commitHash").fetchall()
        assert [r[0] for r in result] == ["abc123", "def456"]

        conn.close()

    def test_load_parquet_idempotent(self, sample_parquet_files, tmp_path):
        """Verify loader is idempotent (CREATE OR REPLACE)."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"

        # Load once
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        # Load again
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        # Should still have same data (not duplicated)
        conn = duckdb.connect(str(db_path))
        result = conn.execute("SELECT COUNT(*) FROM raw.repos").fetchone()
        assert result[0] == 2
        conn.close()

    def test_load_parquet_empty_directory(self, tmp_path):
        """Handle empty parquet directory gracefully."""
        from duckdb_loader import load_parquet_to_duckdb

        empty_dir = tmp_path / "empty"
        empty_dir.mkdir()
        db_path = tmp_path / "test.duckdb"

        # Should not raise
        load_parquet_to_duckdb(empty_dir, db_path)

        # Database should exist with raw schema
        conn = duckdb.connect(str(db_path))
        result = conn.execute("""
            SELECT schema_name FROM information_schema.schemata
            WHERE schema_name = 'raw'
        """).fetchone()
        assert result is not None
        conn.close()

    def test_load_parquet_creates_db_directory(self, sample_parquet_files, tmp_path):
        """Verify loader creates database directory if needed."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "nested" / "dir" / "test.duckdb"
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        assert db_path.exists()


class TestDuckDBLoaderIncremental:
    """Tests for incremental loading."""

    def test_full_refresh_replaces_data(self, sample_parquet_files, tmp_path):
        """Verify full refresh replaces existing data."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"

        # Initial load
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        # Modify source data
        new_repos_df = pd.DataFrame([
            {"full_name": "acme/new-repo", "default_branch": "main"}
        ])
        new_repos_df.to_parquet(sample_parquet_files / "repos.parquet", index=False)

        # Full refresh
        load_parquet_to_duckdb(sample_parquet_files, db_path, incremental=False)

        conn = duckdb.connect(str(db_path))
        result = conn.execute("SELECT COUNT(*) FROM raw.repos").fetchone()
        assert result[0] == 1  # Only new data
        conn.close()

    def test_incremental_appends_data(self, sample_parquet_files, tmp_path):
        """Verify incremental mode appends to existing data."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"

        # Initial load
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        # Create new parquet with different data
        new_repos_df = pd.DataFrame([
            {"full_name": "acme/new-repo", "default_branch": "main"}
        ])
        new_repos_df.to_parquet(sample_parquet_files / "repos.parquet", index=False)

        # Incremental load
        load_parquet_to_duckdb(sample_parquet_files, db_path, incremental=True)

        conn = duckdb.connect(str(db_path))
        result = conn.execute("SELECT COUNT(*) FROM raw.repos").fetchone()
        assert result[0] == 3  # Original 2 + new 1
        conn.close()


class TestDuckDBLoaderCLI:
    """Tests for DuckDB loader CLI."""

    def test_cli_with_parquet_dir_and_db_path(self, sample_parquet_files, tmp_path):
        """Test CLI with --parquet-dir and --db-path flags."""
        from duckdb_loader import main
        from unittest.mock import patch

        db_path = tmp_path / "test.duckdb"

        with patch('sys.argv', [
            'duckdb_loader.py',
            '--parquet-dir', str(sample_parquet_files),
            '--db-path', str(db_path)
        ]):
            main()

        assert db_path.exists()

        conn = duckdb.connect(str(db_path))
        result = conn.execute("SELECT COUNT(*) FROM raw.repos").fetchone()
        assert result[0] == 2
        conn.close()

    def test_cli_default_values(self):
        """Test CLI has sensible defaults."""
        from duckdb_loader import create_parser

        parser = create_parser()
        args = parser.parse_args([])

        assert args.parquet_dir == "data/raw"
        assert args.db_path == "data/analytics.duckdb"
        assert args.incremental is False

    def test_cli_incremental_flag(self):
        """Test CLI --incremental flag."""
        from duckdb_loader import create_parser

        parser = create_parser()

        # Without flag
        args = parser.parse_args([])
        assert args.incremental is False

        # With flag
        args = parser.parse_args(['--incremental'])
        assert args.incremental is True


class TestDuckDBLoaderLogging:
    """Tests for DuckDB loader logging."""

    def test_loader_logs_progress(self, sample_parquet_files, tmp_path, capsys):
        """Verify loader outputs progress messages."""
        from duckdb_loader import load_parquet_to_duckdb

        db_path = tmp_path / "test.duckdb"
        load_parquet_to_duckdb(sample_parquet_files, db_path)

        captured = capsys.readouterr()
        # Should log which tables were loaded
        assert "repos" in captured.out.lower() or "Loaded" in captured.out
