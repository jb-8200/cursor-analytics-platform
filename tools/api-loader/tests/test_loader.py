"""
Tests for the main loader script.

Tests orchestration of extractors and CLI interface.
"""

import pytest
import pandas as pd
from unittest.mock import Mock, patch, MagicMock
from pathlib import Path
import sys


# Mock responses for extractors
@pytest.fixture
def mock_repos_response():
    """Mock /repos endpoint response (GitHub-style raw array)."""
    return [
        {"full_name": "acme/platform", "default_branch": "main"},
        {"full_name": "acme/frontend", "default_branch": "main"}
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
            }
        ],
        "pagination": {"page": 1, "pageSize": 500, "totalPages": 1, "hasNextPage": False}
    }


@pytest.fixture
def mock_prs_response():
    """Mock /repos/{o}/{r}/pulls response (GitHub-style raw array)."""
    return [
        {"number": 1, "state": "merged", "title": "Feature X"},
        {"number": 2, "state": "open", "title": "Fix Y"}
    ]


@pytest.fixture
def mock_reviews_response():
    """Mock /repos/{o}/{r}/pulls/{n}/reviews response."""
    return [
        {"id": 101, "state": "APPROVED"},
        {"id": 102, "state": "CHANGES_REQUESTED"}
    ]


class TestDataLoader:
    """Tests for DataLoader class."""

    def test_loader_writes_all_parquet_files(self, tmp_path, mock_repos_response,
                                              mock_commits_response, mock_prs_response,
                                              mock_reviews_response):
        """Verify loader writes all 4 Parquet files."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            # Set up mock responses in order:
            # 1. repos (non-paginated)
            # 2. commits (paginated cursor-style)
            # 3. PRs for acme/platform (paginated github-style)
            # 4. PRs for acme/platform end
            # 5. PRs for acme/frontend
            # 6. PRs for acme/frontend end
            # 7. Reviews for PR 1
            # 8. Reviews for PR 2
            # 9. Reviews for PR 1 (frontend)
            # 10. Reviews for PR 2 (frontend)
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,       # /repos
                mock_commits_response,     # /analytics/ai-code/commits
                mock_prs_response,         # /repos/acme/platform/pulls page 1
                [],                        # /repos/acme/platform/pulls page 2 (end)
                mock_prs_response,         # /repos/acme/frontend/pulls page 1
                [],                        # /repos/acme/frontend/pulls page 2 (end)
                mock_reviews_response,     # PR 1 reviews (platform)
                mock_reviews_response,     # PR 2 reviews (platform)
                mock_reviews_response,     # PR 1 reviews (frontend)
                mock_reviews_response,     # PR 2 reviews (frontend)
            ]

            loader = DataLoader("http://localhost:8080")
            loader.run(tmp_path)

            # Verify all parquet files exist
            assert (tmp_path / "repos.parquet").exists(), "repos.parquet missing"
            assert (tmp_path / "commits.parquet").exists(), "commits.parquet missing"
            assert (tmp_path / "pull_requests.parquet").exists(), "pull_requests.parquet missing"
            assert (tmp_path / "reviews.parquet").exists(), "reviews.parquet missing"

    def test_loader_repo_discovery_from_endpoint(self, tmp_path, mock_repos_response,
                                                  mock_commits_response):
        """Verify repos come from /repos endpoint, not from commits."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,       # /repos
                mock_commits_response,     # /analytics/ai-code/commits
                [],                        # No PRs
                [],                        # No PRs
                # No reviews since no PRs
            ]

            loader = DataLoader("http://localhost:8080")
            repos = loader._get_repo_list()

            # Should call /repos endpoint and get repos
            assert "acme/platform" in repos
            assert "acme/frontend" in repos

    def test_loader_creates_output_directory(self, tmp_path, mock_repos_response,
                                              mock_commits_response):
        """Verify loader creates output directory if it doesn't exist."""
        from loader import DataLoader

        output_dir = tmp_path / "nested" / "output"

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,
                mock_commits_response,
                [],  # No PRs
                [],  # No PRs
            ]

            loader = DataLoader("http://localhost:8080")
            loader.run(output_dir)

            assert output_dir.exists()

    def test_loader_empty_repos(self, tmp_path, mock_commits_response):
        """Handle empty repos response gracefully."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                [],                        # Empty repos
                mock_commits_response,     # commits
            ]

            loader = DataLoader("http://localhost:8080")
            loader.run(tmp_path)

            # Should still create parquet files (even if empty)
            assert (tmp_path / "repos.parquet").exists()
            assert (tmp_path / "pull_requests.parquet").exists()
            assert (tmp_path / "reviews.parquet").exists()

    def test_loader_progress_logging(self, tmp_path, mock_repos_response,
                                      mock_commits_response, capsys):
        """Verify loader outputs progress messages."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,
                mock_commits_response,
                [],
                [],
            ]

            loader = DataLoader("http://localhost:8080")
            loader.run(tmp_path)

            captured = capsys.readouterr()
            assert "repos" in captured.out.lower() or "Extracting" in captured.out

    def test_loader_extraction_order(self, tmp_path, mock_repos_response,
                                      mock_commits_response, mock_prs_response,
                                      mock_reviews_response):
        """Verify correct extraction order: repos -> commits -> PRs -> reviews."""
        from loader import DataLoader

        call_order = []

        def track_calls(*args, **kwargs):
            url = args[0] if args else kwargs.get('url', '')
            call_order.append(url)
            mock_response = Mock()
            mock_response.status_code = 200

            if '/repos' in url and '/pulls' not in url:
                mock_response.json.return_value = mock_repos_response
            elif '/analytics/ai-code/commits' in url:
                mock_response.json.return_value = mock_commits_response
            elif '/pulls' in url and '/reviews' in url:
                mock_response.json.return_value = mock_reviews_response
            elif '/pulls' in url:
                # Check if this is a pagination request (will have page param)
                if 'page=2' in str(kwargs.get('params', {})) or 'page': 2 in str(kwargs.get('params', {})):
                    mock_response.json.return_value = []
                else:
                    mock_response.json.return_value = mock_prs_response
            else:
                mock_response.json.return_value = []

            return mock_response

        with patch('requests.get', side_effect=track_calls):
            loader = DataLoader("http://localhost:8080")
            loader.run(tmp_path)

        # Verify repos is extracted before PRs
        repos_idx = next((i for i, url in enumerate(call_order) if '/repos' in url and '/pulls' not in url), -1)
        prs_idx = next((i for i, url in enumerate(call_order) if '/pulls' in url and '/reviews' not in url), -1)

        assert repos_idx != -1, "repos endpoint was not called"
        assert repos_idx < prs_idx or prs_idx == -1, "repos should be extracted before PRs"

    def test_loader_continue_on_error_option(self, tmp_path, mock_repos_response,
                                              mock_commits_response):
        """Verify continue-on-error option allows extraction to continue."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200

            def raise_on_prs(*args, **kwargs):
                url = args[0] if args else kwargs.get('url', '')
                mock_response = Mock()
                mock_response.status_code = 200

                if '/repos' in url and '/pulls' not in url:
                    mock_response.json.return_value = mock_repos_response
                elif '/analytics/ai-code/commits' in url:
                    mock_response.json.return_value = mock_commits_response
                elif '/pulls' in url:
                    mock_response.raise_for_status.side_effect = Exception("API Error")
                else:
                    mock_response.json.return_value = []
                return mock_response

            mock_get.side_effect = raise_on_prs

            loader = DataLoader("http://localhost:8080")

            # Should not raise when continue_on_error=True
            loader.run(tmp_path, continue_on_error=True)

            # Repos and commits should still be saved
            assert (tmp_path / "repos.parquet").exists()
            assert (tmp_path / "commits.parquet").exists()


class TestLoaderCLI:
    """Tests for CLI interface."""

    def test_cli_with_url_and_output(self, tmp_path, mock_repos_response,
                                      mock_commits_response):
        """Test CLI with --url and --output flags."""
        from loader import main
        import argparse

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,
                mock_commits_response,
                [],
                [],
            ]

            with patch('sys.argv', ['loader.py', '--url', 'http://localhost:8080',
                                    '--output', str(tmp_path)]):
                main()

            assert (tmp_path / "repos.parquet").exists()

    def test_cli_with_api_key(self, tmp_path, mock_repos_response, mock_commits_response):
        """Test CLI with --api-key flag."""
        from loader import main

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,
                mock_commits_response,
                [],
                [],
            ]

            with patch('sys.argv', ['loader.py', '--url', 'http://localhost:8080',
                                    '--output', str(tmp_path),
                                    '--api-key', 'custom-key']):
                main()

            # Verify the api key was used in requests
            # (Check that auth tuple was passed)
            call_kwargs = mock_get.call_args_list[0][1]
            assert 'auth' in call_kwargs
            assert call_kwargs['auth'][0] == 'custom-key'

    def test_cli_default_values(self):
        """Test CLI has sensible defaults."""
        from loader import create_parser

        parser = create_parser()
        args = parser.parse_args(['--url', 'http://localhost:8080'])

        assert args.url == 'http://localhost:8080'
        assert args.output == 'data/raw'  # default
        assert args.api_key == 'cursor-sim-dev-key'  # default

    def test_cli_continue_on_error_flag(self):
        """Test CLI --continue-on-error flag."""
        from loader import create_parser

        parser = create_parser()

        # Without flag
        args = parser.parse_args(['--url', 'http://localhost:8080'])
        assert args.continue_on_error is False

        # With flag
        args = parser.parse_args(['--url', 'http://localhost:8080', '--continue-on-error'])
        assert args.continue_on_error is True


class TestDataLoaderIntegration:
    """Integration-level tests for data loader."""

    def test_full_extraction_pipeline(self, tmp_path, mock_repos_response,
                                       mock_commits_response, mock_prs_response,
                                       mock_reviews_response):
        """Test full extraction pipeline from repos to reviews."""
        from loader import DataLoader

        with patch('requests.get') as mock_get:
            mock_get.return_value.status_code = 200
            mock_get.return_value.json.side_effect = [
                mock_repos_response,       # /repos
                mock_commits_response,     # /analytics/ai-code/commits
                mock_prs_response,         # /repos/acme/platform/pulls
                [],                        # end pagination
                mock_prs_response,         # /repos/acme/frontend/pulls
                [],                        # end pagination
                mock_reviews_response,     # reviews for PR 1 (platform)
                mock_reviews_response,     # reviews for PR 2 (platform)
                mock_reviews_response,     # reviews for PR 1 (frontend)
                mock_reviews_response,     # reviews for PR 2 (frontend)
            ]

            loader = DataLoader("http://localhost:8080")
            loader.run(tmp_path)

            # Verify repos
            repos_df = pd.read_parquet(tmp_path / "repos.parquet")
            assert len(repos_df) == 2
            assert "full_name" in repos_df.columns

            # Verify commits
            commits_df = pd.read_parquet(tmp_path / "commits.parquet")
            assert len(commits_df) == 1
            assert "commitHash" in commits_df.columns

            # Verify PRs have repo_name column
            prs_df = pd.read_parquet(tmp_path / "pull_requests.parquet")
            assert len(prs_df) == 4  # 2 PRs x 2 repos
            assert "repo_name" in prs_df.columns

            # Verify reviews have repo_name and pr_number columns
            reviews_df = pd.read_parquet(tmp_path / "reviews.parquet")
            assert len(reviews_df) == 8  # 2 reviews x 2 PRs x 2 repos
            assert "repo_name" in reviews_df.columns
            assert "pr_number" in reviews_df.columns
