"""
Tests for refresh_data pipeline function.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-10 Refresh Pipeline

TDD Approach:
1. Write failing tests first (RED)
2. Implement minimal functionality (GREEN)
3. Refactor while keeping tests passing
"""

import os
import subprocess
from unittest.mock import Mock, patch, MagicMock
import pytest


def test_refresh_data_returns_false_in_snowflake_mode(monkeypatch):
    """Verify refresh_data returns False in Snowflake/production mode."""
    monkeypatch.setenv("DB_MODE", "snowflake")

    # Import after setting env vars
    from db.connector import refresh_data

    result = refresh_data()
    assert result is False


def test_refresh_data_returns_true_on_success(monkeypatch, tmp_path):
    """Verify refresh_data returns True when all steps succeed."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))
    monkeypatch.setenv("CURSOR_SIM_URL", "http://localhost:8080")

    from db.connector import refresh_data

    # Mock subprocess.run to simulate success
    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        # Mock the duckdb_loader import
        with patch("db.connector.load_parquet_to_duckdb") as mock_loader:
            result = refresh_data()

            assert result is True
            # Verify subprocess was called for loader and dbt
            assert mock_run.call_count == 2
            # Verify duckdb loader was called
            mock_loader.assert_called_once()


def test_refresh_data_calls_loader_with_correct_args(monkeypatch, tmp_path):
    """Verify refresh_data calls loader.py with correct arguments."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("CURSOR_SIM_URL", "http://test:8080")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb"):
            refresh_data()

            # Check first call (loader)
            first_call = mock_run.call_args_list[0]
            args = first_call[0][0]

            assert "python" in args
            assert "tools/api-loader/loader.py" in args
            assert "--url" in args
            assert "http://test:8080" in args
            assert "--output" in args
            assert "/data/raw" in args


def test_refresh_data_calls_duckdb_loader(monkeypatch, tmp_path):
    """Verify refresh_data calls duckdb_loader with correct paths."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    db_path = str(tmp_path / "test.duckdb")
    monkeypatch.setenv("DUCKDB_PATH", db_path)

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb") as mock_loader:
            refresh_data()

            # Verify loader was called with correct paths
            mock_loader.assert_called_once_with("/data/raw", db_path)


def test_refresh_data_calls_dbt_build(monkeypatch, tmp_path):
    """Verify refresh_data runs dbt build with correct config."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb"):
            refresh_data()

            # Check second call (dbt)
            second_call = mock_run.call_args_list[1]
            args = second_call[0][0]
            kwargs = second_call[1]

            assert "dbt" in args
            assert "build" in args
            assert "--target" in args
            assert "dev" in args
            assert kwargs.get("cwd") == "/app/dbt"


def test_refresh_data_handles_loader_failure(monkeypatch, tmp_path):
    """Verify refresh_data returns False when loader fails."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        # Simulate loader failure
        mock_run.side_effect = subprocess.CalledProcessError(
            1, "loader.py", stderr="Connection failed"
        )

        result = refresh_data()

        assert result is False


def test_refresh_data_handles_duckdb_loader_failure(monkeypatch, tmp_path):
    """Verify refresh_data returns False when duckdb_loader fails."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb") as mock_loader:
            # Simulate loader failure
            mock_loader.side_effect = Exception("DuckDB connection failed")

            result = refresh_data()

            assert result is False


def test_refresh_data_handles_dbt_failure(monkeypatch, tmp_path):
    """Verify refresh_data returns False when dbt fails."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        # First call (loader) succeeds, second call (dbt) fails
        mock_run.side_effect = [
            Mock(returncode=0, stderr="", stdout=""),
            subprocess.CalledProcessError(1, "dbt", stderr="Model failed")
        ]

        with patch("db.connector.load_parquet_to_duckdb"):
            result = refresh_data()

            assert result is False


def test_refresh_data_uses_default_cursor_sim_url(monkeypatch, tmp_path):
    """Verify refresh_data uses default CURSOR_SIM_URL when not set."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))
    # Don't set CURSOR_SIM_URL - should use default
    monkeypatch.delenv("CURSOR_SIM_URL", raising=False)

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb"):
            refresh_data()

            # Check that default URL is used
            first_call = mock_run.call_args_list[0]
            args = first_call[0][0]

            assert "http://localhost:8080" in args


def test_refresh_data_uses_custom_cursor_sim_url(monkeypatch, tmp_path):
    """Verify refresh_data uses custom CURSOR_SIM_URL when set."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))
    monkeypatch.setenv("CURSOR_SIM_URL", "http://custom:9999")

    from db.connector import refresh_data

    with patch("subprocess.run") as mock_run:
        mock_run.return_value = Mock(returncode=0, stderr="", stdout="")

        with patch("db.connector.load_parquet_to_duckdb"):
            refresh_data()

            # Check that custom URL is used
            first_call = mock_run.call_args_list[0]
            args = first_call[0][0]

            assert "http://custom:9999" in args


def test_refresh_data_pipeline_order(monkeypatch, tmp_path):
    """Verify refresh_data executes steps in correct order."""
    monkeypatch.setenv("DB_MODE", "duckdb")
    monkeypatch.setenv("DUCKDB_PATH", str(tmp_path / "test.duckdb"))

    from db.connector import refresh_data

    call_order = []

    def mock_subprocess_run(*args, **kwargs):
        if "loader.py" in args[0]:
            call_order.append("loader")
        elif "dbt" in args[0]:
            call_order.append("dbt")
        return Mock(returncode=0, stderr="", stdout="")

    def mock_duckdb_loader(*args, **kwargs):
        call_order.append("duckdb_loader")

    with patch("subprocess.run", side_effect=mock_subprocess_run):
        with patch("db.connector.load_parquet_to_duckdb", side_effect=mock_duckdb_loader):
            refresh_data()

            # Verify correct order: loader -> duckdb_loader -> dbt
            assert call_order == ["loader", "duckdb_loader", "dbt"]
