#!/usr/bin/env python3
import json
import sys
import os

data = json.load(sys.stdin)
file_path = data.get('tool_input', {}).get('file_path', '')

# Only process markdown files
if file_path and (file_path.endswith('.md') or file_path.endswith('.mdx')):
    try:
        # Check if file exists and is readable
        if os.path.isfile(file_path):
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()

            # Basic markdown formatting (could be extended with more rules)
            # For now, just ensure consistent line endings
            formatted = content.replace('\r\n', '\n')

            # Write back if changed
            if formatted != content:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(formatted)
                print(f"Markdown formatted: {file_path}")
    except Exception as e:
        # Non-blocking error
        print(f"Warning: Could not format {file_path}: {e}", file=sys.stderr)

sys.exit(0)
