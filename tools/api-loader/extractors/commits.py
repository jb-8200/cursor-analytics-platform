"""
Commits extractor for /analytics/ai-code/commits endpoint.

Extracts AI code tracking commit data from cursor-sim Cursor Analytics-style endpoint.
"""

from pathlib import Path
from typing import Optional
from extractors.base import BaseAPIExtractor


class CommitsExtractor(BaseAPIExtractor):
    """
    Extract commits from /analytics/ai-code/commits endpoint.

    Endpoint: GET /analytics/ai-code/commits
    Style: Cursor Analytics (wrapped with pagination metadata)
    Pagination: page/page_size with hasNextPage
    """

    def extract(self, output_dir: Path, start_date: str = "90d") -> None:
        """
        Extract commits and write to Parquet file.

        Args:
            output_dir: Directory to write commits.parquet
            start_date: Start date filter (e.g., "90d", "2025-01-01")
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        # Fetch commits from Cursor Analytics-style endpoint with pagination
        params = {"startDate": start_date}
        df = self.fetch_cursor_style_paginated("/analytics/ai-code/commits", params=params)

        # Write to parquet
        output_file = output_dir / "commits.parquet"
        self.write_parquet(df, output_file)
