"""
Pull Requests extractor for /repos/{o}/{r}/pulls endpoint.

Extracts PR data from cursor-sim GitHub-style endpoint.
"""

from pathlib import Path
from typing import List
import pandas as pd
from extractors.base import BaseAPIExtractor


class PRsExtractor(BaseAPIExtractor):
    """
    Extract pull requests from /repos/{owner}/{repo}/pulls endpoint.

    Endpoint: GET /repos/{owner}/{repo}/pulls
    Style: GitHub (raw array)
    Pagination: page/per_page with empty array termination
    """

    def extract(self, output_dir: Path, repos: List[str]) -> None:
        """
        Extract pull requests for specified repos and write to Parquet file.

        Args:
            output_dir: Directory to write pull_requests.parquet
            repos: List of repo names in format "owner/repo"
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        all_prs = []

        for repo in repos:
            # Fetch PRs with pagination (GitHub-style)
            params = {"state": "all"}
            df = self.fetch_github_style_paginated(f"/repos/{repo}/pulls", params=params)

            # Add repo_name column to identify which repo these PRs belong to
            if not df.empty:
                df["repo_name"] = repo
                all_prs.append(df)

        # Combine all PRs from all repos
        if all_prs:
            combined_df = pd.concat(all_prs, ignore_index=True)
        else:
            # Create empty DataFrame if no PRs found
            combined_df = pd.DataFrame()

        # Write to parquet
        output_file = output_dir / "pull_requests.parquet"
        self.write_parquet(combined_df, output_file)
