"""
Pre-patch hook for Claude Code.

This hook runs before Claude applies a code patch.  It performs
syntactic checks on Markdown files using Vale and markdownlint to
ensure documentation meets quality standards.  If lint errors are
found, they are included as a comment in the patch to inform the
assistant.
"""

import subprocess
from pathlib import Path
from typing import List

def run_lint(file_paths: List[str]) -> str:
    messages = []
    for path in file_paths:
        p = Path(path)
        if p.suffix.lower() == ".md" or p.suffix.lower() == ".mdc":
            try:
                # run vale
                result = subprocess.run(["vale", path], capture_output=True, text=True)
                if result.stdout:
                    messages.append(result.stdout)
                # run markdownlint
                result2 = subprocess.run(["markdownlint", path], capture_output=True, text=True)
                if result2.stdout:
                    messages.append(result2.stdout)
            except FileNotFoundError:
                # Vale or markdownlint not installed; skip
                pass
    return "\n".join(messages)

def pre_patch_hook(files: List[str], diff: str, context: dict) -> str:
    lint_output = run_lint(files)
    if lint_output:
        # Prepend lint messages as comments to the diff
        comment = "\n".join([f"<!-- {line} -->" for line in lint_output.splitlines()])
        return comment + "\n" + diff
    return diff
