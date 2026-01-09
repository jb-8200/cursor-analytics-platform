"""
Tests for schema validation logic.
"""

import pytest
import pandas as pd
import json
from pathlib import Path
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from validators import validate_dataframe, SchemaValidationError


@pytest.fixture
def schema_dir(tmp_path):
    """Create temporary schema directory with test schemas."""
    schema_dir = tmp_path / "schemas"
    schema_dir.mkdir()

    # Create commits schema
    commits_schema = {
        "required_columns": [
            "commitHash",
            "userEmail",
            "repoName",
            "tabLinesAdded",
            "composerLinesAdded",
            "nonAiLinesAdded",
            "commitTs"
        ]
    }
    (schema_dir / "commits.json").write_text(json.dumps(commits_schema))

    # Create pull_requests schema
    prs_schema = {
        "required_columns": [
            "number",
            "repo_name",
            "author_email",
            "state",
            "additions",
            "deletions",
            "ai_ratio",
            "was_reverted",
            "created_at"
        ]
    }
    (schema_dir / "pull_requests.json").write_text(json.dumps(prs_schema))

    # Create reviews schema
    reviews_schema = {
        "required_columns": [
            "id",
            "repo_name",
            "pr_number",
            "state",
            "submitted_at"
        ]
    }
    (schema_dir / "reviews.json").write_text(json.dumps(reviews_schema))

    # Create repos schema
    repos_schema = {
        "required_columns": [
            "full_name",
            "default_branch",
            "created_at"
        ]
    }
    (schema_dir / "repos.json").write_text(json.dumps(repos_schema))

    return schema_dir


class TestSchemaValidation:
    """Test schema validation for extracted data."""

    def test_validate_commits_valid_schema(self, schema_dir):
        """Valid commits DataFrame passes validation."""
        df = pd.DataFrame([
            {
                "commitHash": "abc123",
                "userEmail": "dev@example.com",
                "repoName": "acme/platform",
                "tabLinesAdded": 50,
                "composerLinesAdded": 30,
                "nonAiLinesAdded": 20,
                "commitTs": "2026-01-01T10:00:00Z"
            }
        ])

        schema_path = schema_dir / "commits.json"

        # Should not raise any exception
        validate_dataframe(df, schema_path)

    def test_validate_commits_missing_columns(self, schema_dir):
        """Missing columns raises SchemaValidationError with clear message."""
        df = pd.DataFrame([
            {
                "commitHash": "abc123",
                "userEmail": "dev@example.com",
                "repoName": "acme/platform"
                # Missing: tabLinesAdded, composerLinesAdded, nonAiLinesAdded, commitTs
            }
        ])

        schema_path = schema_dir / "commits.json"

        with pytest.raises(SchemaValidationError) as exc_info:
            validate_dataframe(df, schema_path)

        error_msg = str(exc_info.value)
        assert "Missing required columns" in error_msg
        assert "tabLinesAdded" in error_msg
        assert "composerLinesAdded" in error_msg
        assert "nonAiLinesAdded" in error_msg
        assert "commitTs" in error_msg

    def test_validate_pull_requests_valid_schema(self, schema_dir):
        """Valid pull requests DataFrame passes validation."""
        df = pd.DataFrame([
            {
                "number": 1,
                "repo_name": "acme/platform",
                "author_email": "dev@example.com",
                "state": "merged",
                "additions": 100,
                "deletions": 20,
                "ai_ratio": 0.65,
                "was_reverted": False,
                "created_at": "2026-01-01T10:00:00Z"
            }
        ])

        schema_path = schema_dir / "pull_requests.json"

        # Should not raise any exception
        validate_dataframe(df, schema_path)

    def test_validate_pull_requests_missing_columns(self, schema_dir):
        """Missing PR columns raises SchemaValidationError."""
        df = pd.DataFrame([
            {
                "number": 1,
                "state": "merged"
                # Missing many required columns
            }
        ])

        schema_path = schema_dir / "pull_requests.json"

        with pytest.raises(SchemaValidationError) as exc_info:
            validate_dataframe(df, schema_path)

        error_msg = str(exc_info.value)
        assert "Missing required columns" in error_msg
        assert "repo_name" in error_msg
        assert "author_email" in error_msg

    def test_validate_reviews_valid_schema(self, schema_dir):
        """Valid reviews DataFrame passes validation."""
        df = pd.DataFrame([
            {
                "id": 101,
                "repo_name": "acme/platform",
                "pr_number": 1,
                "state": "APPROVED",
                "submitted_at": "2026-01-01T11:00:00Z"
            }
        ])

        schema_path = schema_dir / "reviews.json"

        # Should not raise any exception
        validate_dataframe(df, schema_path)

    def test_validate_repos_valid_schema(self, schema_dir):
        """Valid repos DataFrame passes validation."""
        df = pd.DataFrame([
            {
                "full_name": "acme/platform",
                "default_branch": "main",
                "created_at": "2025-01-01T00:00:00Z"
            }
        ])

        schema_path = schema_dir / "repos.json"

        # Should not raise any exception
        validate_dataframe(df, schema_path)

    def test_validate_empty_dataframe(self, schema_dir):
        """Empty DataFrame passes validation (no rows to validate)."""
        df = pd.DataFrame()

        schema_path = schema_dir / "commits.json"

        # Empty DataFrame should pass (no data to validate)
        validate_dataframe(df, schema_path)

    def test_validate_extra_columns_ok(self, schema_dir):
        """Extra columns beyond required are allowed."""
        df = pd.DataFrame([
            {
                "commitHash": "abc123",
                "userEmail": "dev@example.com",
                "repoName": "acme/platform",
                "tabLinesAdded": 50,
                "composerLinesAdded": 30,
                "nonAiLinesAdded": 20,
                "commitTs": "2026-01-01T10:00:00Z",
                "extraColumn1": "extra",
                "extraColumn2": 123
            }
        ])

        schema_path = schema_dir / "commits.json"

        # Extra columns should be allowed
        validate_dataframe(df, schema_path)

    def test_validate_missing_schema_file(self):
        """Missing schema file raises FileNotFoundError."""
        df = pd.DataFrame([{"col1": "value"}])

        with pytest.raises(FileNotFoundError):
            validate_dataframe(df, Path("nonexistent_schema.json"))

    def test_validate_invalid_json_schema(self, tmp_path):
        """Invalid JSON in schema file raises ValueError."""
        invalid_schema = tmp_path / "invalid.json"
        invalid_schema.write_text("not valid json {")

        df = pd.DataFrame([{"col1": "value"}])

        with pytest.raises(ValueError) as exc_info:
            validate_dataframe(df, invalid_schema)

        assert "Failed to parse schema file" in str(exc_info.value)

    def test_validate_schema_missing_required_columns_key(self, tmp_path):
        """Schema without 'required_columns' key raises ValueError."""
        invalid_schema = tmp_path / "invalid.json"
        invalid_schema.write_text(json.dumps({"some_other_key": []}))

        df = pd.DataFrame([{"col1": "value"}])

        with pytest.raises(ValueError) as exc_info:
            validate_dataframe(df, invalid_schema)

        assert "Schema must contain 'required_columns' key" in str(exc_info.value)

    def test_error_message_lists_all_missing_columns(self, schema_dir):
        """Error message clearly lists all missing columns."""
        df = pd.DataFrame([
            {
                "commitHash": "abc123"
                # Missing all other required columns
            }
        ])

        schema_path = schema_dir / "commits.json"

        with pytest.raises(SchemaValidationError) as exc_info:
            validate_dataframe(df, schema_path)

        error_msg = str(exc_info.value)
        # All missing columns should be listed
        missing_cols = [
            "userEmail", "repoName", "tabLinesAdded",
            "composerLinesAdded", "nonAiLinesAdded", "commitTs"
        ]
        for col in missing_cols:
            assert col in error_msg
