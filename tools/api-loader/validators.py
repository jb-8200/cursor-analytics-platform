"""
Schema validation for cursor-sim API data.

Feature: P8-F01 Data Tier ETL
Task: TASK-P8-06 Schema Validation

Validates DataFrames against JSON schema definitions to ensure
data integrity during ETL processing.
"""

import json
from pathlib import Path
from typing import Union
import pandas as pd


class SchemaValidationError(Exception):
    """Raised when DataFrame fails schema validation."""
    pass


def validate_dataframe(df: pd.DataFrame, schema_path: Union[str, Path]) -> None:
    """
    Validate DataFrame against JSON schema.

    Args:
        df: pandas DataFrame to validate
        schema_path: Path to JSON schema file

    Raises:
        FileNotFoundError: If schema file doesn't exist
        ValueError: If schema is invalid JSON or missing required_columns key
        SchemaValidationError: If DataFrame is missing required columns
    """
    schema_path = Path(schema_path)

    # Check schema file exists
    if not schema_path.exists():
        raise FileNotFoundError(f"Schema file not found: {schema_path}")

    # Load schema
    try:
        with open(schema_path, 'r') as f:
            schema = json.load(f)
    except json.JSONDecodeError as e:
        raise ValueError(f"Failed to parse schema file: {e}")

    # Validate schema structure
    if "required_columns" not in schema:
        raise ValueError("Schema must contain 'required_columns' key")

    required_columns = schema["required_columns"]

    # Empty DataFrames pass validation (no data to validate)
    if df.empty:
        return

    # Check for missing columns
    df_columns = set(df.columns)
    required_set = set(required_columns)
    missing = required_set - df_columns

    if missing:
        missing_list = sorted(missing)
        raise SchemaValidationError(
            f"Missing required columns: {', '.join(missing_list)}"
        )
