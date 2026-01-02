#!/usr/bin/env python3
"""
Markdown formatter hook for Claude Code.

Configured as PostToolUse hook with Edit|Write matcher.
Auto-formats markdown files after editing.

Setup: /hooks → PostToolUse → Edit|Write → this script
"""
import json
import sys
import re
import os

def detect_language(code):
    """Best-effort language detection from code content."""
    s = code.strip()

    # JSON detection
    if re.search(r'^\s*[{\[]', s):
        try:
            json.loads(s)
            return 'json'
        except:
            pass

    # Go detection
    if re.search(r'^\s*package\s+\w+', s, re.M) or \
       re.search(r'^\s*func\s+', s, re.M) or \
       re.search(r'^\s*type\s+\w+\s+struct', s, re.M):
        return 'go'

    # Python detection
    if re.search(r'^\s*def\s+\w+\s*\(', s, re.M) or \
       re.search(r'^\s*(import|from)\s+\w+', s, re.M):
        return 'python'

    # JavaScript/TypeScript detection
    if re.search(r'\b(function\s+\w+\s*\(|const\s+\w+\s*=)', s) or \
       re.search(r'=>|console\.(log|error)', s):
        return 'javascript'

    # Bash detection
    if re.search(r'^#!.*\b(bash|sh)\b', s, re.M) or \
       re.search(r'^\s*(if|then|fi|for|in|do|done)\b', s, re.M) or \
       re.search(r'^\s*(git|npm|go|make|curl)\s+', s, re.M):
        return 'bash'

    # SQL detection
    if re.search(r'\b(SELECT|INSERT|UPDATE|DELETE|CREATE)\s+', s, re.I):
        return 'sql'

    # Markdown detection (nested)
    if re.search(r'^#+\s+', s, re.M):
        return 'markdown'

    return 'text'

def format_markdown(content):
    """Format markdown content with language detection."""

    def add_lang_to_fence(match):
        indent, info, body, closing = match.groups()
        if not info.strip():
            lang = detect_language(body)
            return f"{indent}```{lang}\n{body}{closing}\n"
        return match.group(0)

    # Fix unlabeled code fences
    fence_pattern = r'(?ms)^([ \t]{0,3})```([^\n]*)\n(.*?)(\n\1```)\s*'
    content = re.sub(fence_pattern, add_lang_to_fence, content)

    # Fix excessive blank lines (more than 2)
    content = re.sub(r'\n{4,}', '\n\n\n', content)

    return content.rstrip() + '\n'

def main():
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError:
        sys.exit(0)

    tool_input = data.get('tool_input', {})
    file_path = tool_input.get('file_path', '')

    # Only process markdown files
    if not file_path.endswith(('.md', '.mdx', '.mdc')):
        sys.exit(0)

    if os.path.exists(file_path):
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            formatted = format_markdown(content)

            if formatted != content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(formatted)
                print(f"✓ Formatted: {file_path}")

        except Exception as e:
            print(f"Warning: Could not format {file_path}: {e}", file=sys.stderr)

    sys.exit(0)

if __name__ == '__main__':
    main()
