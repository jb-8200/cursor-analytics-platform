"""
Tests for SQL query modules.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-04 SQL Query Modules

TDD Approach:
1. Write tests first (RED)
2. Implement minimal code (GREEN)
3. Refactor if needed
"""

import pytest
import pandas as pd
from unittest.mock import patch, MagicMock


class TestVelocityQueries:
    """Test velocity query module."""

    @patch("queries.velocity.query")
    def test_get_velocity_data_returns_expected_columns(self, mock_query):
        """Verify velocity query returns required columns."""
        # Mock data with all required columns
        mock_data = pd.DataFrame(
            {
                "week": ["2026-01-06"],
                "repo_name": ["test/repo"],
                "active_developers": [5],
                "total_prs": [10],
                "avg_pr_size": [250],
                "coding_lead_time": [1.5],
                "pickup_time": [0.5],
                "review_lead_time": [2.0],
                "total_cycle_time": [4.0],
                "p50_cycle_time": [3.5],
                "p90_cycle_time": [6.0],
                "avg_ai_ratio": [0.45],
            }
        )
        mock_query.return_value = mock_data

        from queries.velocity import get_velocity_data

        df = get_velocity_data()

        expected_columns = [
            "week",
            "repo_name",
            "active_developers",
            "total_prs",
            "avg_pr_size",
            "coding_lead_time",
            "pickup_time",
            "review_lead_time",
            "total_cycle_time",
            "p50_cycle_time",
            "p90_cycle_time",
            "avg_ai_ratio",
        ]

        for col in expected_columns:
            assert col in df.columns, f"Missing column: {col}"

    @patch("queries.velocity.query")
    def test_get_velocity_data_with_filter(self, mock_query):
        """Verify filter clause is applied to query."""
        mock_query.return_value = pd.DataFrame(
            {"week": ["2026-01-06"], "repo_name": ["acme/platform"], "total_prs": [5]}
        )

        from queries.velocity import get_velocity_data

        df = get_velocity_data("WHERE repo_name = 'acme/platform'")

        # Verify query was called
        mock_query.assert_called_once()

        # Verify WHERE clause is in SQL
        call_args = mock_query.call_args[0][0]
        assert "WHERE repo_name = 'acme/platform'" in call_args

    @patch("queries.velocity.query")
    def test_get_cycle_time_breakdown(self, mock_query):
        """Verify cycle time breakdown query returns component and hours."""
        mock_data = pd.DataFrame(
            {"component": ["Coding", "Pickup", "Review"], "hours": [36.0, 12.0, 48.0]}
        )
        mock_query.return_value = mock_data

        from queries.velocity import get_cycle_time_breakdown

        df = get_cycle_time_breakdown()

        assert "component" in df.columns
        assert "hours" in df.columns
        assert len(df) == 3
        assert set(df["component"].unique()) == {"Coding", "Pickup", "Review"}


class TestAIImpactQueries:
    """Test AI impact query module."""

    @patch("queries.ai_impact.query")
    def test_get_ai_impact_data_returns_expected_columns(self, mock_query):
        """Verify AI impact query returns required columns."""
        mock_data = pd.DataFrame(
            {
                "week": ["2026-01-06"],
                "ai_usage_band": ["medium"],
                "pr_count": [10],
                "avg_ai_ratio": [0.5],
                "avg_coding_lead_time": [1.5],
                "avg_review_cycle_time": [2.0],
                "revert_rate": [0.05],
            }
        )
        mock_query.return_value = mock_data

        from queries.ai_impact import get_ai_impact_data

        df = get_ai_impact_data()

        expected_columns = [
            "week",
            "ai_usage_band",
            "pr_count",
            "avg_ai_ratio",
            "avg_coding_lead_time",
            "avg_review_cycle_time",
            "revert_rate",
        ]

        for col in expected_columns:
            assert col in df.columns, f"Missing column: {col}"

    @patch("queries.ai_impact.query")
    def test_get_ai_impact_data_has_valid_bands(self, mock_query):
        """Verify AI impact query groups by valid bands."""
        mock_data = pd.DataFrame(
            {
                "week": ["2026-01-06", "2026-01-06", "2026-01-06"],
                "ai_usage_band": ["low", "medium", "high"],
                "pr_count": [5, 10, 3],
                "avg_ai_ratio": [0.1, 0.5, 0.9],
                "avg_coding_lead_time": [2.0, 1.5, 1.0],
                "avg_review_cycle_time": [3.0, 2.5, 2.0],
                "revert_rate": [0.08, 0.05, 0.03],
            }
        )
        mock_query.return_value = mock_data

        from queries.ai_impact import get_ai_impact_data

        df = get_ai_impact_data()

        valid_bands = {"low", "medium", "high"}
        actual_bands = set(df["ai_usage_band"].unique())
        assert actual_bands.issubset(valid_bands)


class TestQualityQueries:
    """Test quality query module."""

    @patch("queries.quality.query")
    def test_get_quality_data_returns_expected_columns(self, mock_query):
        """Verify quality query returns required columns."""
        mock_data = pd.DataFrame(
            {
                "week": ["2026-01-06"],
                "repo_name": ["test/repo"],
                "total_prs": [10],
                "reverted_prs": [1],
                "revert_rate": [0.1],
                "avg_time_to_revert": [2.5],
                "bug_fix_prs": [2],
                "bug_fix_rate": [0.2],
            }
        )
        mock_query.return_value = mock_data

        from queries.quality import get_quality_data

        df = get_quality_data()

        expected_columns = [
            "week",
            "repo_name",
            "total_prs",
            "reverted_prs",
            "revert_rate",
            "avg_time_to_revert",
            "bug_fix_prs",
            "bug_fix_rate",
        ]

        for col in expected_columns:
            assert col in df.columns, f"Missing column: {col}"

    @patch("queries.quality.query")
    def test_get_quality_data_with_filter(self, mock_query):
        """Verify filter clause is applied."""
        mock_query.return_value = pd.DataFrame(
            {"week": ["2026-01-06"], "repo_name": ["test/repo"], "revert_rate": [0.05]}
        )

        from queries.quality import get_quality_data

        df = get_quality_data("WHERE repo_name = 'test/repo'")

        # Verify query was called
        mock_query.assert_called_once()

        # Verify WHERE clause is in SQL
        call_args = mock_query.call_args[0][0]
        assert "WHERE repo_name = 'test/repo'" in call_args


class TestReviewCostsQueries:
    """Test review costs query module."""

    @patch("queries.review_costs.query")
    def test_get_review_costs_data_returns_expected_columns(self, mock_query):
        """Verify review costs query returns required columns."""
        mock_data = pd.DataFrame(
            {
                "week": ["2026-01-06"],
                "repo_name": ["test/repo"],
                "total_prs": [10],
                "avg_review_iterations": [2.5],
                "avg_reviewers_per_pr": [2.0],
                "avg_review_comments": [5.0],
                "avg_review_time": [3.0],
                "total_review_hours": [30.0],
            }
        )
        mock_query.return_value = mock_data

        from queries.review_costs import get_review_costs_data

        df = get_review_costs_data()

        expected_columns = [
            "week",
            "repo_name",
            "total_prs",
            "avg_review_iterations",
            "avg_reviewers_per_pr",
            "avg_review_comments",
            "avg_review_time",
            "total_review_hours",
        ]

        for col in expected_columns:
            assert col in df.columns, f"Missing column: {col}"

    @patch("queries.review_costs.query")
    def test_get_review_costs_data_with_filter(self, mock_query):
        """Verify filter clause is applied."""
        mock_query.return_value = pd.DataFrame(
            {
                "week": ["2026-01-06"],
                "repo_name": ["test/repo"],
                "avg_review_iterations": [2.5],
            }
        )

        from queries.review_costs import get_review_costs_data

        df = get_review_costs_data("WHERE repo_name = 'test/repo'")

        # Verify query was called
        mock_query.assert_called_once()

        # Verify WHERE clause is in SQL
        call_args = mock_query.call_args[0][0]
        assert "WHERE repo_name = 'test/repo'" in call_args


class TestQueryIntegration:
    """Integration tests for query modules (when P8 marts exist)."""

    @pytest.mark.skip(reason="Requires P8 mart tables to exist")
    def test_velocity_query_executes_successfully(self):
        """Verify velocity query executes against real database."""
        from queries.velocity import get_velocity_data

        df = get_velocity_data()
        assert isinstance(df, pd.DataFrame)

    @pytest.mark.skip(reason="Requires P8 mart tables to exist")
    def test_ai_impact_query_executes_successfully(self):
        """Verify AI impact query executes against real database."""
        from queries.ai_impact import get_ai_impact_data

        df = get_ai_impact_data()
        assert isinstance(df, pd.DataFrame)

    @pytest.mark.skip(reason="Requires P8 mart tables to exist")
    def test_quality_query_executes_successfully(self):
        """Verify quality query executes against real database."""
        from queries.quality import get_quality_data

        df = get_quality_data()
        assert isinstance(df, pd.DataFrame)

    @pytest.mark.skip(reason="Requires P8 mart tables to exist")
    def test_review_costs_query_executes_successfully(self):
        """Verify review costs query executes against real database."""
        from queries.review_costs import get_review_costs_data

        df = get_review_costs_data()
        assert isinstance(df, pd.DataFrame)
