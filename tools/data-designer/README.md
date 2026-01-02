# Data Designer Seed Generator

This tool uses NVIDIA NeMo DataDesigner to generate realistic seed data for cursor-sim.

## What This Generates

| Category | Examples |
|----------|----------|
| **Developer roster** | Names, emails, seniority, org structure |
| **Repository catalog** | Repo names, languages, service types |
| **Text templates** | Commit messages, PR titles |
| **Static correlations** | Seniority → acceptance rate band |

## What cursor-sim Generates (Not This Tool)

- Event timelines and time-series data
- Per-commit AI/non-AI line counts
- Daily/hourly aggregates
- Cursor API-shaped responses

## Quick Start

### With DataDesigner (Recommended)

```bash
# Install DataDesigner
pip install data-designer

# Set API key (free at build.nvidia.com)
export NVIDIA_API_KEY="your-key"

# Generate seed data
python generate_seed.py --developers 100 --repos 20 -o ../../seed.json
```

### Without DataDesigner (Fallback)

```bash
# Generate with basic random data (no API key needed)
python generate_seed.py --developers 100 --repos 20 --fallback -o ../../seed.json
```

### Preview Mode

```bash
# Preview output (5 developers, 3 repos)
python generate_seed.py --preview
```

## Output Format

```json
{
  "version": "1.0",
  "generated_at": "2025-01-15T10:30:00Z",
  "generator": "datadesigner",
  "developers": [
    {
      "user_id": "user_abc123xyz",
      "email": "alex.chen@acmecorp.com",
      "name": "Alex Chen",
      "org": "acme-corp",
      "division": "Engineering",
      "team": "Backend",
      "region": "US",
      "timezone": "America/New_York",
      "seniority": "senior",
      "params": {
        "acceptance_rate": 0.847,
        "events_per_day": {"mean": 150, "std": 50},
        "tab_to_composer_ratio": 0.65,
        "working_hours": {"start": 9, "end": 18, "peak": 14}
      }
    }
  ],
  "repositories": [...],
  "text_templates": {...},
  "correlations": {...}
}
```

## Using with cursor-sim

```bash
# Generate seed (one time)
make seed-generate

# Start simulator with seed
cursor-sim --mode=runtime --seed=seed.json --port=8080
```

## Files

- `config/seed_schema.yaml` - Schema definition for seed data
- `generate_seed.py` - Main generator script
- `templates/` - LLM prompt templates (optional)

## Requirements

- Python 3.10+
- data-designer (optional, for LLM-enhanced generation)
- NVIDIA API key (only if using DataDesigner)
