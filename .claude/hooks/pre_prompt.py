"""
Pre-prompt hook for Claude Code.

This hook ensures that Claude always has access to the core rules
before processing user input.  It inspects the current prompt and
prepends includes for the specification files if they are missing.
Claude Code calls this function before each user message.
"""

def pre_prompt_hook(prompt: str, context: dict) -> str:
    # Include project-specific development instructions and genai-specs rules.
    # DEVELOPMENT.md provides project context, tech stack, and workflow.
    # genai-specs rules provide core development standards.
    include_lines = [
        "@./.claude/DEVELOPMENT.md",
        "@./.claude/skill/process-01-core.mdc",
        "@./.claude/skill/process-02-project.mdc",
        "@./.claude/skill/process-03-development.mdc",
        "@./.claude/skill/process-04-operational.mdc",
        "@./.claude/skill/process-05-coding.mdc",
        "@./.claude/skill/standards-user-story.mdc",
        "@./.claude/skill/standards-design.mdc",
        "@./.claude/skill/standards-task.mdc",
        "@./.claude/skill/standards-architecture.mdc",
        "@./.claude/skill/standards-decision.mdc",
        "@./.claude/skill/standards-guidelines.mdc",
    ]
    if any(line in prompt for line in include_lines):
        # Context already loaded.
        return prompt
    header = "\n".join(include_lines) + "\n\n"
    return header + prompt
