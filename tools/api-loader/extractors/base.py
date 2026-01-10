"""
Base API extractor for cursor-sim endpoints.

Handles two response formats:
1. GitHub-style: Raw arrays (e.g., /repos, /repos/{o}/{r}/pulls)
2. Cursor Analytics-style: Wrapped objects with {data: [...], pagination: {...}}
"""

from pathlib import Path
from typing import Any, Dict, List, Optional
import requests
import pandas as pd


class BaseAPIExtractor:
    """
    Base class for extracting data from cursor-sim API.

    Supports both GitHub-style (raw arrays) and Cursor Analytics-style
    (wrapped with pagination metadata) endpoints.
    """

    def __init__(self, base_url: str, api_key: str = "cursor-sim-dev-key"):
        """
        Initialize the API extractor.

        Args:
            base_url: Base URL for cursor-sim API (e.g., http://localhost:8080)
            api_key: API key for authentication (default: cursor-sim-dev-key)
        """
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key
        self.auth = (api_key, '')  # Basic auth with empty password

    def fetch_github_style(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> pd.DataFrame:
        """
        Fetch from GitHub-style endpoint that returns raw array.

        These endpoints return arrays directly without wrapper objects:
        - GET /repos -> [{full_name: "...", ...}, ...]
        - GET /repos/{o}/{r} -> {full_name: "...", ...}

        Args:
            endpoint: API endpoint path (e.g., /repos)
            params: Optional query parameters

        Returns:
            DataFrame with response data
        """
        url = f"{self.base_url}{endpoint}"
        resp = requests.get(url, params=params, auth=self.auth)
        resp.raise_for_status()

        data = resp.json()

        # Handle both single object and array responses
        if isinstance(data, dict):
            data = [data]
        elif not isinstance(data, list):
            data = []

        return pd.DataFrame(data)

    def fetch_github_style_paginated(
        self,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        per_page: int = 100
    ) -> pd.DataFrame:
        """
        Fetch from GitHub-style endpoint with pagination.

        GitHub-style pagination uses:
        - Query params: page, per_page
        - Termination: Empty array response

        Args:
            endpoint: API endpoint path
            params: Optional query parameters
            per_page: Items per page (default 100)

        Returns:
            DataFrame with all paginated results
        """
        all_data = []
        page = 1

        if params is None:
            params = {}

        while True:
            page_params = {**params, "page": page, "per_page": per_page}
            url = f"{self.base_url}{endpoint}"

            resp = requests.get(url, params=page_params, auth=self.auth)
            resp.raise_for_status()

            data = resp.json()

            if not isinstance(data, list):
                raise ValueError(f"Expected array response, got {type(data)}")

            # Empty array signals end of pagination
            if not data:
                break

            all_data.extend(data)

            # If we got fewer items than requested, we're done
            if len(data) < per_page:
                break

            page += 1

        return pd.DataFrame(all_data)

    def fetch_cursor_style(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> pd.DataFrame:
        """
        Fetch from Cursor Analytics-style endpoint with wrapped response.

        These endpoints return:
        {
            "data": [...],
            "pagination": {"page": 1, "pageSize": 100, "hasNextPage": false, ...},
            "params": {...}
        }

        Args:
            endpoint: API endpoint path (e.g., /analytics/ai-code/commits)
            params: Optional query parameters

        Returns:
            DataFrame with data from response
        """
        url = f"{self.base_url}{endpoint}"
        resp = requests.get(url, params=params, auth=self.auth)
        resp.raise_for_status()

        response = resp.json()

        # Extract data array from wrapper
        data = response.get("data", [])

        return pd.DataFrame(data)

    def fetch_cursor_style_paginated(
        self,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        page_size: int = 500
    ) -> pd.DataFrame:
        """
        Fetch from Cursor Analytics-style endpoint with pagination.

        Handles two cursor-sim response formats:
        1. {data: [...], pagination: {hasNextPage: bool}}
        2. {items: [...], totalCount: int, page: int, pageSize: int}

        Args:
            endpoint: API endpoint path
            params: Optional query parameters
            page_size: Items per page (default 500, max allowed by cursor-sim)

        Returns:
            DataFrame with all paginated results
        """
        all_data = []
        page = 1

        if params is None:
            params = {}

        while True:
            page_params = {**params, "page": page, "page_size": page_size}
            url = f"{self.base_url}{endpoint}"

            resp = requests.get(url, params=page_params, auth=self.auth)
            resp.raise_for_status()

            response = resp.json()

            # Handle both response formats
            # Format 1: {data: [...], pagination: {hasNextPage}}
            # Format 2: {items: [...], totalCount, page, pageSize}
            if "items" in response:
                # cursor-sim actual format
                data = response.get("items", [])
                total_count = response.get("totalCount", 0)
                current_page = response.get("page", page)
                current_page_size = response.get("pageSize", page_size)

                all_data.extend(data)

                # Check if we got all items or if there's more
                if len(data) < current_page_size:
                    # Got fewer items than requested - we're done
                    break
                if len(all_data) >= total_count:
                    # We've fetched all available items
                    break
            else:
                # Standard format: {data: [...], pagination: {hasNextPage}}
                data = response.get("data", [])
                all_data.extend(data)

                # Check pagination metadata
                pagination = response.get("pagination", {})
                has_next = pagination.get("hasNextPage", False)

                if not has_next:
                    break

            page += 1

        return pd.DataFrame(all_data)

    def write_parquet(self, df: pd.DataFrame, output_path: Path) -> None:
        """
        Write DataFrame to Parquet file.

        Args:
            df: DataFrame to write
            output_path: Path to output Parquet file
        """
        output_path = Path(output_path)
        output_path.parent.mkdir(parents=True, exist_ok=True)

        df.to_parquet(output_path, index=False, engine='pyarrow')
