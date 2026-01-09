# Spawn Planning Subagent

Spawns a planning agent (Opus model) for research and design work.

## Usage

```
/subagent:planning {feature-description}
```

## Agent Configuration

- **Model**: Opus (for complex reasoning and research)
- **Type**: Planning and design
- **Output**: Work items (user-story.md, design.md, task.md)

## When to Use

- Designing new features that require API research
- Breaking down complex features into tasks
- Creating technical specifications
- Planning multi-service integrations

## Example

```
/subagent:planning P4-F04 Add external data source simulators (Harvey, Copilot, Qualtrics)
```

## Spawn Command

```
Task tool with:
- subagent_type: Plan
- model: opus
- prompt: Include feature context, research requirements, and SDD compliance
```
