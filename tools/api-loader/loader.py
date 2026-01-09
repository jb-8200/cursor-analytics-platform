"""
Main loader script for cursor-sim data extraction.

Orchestrates extraction of repos, commits, PRs, and reviews from cursor-sim API
and writes them to Parquet files.
"""

import argparse
import logging
import sys
from pathlib import Path
from typing import List, Optional

import pandas as pd

from extractors import ReposExtractor, CommitsExtractor, PRsExtractor, ReviewsExtractor

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class DataLoader:
    """
    Orchestrates data extraction from cursor-sim API.

    Extracts data in the correct order:
    1. Repos (discover available repositories)
    2. Commits (from Cursor Analytics endpoint)
    3. PRs (for each repo)
    4. Reviews (for each PR in each repo)
    """

    def __init__(self, base_url: str, api_key: str = "cursor-sim-dev-key"):
        """
        Initialize the data loader.

        Args:
            base_url: Base URL for cursor-sim API (e.g., http://localhost:8080)
            api_key: API key for authentication
        """
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key

        # Initialize extractors
        self.repos_extractor = ReposExtractor(base_url, api_key)
        self.commits_extractor = CommitsExtractor(base_url, api_key)
        self.prs_extractor = PRsExtractor(base_url, api_key)
        self.reviews_extractor = ReviewsExtractor(base_url, api_key)

    def _get_repo_list(self) -> List[str]:
        """
        Get list of repositories from /repos endpoint.

        Returns:
            List of repository names in format "owner/repo"
        """
        df = self.repos_extractor.fetch_github_style("/repos")
        if df.empty:
            return []
        return df["full_name"].tolist()

    def run(
        self,
        output_dir: Path,
        start_date: str = "90d",
        continue_on_error: bool = False
    ) -> None:
        """
        Run the full extraction pipeline.

        Extraction order:
        1. Repos from /repos endpoint
        2. Commits from /analytics/ai-code/commits
        3. PRs for each repo from /repos/{o}/{r}/pulls
        4. Reviews for each PR from /repos/{o}/{r}/pulls/{n}/reviews

        Args:
            output_dir: Directory to write Parquet files
            start_date: Start date filter for commits (e.g., "90d", "2025-01-01")
            continue_on_error: If True, continue extraction even if some steps fail
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        # Step 1: Extract repos
        logger.info("Extracting repositories...")
        try:
            self.repos_extractor.extract(output_dir)
            repos_df = pd.read_parquet(output_dir / "repos.parquet")
            repo_names = repos_df["full_name"].tolist() if not repos_df.empty else []
            logger.info(f"Extracted {len(repo_names)} repositories")
        except Exception as e:
            logger.error(f"Failed to extract repos: {e}")
            if not continue_on_error:
                raise
            repo_names = []
            # Write empty parquet
            pd.DataFrame().to_parquet(output_dir / "repos.parquet", index=False)

        # Step 2: Extract commits
        logger.info("Extracting commits...")
        try:
            self.commits_extractor.extract(output_dir, start_date=start_date)
            commits_df = pd.read_parquet(output_dir / "commits.parquet")
            logger.info(f"Extracted {len(commits_df)} commits")
        except Exception as e:
            logger.error(f"Failed to extract commits: {e}")
            if not continue_on_error:
                raise
            pd.DataFrame().to_parquet(output_dir / "commits.parquet", index=False)

        # Step 3: Extract PRs for each repo
        logger.info("Extracting pull requests...")
        try:
            self.prs_extractor.extract(output_dir, repos=repo_names)
            prs_df = pd.read_parquet(output_dir / "pull_requests.parquet")
            logger.info(f"Extracted {len(prs_df)} pull requests")
        except Exception as e:
            logger.error(f"Failed to extract PRs: {e}")
            if not continue_on_error:
                raise
            prs_df = pd.DataFrame()
            pd.DataFrame().to_parquet(output_dir / "pull_requests.parquet", index=False)

        # Step 4: Extract reviews for each PR
        logger.info("Extracting reviews...")
        try:
            all_reviews = []

            if not prs_df.empty and "number" in prs_df.columns and "repo_name" in prs_df.columns:
                # Group PRs by repo
                pr_groups = prs_df.groupby("repo_name")["number"].apply(list).to_dict()

                for repo, pr_numbers in pr_groups.items():
                    logger.info(f"Extracting reviews for {repo} ({len(pr_numbers)} PRs)...")
                    # Use the reviews extractor directly but collect results
                    for pr_number in pr_numbers:
                        try:
                            df = self.reviews_extractor.fetch_github_style(
                                f"/repos/{repo}/pulls/{pr_number}/reviews"
                            )
                            if not df.empty:
                                df["repo_name"] = repo
                                df["pr_number"] = pr_number
                                all_reviews.append(df)
                        except Exception as e:
                            logger.warning(f"Failed to extract reviews for {repo}#{pr_number}: {e}")
                            if not continue_on_error:
                                raise

            if all_reviews:
                combined_df = pd.concat(all_reviews, ignore_index=True)
            else:
                combined_df = pd.DataFrame()

            combined_df.to_parquet(output_dir / "reviews.parquet", index=False)
            logger.info(f"Extracted {len(combined_df)} reviews")

        except Exception as e:
            logger.error(f"Failed to extract reviews: {e}")
            if not continue_on_error:
                raise
            pd.DataFrame().to_parquet(output_dir / "reviews.parquet", index=False)

        logger.info(f"Done! Files written to {output_dir}")


def create_parser() -> argparse.ArgumentParser:
    """
    Create argument parser for CLI.

    Returns:
        Configured ArgumentParser
    """
    parser = argparse.ArgumentParser(
        description="Extract data from cursor-sim API to Parquet files"
    )
    parser.add_argument(
        "--url",
        required=True,
        help="Base URL for cursor-sim API (e.g., http://localhost:8080)"
    )
    parser.add_argument(
        "--output",
        default="data/raw",
        help="Output directory for Parquet files (default: data/raw)"
    )
    parser.add_argument(
        "--api-key",
        default="cursor-sim-dev-key",
        help="API key for authentication (default: cursor-sim-dev-key)"
    )
    parser.add_argument(
        "--start-date",
        default="90d",
        help="Start date for commits filter (default: 90d)"
    )
    parser.add_argument(
        "--continue-on-error",
        action="store_true",
        help="Continue extraction even if some steps fail"
    )
    return parser


def main() -> None:
    """Main entry point for CLI."""
    parser = create_parser()
    args = parser.parse_args()

    loader = DataLoader(args.url, args.api_key)
    loader.run(
        Path(args.output),
        start_date=args.start_date,
        continue_on_error=args.continue_on_error
    )


if __name__ == "__main__":
    main()
