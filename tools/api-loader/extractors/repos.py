"""
Repos extractor for /repos endpoint.

Extracts repository metadata from cursor-sim GitHub-style endpoint.
"""

from pathlib import Path
from extractors.base import BaseAPIExtractor


class ReposExtractor(BaseAPIExtractor):
    """
    Extract repositories from /repos endpoint.

    Endpoint: GET /repos
    Style: GitHub (raw array)
    Pagination: Not paginated
    """

    def extract(self, output_dir: Path) -> None:
        """
        Extract repositories and write to Parquet file.

        Args:
            output_dir: Directory to write repos.parquet
        """
        output_dir = Path(output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)

        # Fetch repos from GitHub-style endpoint (returns raw array)
        df = self.fetch_github_style("/repos")

        # Write to parquet
        output_file = output_dir / "repos.parquet"
        self.write_parquet(df, output_file)
