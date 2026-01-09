"""
Reviews extractor for /repos/{o}/{r}/pulls/{n}/reviews endpoint.

Extracts PR review data from cursor-sim GitHub-style endpoint.
"""

from pathlib import Path
from typing import List
import pandas as pd
from extractors.base import BaseAPIExtractor


class ReviewsExtractor(BaseAPIExtractor):
    """
    Extract reviews from /repos/{owner}/{repo}/pulls/{number}/reviews endpoint.

    Endpoint: GET /repos/{owner}/{repo}/pulls/{number}/reviews
    Style: GitHub (raw array)
    Pagination: Not paginated (single endpoint per PR)
    """

    def extract(self, output_dir: Path, repo: str, pr_numbers: List[int]) -> None:
        """
        Extract reviews for specified PRs and write to Parquet file.

        Args:
            output_dir: Directory to write reviews.parquet
            repo: Repository name in format "owner/repo"
            pr_numbers: List of PR numbers to fetch reviews for
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        all_reviews = []

        for pr_number in pr_numbers:
            # Fetch reviews for this PR (GitHub-style, not paginated)
            df = self.fetch_github_style(f"/repos/{repo}/pulls/{pr_number}/reviews")

            # Add repo_name and pr_number columns to identify which PR these reviews belong to
            if not df.empty:
                df["repo_name"] = repo
                df["pr_number"] = pr_number
                all_reviews.append(df)

        # Combine all reviews from all PRs
        if all_reviews:
            combined_df = pd.concat(all_reviews, ignore_index=True)
        else:
            # Create empty DataFrame if no reviews found
            combined_df = pd.DataFrame()

        # Write to parquet
        output_file = output_dir / "reviews.parquet"
        self.write_parquet(combined_df, output_file)
