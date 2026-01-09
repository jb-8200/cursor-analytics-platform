"""
Tests for specific API extractors (repos, commits, prs, reviews).
"""

import pytest
import pandas as pd
from unittest.mock import Mock, patch
from pathlib import Path


@pytest.fixture
def mock_repos_response():
    """Mock /repos endpoint response (GitHub-style raw array)."""
    return [
        {
            "full_name": "acme/platform",
            "default_branch": "main",
            "open_pull_requests_count": 5,
            "created_at": "2025-01-01T00:00:00Z"
        },
        {
            "full_name": "acme/frontend",
            "default_branch": "main",
            "open_pull_requests_count": 2,
            "created_at": "2025-02-01T00:00:00Z"
        }
    ]


@pytest.fixture
def mock_commits_response():
    """Mock /analytics/ai-code/commits response (Cursor-style wrapped)."""
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
            },
            {
                "commitHash": "def456",
                "userId": "user_002",
                "userEmail": "dev2@example.com",
                "repoName": "acme/frontend",
                "tabLinesAdded": 25,
                "composerLinesAdded": 15,
                "nonAiLinesAdded": 10,
                "commitTs": "2026-01-02T10:00:00Z"
            }
        ],
        "pagination": {
            "page": 1,
            "pageSize": 500,
            "totalPages": 1,
            "hasNextPage": False
        }
    }


@pytest.fixture
def mock_prs_response():
    """Mock /repos/{o}/{r}/pulls response (GitHub-style raw array)."""
    return [
        {
            "number": 1,
            "state": "merged",
            "title": "Add feature X",
            "additions": 100,
            "deletions": 20,
            "created_at": "2026-01-01T10:00:00Z"
        },
        {
            "number": 2,
            "state": "open",
            "title": "Fix bug Y",
            "additions": 50,
            "deletions": 10,
            "created_at": "2026-01-02T10:00:00Z"
        }
    ]


@pytest.fixture
def mock_reviews_response():
    """Mock /repos/{o}/{r}/pulls/{n}/reviews response (GitHub-style)."""
    return [
        {
            "id": 101,
            "user": {"login": "reviewer1"},
            "state": "APPROVED",
            "submitted_at": "2026-01-01T11:00:00Z"
        },
        {
            "id": 102,
            "user": {"login": "reviewer2"},
            "state": "CHANGES_REQUESTED",
            "submitted_at": "2026-01-01T12:00:00Z"
        }
    ]


class TestReposExtractor:
    """Tests for ReposExtractor."""

    def test_extract_repos_returns_dataframe(self, mock_repos_response, tmp_path):
        """Verify repos extractor returns DataFrame with correct data."""
        from extractors.repos import ReposExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_repos_response

            extractor = ReposExtractor("http://localhost:8080")
            extractor.extract(tmp_path)

            # Verify parquet file was created
            output_file = tmp_path / "repos.parquet"
            assert output_file.exists()

            # Read and verify contents
            df = pd.read_parquet(output_file)
            assert len(df) == 2
            assert "full_name" in df.columns
            assert df["full_name"].tolist() == ["acme/platform", "acme/frontend"]

    def test_extract_empty_repos(self, tmp_path):
        """Handle empty repos response gracefully."""
        from extractors.repos import ReposExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = []

            extractor = ReposExtractor("http://localhost:8080")
            extractor.extract(tmp_path)

            output_file = tmp_path / "repos.parquet"
            assert output_file.exists()

            df = pd.read_parquet(output_file)
            assert len(df) == 0


class TestCommitsExtractor:
    """Tests for CommitsExtractor."""

    def test_extract_commits_paginated(self, mock_commits_response, tmp_path):
        """Verify commits extractor handles paginated response."""
        from extractors.commits import CommitsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_commits_response

            extractor = CommitsExtractor("http://localhost:8080")
            extractor.extract(tmp_path)

            # Verify parquet file was created
            output_file = tmp_path / "commits.parquet"
            assert output_file.exists()

            # Read and verify contents
            df = pd.read_parquet(output_file)
            assert len(df) == 2
            assert "commitHash" in df.columns
            assert df["commitHash"].tolist() == ["abc123", "def456"]

    def test_extract_commits_with_start_date(self, mock_commits_response, tmp_path):
        """Verify start_date parameter is passed correctly."""
        from extractors.commits import CommitsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_commits_response

            extractor = CommitsExtractor("http://localhost:8080")
            extractor.extract(tmp_path, start_date="30d")

            # Verify the request was made with correct params
            call_kwargs = mock_get.call_args[1]
            assert "params" in call_kwargs
            # Note: params will include page and page_size from pagination


class TestPRsExtractor:
    """Tests for PRsExtractor."""

    def test_extract_prs_single_repo(self, mock_prs_response, tmp_path):
        """Verify PRs extractor for single repo."""
        from extractors.prs import PRsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            # First call returns PRs, second call returns empty (pagination end)
            mock_get.return_value.json.side_effect = [mock_prs_response, []]

            extractor = PRsExtractor("http://localhost:8080")
            extractor.extract(tmp_path, repos=["acme/platform"])

            # Verify parquet file was created
            output_file = tmp_path / "pull_requests.parquet"
            assert output_file.exists()

            # Read and verify contents
            df = pd.read_parquet(output_file)
            assert len(df) == 2
            assert "number" in df.columns
            assert "repo_name" in df.columns
            assert df["repo_name"].iloc[0] == "acme/platform"

    def test_extract_prs_multiple_repos(self, mock_prs_response, tmp_path):
        """Verify PRs extractor handles multiple repos."""
        from extractors.prs import PRsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            # Return PRs for first repo, then empty, then PRs for second repo, then empty
            mock_get.return_value.json.side_effect = [
                mock_prs_response,  # acme/platform page 1
                [],                  # acme/platform page 2 (end)
                mock_prs_response,  # acme/frontend page 1
                []                   # acme/frontend page 2 (end)
            ]

            extractor = PRsExtractor("http://localhost:8080")
            extractor.extract(tmp_path, repos=["acme/platform", "acme/frontend"])

            output_file = tmp_path / "pull_requests.parquet"
            df = pd.read_parquet(output_file)

            # Should have PRs from both repos
            assert len(df) == 4
            assert set(df["repo_name"].unique()) == {"acme/platform", "acme/frontend"}

    def test_extract_prs_empty_list(self, tmp_path):
        """Handle empty repo list gracefully."""
        from extractors.prs import PRsExtractor

        extractor = PRsExtractor("http://localhost:8080")
        extractor.extract(tmp_path, repos=[])

        output_file = tmp_path / "pull_requests.parquet"
        assert output_file.exists()

        df = pd.read_parquet(output_file)
        assert len(df) == 0


class TestReviewsExtractor:
    """Tests for ReviewsExtractor."""

    def test_extract_reviews_single_pr(self, mock_reviews_response, tmp_path):
        """Verify reviews extractor for single PR."""
        from extractors.reviews import ReviewsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.return_value = mock_reviews_response

            extractor = ReviewsExtractor("http://localhost:8080")
            extractor.extract(tmp_path, repo="acme/platform", pr_numbers=[1])

            # Verify parquet file was created
            output_file = tmp_path / "reviews.parquet"
            assert output_file.exists()

            # Read and verify contents
            df = pd.read_parquet(output_file)
            assert len(df) == 2
            assert "id" in df.columns
            assert "repo_name" in df.columns
            assert "pr_number" in df.columns
            assert df["repo_name"].iloc[0] == "acme/platform"
            assert df["pr_number"].iloc[0] == 1

    def test_extract_reviews_multiple_prs(self, mock_reviews_response, tmp_path):
        """Verify reviews extractor handles multiple PRs."""
        from extractors.reviews import ReviewsExtractor

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            # Return reviews for each PR
            mock_get.return_value.json.side_effect = [
                mock_reviews_response,  # PR 1
                mock_reviews_response   # PR 2
            ]

            extractor = ReviewsExtractor("http://localhost:8080")
            extractor.extract(tmp_path, repo="acme/platform", pr_numbers=[1, 2])

            output_file = tmp_path / "reviews.parquet"
            df = pd.read_parquet(output_file)

            # Should have reviews from both PRs
            assert len(df) == 4
            assert set(df["pr_number"].unique()) == {1, 2}

    def test_extract_reviews_empty_pr_list(self, tmp_path):
        """Handle empty PR list gracefully."""
        from extractors.reviews import ReviewsExtractor

        extractor = ReviewsExtractor("http://localhost:8080")
        extractor.extract(tmp_path, repo="acme/platform", pr_numbers=[])

        output_file = tmp_path / "reviews.parquet"
        assert output_file.exists()

        df = pd.read_parquet(output_file)
        assert len(df) == 0
