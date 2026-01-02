#!/usr/bin/env python3
"""
Cursor Sim Seed Generator

Uses NVIDIA NeMo DataDesigner to generate seed data for cursor-sim.
This script produces seed.json containing:
- Developer roster with org structure and behavioral parameters
- Repository catalog with languages and file patterns
- Text templates (commit messages, PR titles, chat themes)
- Pre-computed correlation parameters

DataDesigner handles: realistic names, LLM-generated text, correlated attributes
cursor-sim handles: time-series events, API contract, pagination, aggregations
"""

from __future__ import annotations

import argparse
import json
import os
import uuid
from datetime import datetime
from pathlib import Path
from typing import Any

# Try to import DataDesigner, fall back to simple generation if not available
try:
    from data_designer.essentials import (
        CategorySamplerParams,
        DataDesigner,
        DataDesignerConfigBuilder,
        ExpressionColumnConfig,
        InferenceParameters,
        LLMTextColumnConfig,
        ModelConfig,
        ModelProvider,
        PersonSamplerParams,
        SamplerColumnConfig,
        SamplerType,
        ScipySamplerParams,
        SubcategorySamplerParams,
    )
    DATADESIGNER_AVAILABLE = True
except ImportError:
    DATADESIGNER_AVAILABLE = False
    print("Warning: DataDesigner not installed. Using fallback generator.")


# =============================================================================
# Configuration
# =============================================================================

DEFAULT_CONFIG = {
    "num_developers": 100,
    "num_repos": 20,
    "orgs": ["acme-corp", "techstartup-io", "enterprise-inc"],
    "output_path": "seed.json",
}

# Organization structure
ORG_STRUCTURE = {
    "acme-corp": {
        "Engineering": ["Backend", "Frontend", "Mobile", "DevOps"],
        "Product": ["Design", "PM"],
        "Data": ["Data Engineering", "ML/AI", "Analytics"],
    },
    "techstartup-io": {
        "Core Platform": ["API", "Infrastructure", "Security"],
        "Growth": ["Web", "Mobile", "Experiments"],
    },
    "enterprise-inc": {
        "R&D": ["Research", "Innovation Lab"],
        "Platform Engineering": ["Core", "Integrations", "SRE"],
        "IT": ["Internal Tools", "Support Systems"],
    },
}

# Team to language mapping
TEAM_LANGUAGES = {
    "Backend": ["Python", "Go", "Java", "Rust", "Node.js"],
    "Frontend": ["TypeScript", "JavaScript"],
    "Mobile": ["Swift", "Kotlin", "React Native"],
    "DevOps": ["Python", "Go", "Bash", "Terraform"],
    "ML/AI": ["Python", "Julia"],
    "Data Engineering": ["Python", "Scala", "SQL"],
    "API": ["Go", "Python", "Node.js"],
    "Infrastructure": ["Go", "Terraform", "Python"],
    "Security": ["Python", "Go", "Rust"],
    "SRE": ["Go", "Python", "Bash"],
    "default": ["Python", "TypeScript"],
}

# Region configuration
REGIONS = {
    "US": {
        "weight": 0.50,
        "timezones": ["America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles"],
        "locales": ["en-US"],
        "working_hours": {"start": 9, "end": 18, "peak": 14},
    },
    "EU": {
        "weight": 0.30,
        "timezones": ["Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Warsaw"],
        "locales": ["en-GB", "de-DE", "fr-FR"],
        "working_hours": {"start": 8, "end": 17, "peak": 11},
    },
    "APAC": {
        "weight": 0.20,
        "timezones": ["Asia/Tokyo", "Asia/Singapore", "Asia/Shanghai", "Australia/Sydney"],
        "locales": ["ja-JP", "zh-CN", "en-AU"],
        "working_hours": {"start": 9, "end": 19, "peak": 15},
    },
}

# Seniority to behavior mapping
SENIORITY_BEHAVIOR = {
    "junior": {
        "weight": 0.20,
        "acceptance_rate": {"min": 0.45, "max": 0.65},
        "events_per_day": {"mean": 80, "std": 30},
        "tab_to_composer_ratio": 0.85,
        # PR Behavior (NEW for research)
        "prs_per_week": {"mean": 2.5, "std": 1.2},
        "avg_pr_size_loc": {"mean": 80, "std": 50},
        "avg_files_per_pr": {"mean": 3, "std": 2},
        "coding_speed_hours": {"mean": 8, "std": 4},
        "review_thoroughness": 0.6,
        "revert_rate_modifier": 2.0,
        "hotfix_rate_modifier": 1.8,
        "code_survival_modifier": 0.92,
    },
    "mid": {
        "weight": 0.40,
        "acceptance_rate": {"min": 0.60, "max": 0.80},
        "events_per_day": {"mean": 120, "std": 40},
        "tab_to_composer_ratio": 0.75,
        "prs_per_week": {"mean": 4.0, "std": 1.5},
        "avg_pr_size_loc": {"mean": 150, "std": 80},
        "avg_files_per_pr": {"mean": 5, "std": 3},
        "coding_speed_hours": {"mean": 4, "std": 2},
        "review_thoroughness": 0.75,
        "revert_rate_modifier": 1.2,
        "hotfix_rate_modifier": 1.2,
        "code_survival_modifier": 0.96,
    },
    "senior": {
        "weight": 0.25,
        "acceptance_rate": {"min": 0.75, "max": 0.92},
        "events_per_day": {"mean": 150, "std": 50},
        "tab_to_composer_ratio": 0.65,
        "prs_per_week": {"mean": 5.0, "std": 2.0},
        "avg_pr_size_loc": {"mean": 200, "std": 100},
        "avg_files_per_pr": {"mean": 6, "std": 4},
        "coding_speed_hours": {"mean": 2, "std": 1},
        "review_thoroughness": 0.85,
        "revert_rate_modifier": 0.7,
        "hotfix_rate_modifier": 0.8,
        "code_survival_modifier": 1.02,
    },
    "staff": {
        "weight": 0.10,
        "acceptance_rate": {"min": 0.80, "max": 0.95},
        "events_per_day": {"mean": 130, "std": 45},
        "tab_to_composer_ratio": 0.55,
        "prs_per_week": {"mean": 3.5, "std": 1.5},
        "avg_pr_size_loc": {"mean": 250, "std": 120},
        "avg_files_per_pr": {"mean": 8, "std": 5},
        "coding_speed_hours": {"mean": 3, "std": 1.5},
        "review_thoroughness": 0.90,
        "revert_rate_modifier": 0.5,
        "hotfix_rate_modifier": 0.7,
        "code_survival_modifier": 1.04,
    },
    "principal": {
        "weight": 0.05,
        "acceptance_rate": {"min": 0.82, "max": 0.96},
        "events_per_day": {"mean": 100, "std": 40},
        "tab_to_composer_ratio": 0.50,
        "prs_per_week": {"mean": 2.0, "std": 1.0},
        "avg_pr_size_loc": {"mean": 300, "std": 150},
        "avg_files_per_pr": {"mean": 10, "std": 6},
        "coding_speed_hours": {"mean": 4, "std": 2},
        "review_thoroughness": 0.95,
        "revert_rate_modifier": 0.4,
        "hotfix_rate_modifier": 0.6,
        "code_survival_modifier": 1.05,
    },
}

# AI ratio impact on quality (for research correlations)
AI_RATIO_MODIFIERS = {
    "low": {  # < 30% AI
        "review_iterations_modifier": 0.9,
        "review_density_modifier": 0.8,
        "revert_rate_modifier": 0.8,
        "code_survival_modifier": 1.02,
    },
    "medium": {  # 30-60% AI
        "review_iterations_modifier": 1.0,
        "review_density_modifier": 1.0,
        "revert_rate_modifier": 1.0,
        "code_survival_modifier": 1.0,
    },
    "high": {  # > 60% AI
        "review_iterations_modifier": 1.3,
        "review_density_modifier": 1.4,
        "revert_rate_modifier": 1.5,
        "code_survival_modifier": 0.95,
    },
}

# Cycle time parameters
CYCLE_TIME_PARAMS = {
    "pickup_time": {"mean": 6, "std": 4},  # hours
    "review_lead_time": {"mean": 8, "std": 6},  # hours
    "base_revert_rate": 0.02,  # 2% baseline
    "base_hotfix_rate": 0.08,  # 8% baseline
    "base_survival_rate": 0.85,  # 85% survive 30 days
}


# =============================================================================
# DataDesigner-based Generator
# =============================================================================

def generate_with_datadesigner(
    num_developers: int,
    num_repos: int,
    api_key: str | None = None,
) -> dict[str, Any]:
    """Generate seed data using NVIDIA NeMo DataDesigner."""
    
    # Setup model provider
    providers = []
    if api_key or os.getenv("NVIDIA_API_KEY"):
        providers.append(
            ModelProvider(
                name="nvidia",
                endpoint="https://ai.api.nvidia.com/v1",
                provider_type="openai",
                api_key=api_key or "NVIDIA_API_KEY",
            )
        )
    
    model_configs = [
        ModelConfig(
            alias="text-gen",
            model="nvidia/nemotron-3-nano-30b-a3b",
            inference_parameters=InferenceParameters(
                temperature=0.7,
                max_tokens=256,
            ),
        ),
    ]
    
    data_designer = DataDesigner(model_providers=providers if providers else None)
    
    # =========================================================================
    # Generate Developers
    # =========================================================================
    dev_builder = DataDesignerConfigBuilder(model_configs=model_configs)
    
    # Person sampler for realistic names
    dev_builder.add_column(
        SamplerColumnConfig(
            name="person",
            sampler_type=SamplerType.PERSON,
            params=PersonSamplerParams(include_demographics=True),
        )
    )
    
    # Generate stable user ID
    dev_builder.add_column(
        ExpressionColumnConfig(
            name="user_id",
            expression="'user_' + ''.join(random.choices('abcdefghijklmnopqrstuvwxyz0123456789', k=12))",
            imports=["random"],
        )
    )
    
    # Organization
    dev_builder.add_column(
        SamplerColumnConfig(
            name="org",
            sampler_type=SamplerType.CATEGORY,
            params=CategorySamplerParams(
                values=list(ORG_STRUCTURE.keys()),
            ),
        )
    )
    
    # Region with weights
    dev_builder.add_column(
        SamplerColumnConfig(
            name="region",
            sampler_type=SamplerType.CATEGORY,
            params=CategorySamplerParams(
                values=list(REGIONS.keys()),
                weights=[r["weight"] for r in REGIONS.values()],
            ),
        )
    )
    
    # Seniority with weights
    dev_builder.add_column(
        SamplerColumnConfig(
            name="seniority",
            sampler_type=SamplerType.CATEGORY,
            params=CategorySamplerParams(
                values=list(SENIORITY_BEHAVIOR.keys()),
                weights=[s["weight"] for s in SENIORITY_BEHAVIOR.values()],
            ),
        )
    )
    
    # Generate developers
    dev_result = data_designer.generate(
        config_builder=dev_builder,
        num_records=num_developers,
    )
    developers_raw = dev_result.to_pandas().to_dict(orient="records")
    
    # Post-process developers
    developers = []
    for dev in developers_raw:
        # Pick division and team based on org
        org = dev.get("org", "acme-corp")
        divisions = list(ORG_STRUCTURE.get(org, ORG_STRUCTURE["acme-corp"]).keys())
        division = divisions[hash(dev["user_id"]) % len(divisions)]
        teams = ORG_STRUCTURE[org][division]
        team = teams[hash(dev["user_id"] + "team") % len(teams)]
        
        # Get region config
        region = dev.get("region", "US")
        region_config = REGIONS.get(region, REGIONS["US"])
        timezone = region_config["timezones"][hash(dev["user_id"] + "tz") % len(region_config["timezones"])]
        locale = region_config["locales"][hash(dev["user_id"] + "loc") % len(region_config["locales"])]
        
        # Get seniority config
        seniority = dev.get("seniority", "mid")
        seniority_config = SENIORITY_BEHAVIOR.get(seniority, SENIORITY_BEHAVIOR["mid"])
        
        # Calculate acceptance rate within band
        import random
        random.seed(hash(dev["user_id"] + "acc"))
        acceptance_rate = random.uniform(
            seniority_config["acceptance_rate"]["min"],
            seniority_config["acceptance_rate"]["max"],
        )
        
        developers.append({
            "user_id": dev["user_id"],
            "email": f"{dev['person'].first_name.lower()}.{dev['person'].last_name.lower()}@{org.replace('-', '')}.com"
                     if hasattr(dev.get("person"), "first_name") 
                     else f"dev{hash(dev['user_id']) % 10000}@example.com",
            "name": f"{dev['person'].first_name} {dev['person'].last_name}"
                    if hasattr(dev.get("person"), "first_name")
                    else f"Developer {hash(dev['user_id']) % 10000}",
            "role": "member" if random.random() > 0.15 else ("owner" if random.random() > 0.3 else "free-owner"),
            "org": org,
            "division": division,
            "team": team,
            "region": region,
            "timezone": timezone,
            "locale": locale,
            "seniority": seniority,
            "params": {
                "acceptance_rate": round(acceptance_rate, 3),
                "events_per_day": seniority_config["events_per_day"],
                "tab_to_composer_ratio": seniority_config["tab_to_composer_ratio"],
                "working_hours": region_config["working_hours"],
            },
        })
    
    # =========================================================================
    # Generate Repositories (simplified without LLM for now)
    # =========================================================================
    repositories = generate_repositories_fallback(num_repos, developers)
    
    # =========================================================================
    # Generate Text Templates
    # =========================================================================
    text_templates = generate_text_templates_fallback()
    
    return {
        "version": "1.0",
        "generated_at": datetime.utcnow().isoformat() + "Z",
        "generator": "datadesigner",
        "developers": developers,
        "repositories": repositories,
        "text_templates": text_templates,
        "correlations": {
            "seniority_behavior": SENIORITY_BEHAVIOR,
            "region_activity": {k: {"working_hours": v["working_hours"]} for k, v in REGIONS.items()},
        },
    }


# =============================================================================
# Fallback Generator (no DataDesigner)
# =============================================================================

def generate_fallback(
    num_developers: int,
    num_repos: int,
) -> dict[str, Any]:
    """Generate seed data without DataDesigner (using basic random generation)."""
    import random
    
    # Simple name lists for fallback
    first_names = ["Alex", "Jordan", "Taylor", "Morgan", "Casey", "Riley", "Quinn", "Avery",
                   "Cameron", "Drew", "Sage", "Reese", "Parker", "Blake", "Charlie", "Jamie",
                   "Sam", "Pat", "Chris", "Dana", "Kelly", "Lee", "Robin", "Terry"]
    last_names = ["Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
                  "Rodriguez", "Martinez", "Anderson", "Taylor", "Thomas", "Moore", "Jackson",
                  "Martin", "Lee", "Thompson", "White", "Harris", "Clark", "Lewis", "Walker"]
    
    developers = []
    for i in range(num_developers):
        # Generate stable random values using index as seed
        random.seed(i)
        
        user_id = f"user_{''.join(random.choices('abcdefghijklmnopqrstuvwxyz0123456789', k=12))}"
        first_name = random.choice(first_names)
        last_name = random.choice(last_names)
        
        # Pick org, division, team
        org = random.choice(list(ORG_STRUCTURE.keys()))
        division = random.choice(list(ORG_STRUCTURE[org].keys()))
        team = random.choice(ORG_STRUCTURE[org][division])
        
        # Pick region and derived fields
        region = random.choices(
            list(REGIONS.keys()),
            weights=[r["weight"] for r in REGIONS.values()],
        )[0]
        region_config = REGIONS[region]
        timezone = random.choice(region_config["timezones"])
        locale = random.choice(region_config["locales"])
        
        # Pick seniority and derived params
        seniority = random.choices(
            list(SENIORITY_BEHAVIOR.keys()),
            weights=[s["weight"] for s in SENIORITY_BEHAVIOR.values()],
        )[0]
        seniority_config = SENIORITY_BEHAVIOR[seniority]
        acceptance_rate = random.uniform(
            seniority_config["acceptance_rate"]["min"],
            seniority_config["acceptance_rate"]["max"],
        )
        
        developers.append({
            "user_id": user_id,
            "email": f"{first_name.lower()}.{last_name.lower()}{i}@{org.replace('-', '')}.com",
            "name": f"{first_name} {last_name}",
            "role": random.choices(["member", "owner", "free-owner"], weights=[0.85, 0.10, 0.05])[0],
            "org": org,
            "division": division,
            "team": team,
            "region": region,
            "timezone": timezone,
            "locale": locale,
            "seniority": seniority,
            "params": {
                # AI behavior
                "acceptance_rate": round(acceptance_rate, 3),
                "events_per_day": seniority_config["events_per_day"],
                "tab_to_composer_ratio": seniority_config["tab_to_composer_ratio"],
                "working_hours": region_config["working_hours"],
                # PR behavior (for research)
                "prs_per_week": seniority_config["prs_per_week"],
                "avg_pr_size_loc": seniority_config["avg_pr_size_loc"],
                "avg_files_per_pr": seniority_config["avg_files_per_pr"],
                "coding_speed_hours": seniority_config["coding_speed_hours"],
                "review_thoroughness": seniority_config["review_thoroughness"],
                # Quality modifiers (for research correlations)
                "revert_rate_modifier": seniority_config["revert_rate_modifier"],
                "hotfix_rate_modifier": seniority_config["hotfix_rate_modifier"],
                "code_survival_modifier": seniority_config["code_survival_modifier"],
            },
        })
    
    repositories = generate_repositories_fallback(num_repos, developers)
    text_templates = generate_text_templates_fallback()
    
    return {
        "version": "1.0",
        "generated_at": datetime.utcnow().isoformat() + "Z",
        "generator": "fallback",
        "developers": developers,
        "repositories": repositories,
        "text_templates": text_templates,
        "correlations": {
            "seniority_behavior": SENIORITY_BEHAVIOR,
            "region_activity": {k: {"working_hours": v["working_hours"]} for k, v in REGIONS.items()},
            # Research correlations
            "ai_ratio_modifiers": AI_RATIO_MODIFIERS,
            "cycle_time_params": CYCLE_TIME_PARAMS,
        },
    }


def generate_repositories_fallback(num_repos: int, developers: list[dict]) -> list[dict]:
    """Generate repository catalog with maturity metrics for research."""
    import random
    
    repo_prefixes = ["api", "web", "mobile", "infra", "data", "ml", "auth", "payments",
                     "notifications", "search", "analytics", "dashboard", "admin", "core"]
    repo_suffixes = ["service", "app", "platform", "system", "gateway", "pipeline", "lib"]
    
    repositories = []
    orgs = list(set(d["org"] for d in developers))
    
    for i in range(num_repos):
        random.seed(i + 1000)
        org = random.choice(orgs)
        
        # Pick a team from this org's developers
        org_teams = list(set(d["team"] for d in developers if d["org"] == org))
        if not org_teams:
            org_teams = ["Backend"]
        team = random.choice(org_teams)
        
        # Get language for team
        languages = TEAM_LANGUAGES.get(team, TEAM_LANGUAGES["default"])
        primary_language = random.choice(languages)
        
        repo_name = f"{org}/{random.choice(repo_prefixes)}-{random.choice(repo_suffixes)}"
        
        # Maturity metrics (for research control variables)
        age_days = int(random.lognormvariate(5.5, 0.8))  # ~250 days median
        age_days = max(30, min(2000, age_days))
        total_commits = int(age_days * random.uniform(0.5, 2.0))
        total_prs = int(total_commits * random.uniform(0.3, 0.6))
        total_contributors = min(50, max(3, total_commits // 100))
        size_bytes = int(random.lognormvariate(12, 1.5))  # ~150KB median
        
        # Code quality baseline
        greenfield_ratio = random.betavariate(2, 5)  # Skewed low for mature repos
        
        repositories.append({
            "repo_id": f"repo_{uuid.uuid4().hex[:12]}",
            "repo_name": repo_name,
            "org": org,
            "primary_team": team,
            "primary_language": primary_language,
            "service_type": random.choice(["api", "web-app", "library", "infrastructure", "data-pipeline"]),
            "default_branch": random.choices(["main", "master"], weights=[0.85, 0.15])[0],
            # Maturity metrics (research control variables)
            "maturity": {
                "age_days": age_days,
                "total_commits": total_commits,
                "total_prs": total_prs,
                "total_contributors": total_contributors,
                "size_bytes": size_bytes,
            },
            # Quality baselines
            "quality_baseline": {
                "greenfield_ratio": round(greenfield_ratio, 3),
                "revert_rate_baseline": round(random.betavariate(1, 40), 4),  # ~2-3%
                "hotfix_rate_baseline": round(random.betavariate(2, 20), 4),  # ~8-10%
            },
        })
    
    return repositories


def generate_text_templates_fallback() -> dict[str, Any]:
    """Generate text templates for commit messages, PR titles, review comments, etc."""
    return {
        "commit_messages": {
            "feature": [
                "Add {feature} to {component}",
                "Implement {functionality}",
                "Create {component} module",
                "Introduce {feature} support",
            ],
            "bugfix": [
                "Fix {issue} in {component}",
                "Resolve {error} error",
                "Handle edge case in {function}",
                "Correct {behavior} behavior",
            ],
            "refactor": [
                "Refactor {component} for clarity",
                "Extract {pattern} into helper",
                "Simplify {logic} handling",
                "Clean up {area} code",
            ],
            "chore": [
                "Update dependencies",
                "Configure {tool}",
                "Add tests for {component}",
                "Update documentation",
            ],
        },
        "pr_titles": {
            "feature": "[Feature] {description}",
            "bugfix": "[Fix] {description}",
            "refactor": "[Refactor] {description}",
            "chore": "[Chore] {description}",
        },
        "pr_descriptions": {
            "template": "## Summary\n{summary}\n\n## Changes\n{changes}\n\n## Testing\n{testing}",
        },
        "review_comments": {
            "style": [
                "Consider using const here instead of let",
                "This variable name could be more descriptive",
                "Formatting nit: extra whitespace",
                "Consider extracting this to a named constant",
            ],
            "logic": [
                "What happens if this is null?",
                "This might cause a race condition",
                "Edge case: what about empty input?",
                "Should we handle the error case here?",
                "This condition might not cover all cases",
            ],
            "suggestion": [
                "Could we extract this into a helper function?",
                "Have you considered using the builder pattern here?",
                "This might be cleaner with async/await",
                "We have a utility for this in shared/utils",
            ],
            "approval": [
                "LGTM!",
                "Looks good to me",
                "Nice work!",
                "Ship it 🚀",
                "Approved with minor suggestions",
            ],
            "request_changes": [
                "Please add tests for this change",
                "This needs error handling",
                "Let's discuss the approach before merging",
                "Missing type annotations",
            ],
        },
        "chat_themes": [
            "implement function",
            "fix bug",
            "explain code",
            "optimize performance",
            "add tests",
            "refactor module",
            "create API endpoint",
            "debug issue",
        ],
    }


# =============================================================================
# Main Entry Point
# =============================================================================

def main():
    parser = argparse.ArgumentParser(
        description="Generate seed data for cursor-sim using DataDesigner",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Generate with DataDesigner (requires NVIDIA_API_KEY)
  python generate_seed.py --developers 100 --repos 20 -o seed.json
  
  # Generate with fallback (no API key needed)
  python generate_seed.py --developers 50 --repos 10 --fallback -o seed.json
  
  # Preview output
  python generate_seed.py --preview
        """,
    )
    parser.add_argument(
        "--developers", "-d",
        type=int,
        default=100,
        help="Number of developers to generate (default: 100)",
    )
    parser.add_argument(
        "--repos", "-r",
        type=int,
        default=20,
        help="Number of repositories to generate (default: 20)",
    )
    parser.add_argument(
        "--output", "-o",
        type=Path,
        default=Path("seed.json"),
        help="Output file path (default: seed.json)",
    )
    parser.add_argument(
        "--fallback",
        action="store_true",
        help="Use fallback generator (no DataDesigner/API key needed)",
    )
    parser.add_argument(
        "--preview",
        action="store_true",
        help="Generate preview (5 developers, 3 repos) and print to stdout",
    )
    parser.add_argument(
        "--api-key",
        type=str,
        help="NVIDIA API key (or set NVIDIA_API_KEY env var)",
    )
    
    args = parser.parse_args()
    
    if args.preview:
        args.developers = 5
        args.repos = 3
    
    print(f"Generating seed data: {args.developers} developers, {args.repos} repos")
    
    # Choose generator
    use_datadesigner = DATADESIGNER_AVAILABLE and not args.fallback
    
    if use_datadesigner:
        print("Using DataDesigner for generation...")
        seed_data = generate_with_datadesigner(
            num_developers=args.developers,
            num_repos=args.repos,
            api_key=args.api_key,
        )
    else:
        print("Using fallback generator...")
        seed_data = generate_fallback(
            num_developers=args.developers,
            num_repos=args.repos,
        )
    
    # Output
    if args.preview:
        print("\n" + "=" * 60)
        print("PREVIEW OUTPUT")
        print("=" * 60)
        print(json.dumps(seed_data, indent=2))
    else:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        with open(args.output, "w") as f:
            json.dump(seed_data, f, indent=2)
        print(f"\nSeed data written to: {args.output}")
        print(f"  - {len(seed_data['developers'])} developers")
        print(f"  - {len(seed_data['repositories'])} repositories")
        print(f"\nUse with cursor-sim:")
        print(f"  cursor-sim --seed {args.output}")


if __name__ == "__main__":
    main()
