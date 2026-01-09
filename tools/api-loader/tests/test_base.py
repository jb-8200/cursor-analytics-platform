"""
Tests for base API extractor.
"""

import pytest
import pandas as pd
from unittest.mock import Mock, patch
import requests


@pytest.fixture
def mock_github_response():
    """Mock response for GitHub-style endpoints (raw arrays)."""
    return [
        {"full_name": "acme/platform", "default_branch": "main", "open_pull_requests_count": 5},
        {"full_name": "acme/frontend", "default_branch": "main", "open_pull_requests_count": 2}
    ]


@pytest.fixture
def mock_cursor_response():
    """Mock response for Cursor Analytics endpoints (wrapped object)."""
    return {
        "data": [
            {
                "commitHash": "abc123",
                "userId": "user_001",
                "userEmail": "dev@example.com",
                "repoName": "acme/platform",
                "tabLinesAdded": 50,
                "composerLinesAdded": 30,
                "nonAiLinesAdded": 20,
                "commitTs": "2026-01-01T10:00:00Z"
            }
        ],
        "pagination": {
            "page": 1,
            "pageSize": 100,
            "totalPages": 1,
            "hasNextPage": False,
            "hasPreviousPage": False
        },
        "params": {
            "from": "2025-10-01",
            "to": "2026-01-01",
            "page": 1,
            "pageSize": 100
        }
    }


class TestBaseExtractor:
    """Tests for BaseAPIExtractor."""

    def test_extract_github_repos_raw_array(self, mock_github_response):
        """Verify loader handles raw array responses (GitHub-style)."""
        from extractors.base import BaseAPIExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_github_response

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_github_style("/repos")

            assert len(df) == 2
            assert "full_name" in df.columns
            assert df["full_name"].tolist() == ["acme/platform", "acme/frontend"]

    def test_extract_cursor_commits_wrapped_object(self, mock_cursor_response):
        """Verify loader handles wrapped object responses (Cursor Analytics-style)."""
        from extractors.base import BaseAPIExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_cursor_response

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_cursor_style("/analytics/ai-code/commits")

            assert len(df) == 1
            assert "commitHash" in df.columns
            assert df["commitHash"].iloc[0] == "abc123"

    def test_pagination_cursor_style(self):
        """Verify cursor-style pagination with hasNextPage."""
        from extractors.base import BaseAPIExtractor

        # Page 1: hasNextPage = True
        page1 = {
            "data": [{"id": 1}, {"id": 2}],
            "pagination": {"page": 1, "pageSize": 2, "hasNextPage": True, "hasPreviousPage": False}
        }
        # Page 2: hasNextPage = False
        page2 = {
            "data": [{"id": 3}],
            "pagination": {"page": 2, "pageSize": 2, "hasNextPage": False, "hasPreviousPage": True}
        }

        with patch('requests.get') as mock_get:
            # First call returns page 1, second call returns page 2
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [page1, page2]

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_cursor_style_paginated("/analytics/ai-code/commits")

            # Should have combined both pages
            assert len(df) == 3
            assert df["id"].tolist() == [1, 2, 3]
            # Should have made 2 requests
            assert mock_get.call_count == 2

    def test_pagination_github_style(self):
        """Verify github-style pagination with empty array termination."""
        from extractors.base import BaseAPIExtractor

        page1 = [{"number": 1}, {"number": 2}]
        page2 = [{"number": 3}]
        page3 = []  # Empty array signals end

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [page1, page2, page3]

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_github_style_paginated("/repos/acme/platform/pulls", params={"state": "all"})

            # Should have combined first two pages, stopped at empty
            assert len(df) == 3
            assert df["number"].tolist() == [1, 2, 3]
            # Should have made 3 requests
            assert mock_get.call_count == 3

    def test_empty_response_cursor_style(self):
        """Handle empty Cursor Analytics response gracefully."""
        from extractors.base import BaseAPIExtractor

        empty_response = {
            "data": [],
            "pagination": {"page": 1, "pageSize": 100, "hasNextPage": False, "hasPreviousPage": False}
        }

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = empty_response

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_cursor_style("/analytics/ai-code/commits")

            assert len(df) == 0
            assert isinstance(df, pd.DataFrame)

    def test_empty_response_github_style(self):
        """Handle empty GitHub response gracefully."""
        from extractors.base import BaseAPIExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = []

            extractor = BaseAPIExtractor("http://localhost:8080")
            df = extractor.fetch_github_style("/repos")

            assert len(df) == 0
            assert isinstance(df, pd.DataFrame)

    def test_authentication_included(self):
        """Verify API key is included in requests."""
        from extractors.base import BaseAPIExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = []

            extractor = BaseAPIExtractor("http://localhost:8080", api_key="test-key")
            extractor.fetch_github_style("/repos")

            # Verify auth was included
            call_kwargs = mock_get.call_args[1]
            assert "auth" in call_kwargs
            assert call_kwargs["auth"] == ("test-key", "")

    def test_http_error_handling(self):
        """Verify HTTP errors raise appropriate exceptions."""
        from extractors.base import BaseAPIExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 404
            mock_get.return_value.raise_for_status.side_effect = requests.HTTPError("Not Found")

            extractor = BaseAPIExtractor("http://localhost:8080")

            with pytest.raises(requests.HTTPError):
                extractor.fetch_github_style("/invalid-endpoint")

    def test_write_to_parquet(self, tmp_path):
        """Verify data can be written to Parquet files."""
        from extractors.base import BaseAPIExtractor
        import pyarrow.parquet as pq

        df = pd.DataFrame({
            "id": [1, 2, 3],
            "name": ["Alice", "Bob", "Charlie"]
        })

        extractor = BaseAPIExtractor("http://localhost:8080")
        output_file = tmp_path / "test.parquet"
        extractor.write_parquet(df, output_file)

        # Verify file exists and is readable
        assert output_file.exists()

        # Read back and verify contents
        result_df = pd.read_parquet(output_file)
        assert len(result_df) == 3
        assert result_df["name"].tolist() == ["Alice", "Bob", "Charlie"]
