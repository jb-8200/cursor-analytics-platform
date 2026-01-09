"""
dbt runner utilities for embedded pipeline.

Feature: P9-F01 Streamlit Dashboard
Task: TASK-P9-10 Refresh Pipeline

This module provides utilities to run dbt commands.
"""

import subprocess
from typing import Optional, Dict


def run_dbt_build(
    target: str = "dev",
    dbt_dir: str = "/app/dbt",
    select: Optional[str] = None
) -> tuple[bool, str, str]:
    """
    Run dbt build command.

    Args:
        target: dbt target (dev or prod)
        dbt_dir: Path to dbt project directory
        select: Optional dbt selector (e.g., "tag:daily")

    Returns:
        Tuple of (success, stdout, stderr)
    """
    cmd = ["dbt", "build", "--target", target]

    if select:
        cmd.extend(["--select", select])

    try:
        result = subprocess.run(
            cmd,
            cwd=dbt_dir,
            capture_output=True,
            text=True,
            check=True,
        )
        return True, result.stdout, result.stderr

    except subprocess.CalledProcessError as e:
        return False, e.stdout or "", e.stderr or ""


def run_dbt_test(
    target: str = "dev",
    dbt_dir: str = "/app/dbt"
) -> tuple[bool, str, str]:
    """
    Run dbt test command.

    Args:
        target: dbt target (dev or prod)
        dbt_dir: Path to dbt project directory

    Returns:
        Tuple of (success, stdout, stderr)
    """
    cmd = ["dbt", "test", "--target", target]

    try:
        result = subprocess.run(
            cmd,
            cwd=dbt_dir,
            capture_output=True,
            text=True,
            check=True,
        )
        return True, result.stdout, result.stderr

    except subprocess.CalledProcessError as e:
        return False, e.stdout or "", e.stderr or ""


def get_dbt_version() -> Optional[str]:
    """
    Get dbt version.

    Returns:
        dbt version string or None if dbt not found
    """
    try:
        result = subprocess.run(
            ["dbt", "--version"],
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout.strip()
    except (subprocess.CalledProcessError, FileNotFoundError):
        return None
