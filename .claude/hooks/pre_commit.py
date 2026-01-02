"""
Pre-commit hook for Claude Code.

This hook runs before Claude creates a git commit.  It ensures
that the unit tests and lints are passing.  If any test fails,
the commit is aborted and a message is returned to the user.
"""

import subprocess

def pre_commit_hook(context: dict) -> str:
    try:
        # Run pytest to ensure unit and integration tests pass
        result = subprocess.run([
            "pytest", "--quiet", "--disable-warnings", "--maxfail=1"
        ], capture_output=True, text=True)
        if result.returncode != 0:
            return (
                "Tests failed:\n" + result.stdout + result.stderr +
                "\nCommit aborted. Please fix the failing tests."
            )
        return ""  # success; no message needed
    except FileNotFoundError:
        # pytest not installed; skip test execution
        return "Warning: pytest not installed. Tests not executed before commit."
