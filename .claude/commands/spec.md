# Display Service Specification

When the user runs `/spec [service-name]`, display the specification for the requested service.

## Services

- `cursor-sim` - services/cursor-sim/SPEC.md
- `cursor-analytics-core` - services/cursor-analytics-core/SPEC.md
- `cursor-viz-spa` - services/cursor-viz-spa/SPEC.md

## Behavior

1. If no service name provided, list available services
2. If service name provided, read and display the corresponding SPEC.md
3. Highlight key sections:
   - Overview
   - API Endpoints
   - Data Models
   - Configuration
   - Implementation Checklist

## Example Usage

```
User: /spec cursor-sim
Assistant: [Reads services/cursor-sim/SPEC.md and displays formatted summary]
```

## Implementation Instructions

Read the appropriate SPEC.md file and present it in a structured, easy-to-read format. Focus on actionable information for developers starting work on that service.
