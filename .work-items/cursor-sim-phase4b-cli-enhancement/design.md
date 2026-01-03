# Technical Design: Interactive CLI Configuration

**Feature ID**: cursor-sim-phase4b-cli-enhancement
**Created**: January 3, 2026
**Status**: Design (Ready to Implement)

**Note**: Phase 4A (Empty Dataset Fixes) is complete. This design is for Phase 4B only.

---

## Overview

This design implements interactive CLI prompts for controlling cursor-sim data generation parameters, replacing hardcoded values with user-configurable inputs that accept sensible defaults.

---

## Architecture

### High-Level Flow

```
┌────────────────────────────────────────────────────────────────┐
│                     cursor-sim Startup                          │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────────┐
│             Parse CLI Flags (-mode, -seed, -port)               │
│             (Non-interactive flags take precedence)             │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ Interactive Mode?│
                    │ (--interactive)  │
                    └─────────────────┘
                         YES │  NO
                    ┌────────┴────────┐
                    ▼                 ▼
       ┌─────────────────────┐   ┌──────────────────┐
       │ Show Prompts:       │   │ Use CLI Flags    │
       │ 1. Developers       │   │ or Defaults      │
       │ 2. Period (months)  │   └──────────────────┘
       │ 3. Max Commits      │            │
       └─────────────────────┘            │
                    │                     │
                    └──────────┬──────────┘
                               ▼
                    ┌─────────────────────┐
                    │ Validate Inputs     │
                    │ (Regex + Range)     │
                    └─────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │ Load Seed File      │
                    └─────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │ Replicate/Sample    │
                    │ Developers to Count │
                    └─────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │ Generate Commits    │
                    │ (with max limit)    │
                    └─────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │ Start HTTP Server   │
                    └─────────────────────┘
```

---

## Component Design

### 1. Interactive Prompt Module (`internal/config/interactive.go`)

**Responsibility**: Handle user input prompts with defaults and validation.

```go
package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// PromptConfig holds configuration for a single prompt
type PromptConfig struct {
	Label       string        // "Number of developers"
	DefaultVal  int           // Default value
	Pattern     *regexp.Regexp // Validation regex
	Min         int           // Minimum allowed value
	Max         int           // Maximum allowed value
	MaxRetries  int           // Max retry attempts (default 3)
}

// PromptForInt displays an interactive prompt and returns validated integer input
func PromptForInt(cfg PromptConfig) (int, error) {
	reader := bufio.NewReader(os.Stdin)

	for attempt := 0; attempt < cfg.MaxRetries; attempt++ {
		// Display prompt with default
		fmt.Printf("%s [default: %d]: ", cfg.Label, cfg.DefaultVal)

		// Read input
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Empty input = use default
		if input == "" {
			fmt.Printf("✓ Using default: %d\n", cfg.DefaultVal)
			return cfg.DefaultVal, nil
		}

		// Validate with regex
		if !cfg.Pattern.MatchString(input) {
			fmt.Printf("✗ Invalid input. Expected format: positive integer\n")
			continue
		}

		// Convert to int
		val, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("✗ Invalid number: %s\n", input)
			continue
		}

		// Range validation
		if val < cfg.Min || val > cfg.Max {
			fmt.Printf("✗ Value must be between %d and %d\n", cfg.Min, cfg.Max)
			continue
		}

		fmt.Printf("✓ %s: %d\n", cfg.Label, val)
		return val, nil
	}

	// Max retries exceeded, use default
	fmt.Printf("⚠ Max retries exceeded, using default: %d\n", cfg.DefaultVal)
	return cfg.DefaultVal, nil
}

// InteractiveConfig prompts user for all configuration values
func InteractiveConfig() (*GenerationParams, error) {
	fmt.Println("\nConfiguration (press Enter for defaults):\n")

	// Prompt 1: Number of developers
	developers, err := PromptForInt(PromptConfig{
		Label:      "Number of developers",
		DefaultVal: 2,
		Pattern:    regexp.MustCompile(`^\d+$`),
		Min:        1,
		Max:        100,
		MaxRetries: 3,
	})
	if err != nil {
		return nil, err
	}

	// Prompt 2: Period in months
	months, err := PromptForInt(PromptConfig{
		Label:      "Period in months",
		DefaultVal: 3,
		Pattern:    regexp.MustCompile(`^\d+$`),
		Min:        1,
		Max:        24,
		MaxRetries: 3,
	})
	if err != nil {
		return nil, err
	}

	// Prompt 3: Maximum commits
	maxCommits, err := PromptForInt(PromptConfig{
		Label:      "Maximum commits",
		DefaultVal: 1000,
		Pattern:    regexp.MustCompile(`^\d+$`),
		Min:        1,
		Max:        100000,
		MaxRetries: 3,
	})
	if err != nil {
		return nil, err
	}

	// Convert months to days
	days := months * 30

	// Display summary
	fmt.Println("\nValidating inputs...")
	fmt.Printf("✓ Developers: %d\n", developers)
	fmt.Printf("✓ Period: %d months (%d days)\n", months, days)
	fmt.Printf("✓ Max commits: %d\n\n", maxCommits)

	return &GenerationParams{
		Developers: developers,
		Days:       days,
		MaxCommits: maxCommits,
	}, nil
}

// GenerationParams holds the parsed generation parameters
type GenerationParams struct {
	Developers int
	Days       int
	MaxCommits int
}
```

---

### 2. Developer Replication Logic (`internal/seed/replicator.go`)

**Responsibility**: Replicate or sample from seed file developers to reach target count.

```go
package seed

import (
	"fmt"
	"math/rand"
)

// ReplicateDevelopers creates N developers by replicating/sampling from seed data
// Strategy:
//   - If N <= len(seed.Developers): Sample N developers randomly
//   - If N > len(seed.Developers): Cycle through seed developers, modifying IDs
func ReplicateDevelopers(seed *SeedData, targetCount int, rng *rand.Rand) ([]Developer, error) {
	if targetCount < 1 {
		return nil, fmt.Errorf("target count must be >= 1, got %d", targetCount)
	}

	if len(seed.Developers) == 0 {
		return nil, fmt.Errorf("seed data has no developers")
	}

	result := make([]Developer, 0, targetCount)

	// Case 1: Downsample (N <= seed count)
	if targetCount <= len(seed.Developers) {
		// Shuffle and take first N
		indices := rng.Perm(len(seed.Developers))
		for i := 0; i < targetCount; i++ {
			result = append(result, seed.Developers[indices[i]])
		}
		return result, nil
	}

	// Case 2: Replicate (N > seed count)
	for i := 0; i < targetCount; i++ {
		// Cycle through seed developers
		srcDev := seed.Developers[i % len(seed.Developers)]

		// Clone and modify ID/email to make unique
		clonedDev := srcDev
		clonedDev.UserID = fmt.Sprintf("%s_clone%d", srcDev.UserID, i/len(seed.Developers))
		clonedDev.Email = fmt.Sprintf("clone%d_%s", i/len(seed.Developers), srcDev.Email)
		clonedDev.Name = fmt.Sprintf("%s (Clone %d)", srcDev.Name, i/len(seed.Developers))

		result = append(result, clonedDev)
	}

	return result, nil
}
```

---

### 3. Commit Limit Logic (`internal/generator/commit_generator.go`)

**Modification**: Add max commit tracking to `GenerateCommits`.

```go
// GenerateCommits generates commits for the specified number of days.
// Stops early if maxCommits is reached (0 = unlimited).
func (g *CommitGenerator) GenerateCommits(ctx context.Context, days int, maxCommits int) error {
	startTime := time.Now().AddDate(0, 0, -days)

	totalCommits := 0

	for _, dev := range g.seed.Developers {
		if maxCommits > 0 && totalCommits >= maxCommits {
			log.Printf("Reached max commits (%d), stopping generation early\n", maxCommits)
			return nil
		}

		remaining := maxCommits - totalCommits
		if remaining <= 0 {
			remaining = -1 // unlimited
		}

		commits, err := g.generateForDeveloper(ctx, dev, startTime, remaining)
		if err != nil {
			return err
		}

		totalCommits += commits
	}

	return nil
}

// generateForDeveloper now returns commit count
func (g *CommitGenerator) generateForDeveloper(ctx context.Context, dev seed.Developer, startTime time.Time, maxCommits int) (int, error) {
	// ... existing logic ...

	count := 0
	for current.Before(now) {
		if maxCommits > 0 && count >= maxCommits {
			return count, nil // Stop early
		}

		// ... generate commit ...
		count++
	}

	return count, nil
}
```

---

### 4. Updated Config Struct (`internal/config/config.go`)

**Modification**: Add interactive flag and generation params.

```go
type Config struct {
	// Mode is the operation mode: "runtime" or "replay"
	Mode string

	// SeedPath is the path to seed.json (required for runtime mode)
	SeedPath string

	// Port is the HTTP server port
	Port int

	// Interactive mode: prompt for parameters
	Interactive bool

	// Generation parameters (from interactive mode or flags)
	GenParams GenerationParams

	// Velocity controls event generation rate: "low", "medium", or "high"
	Velocity string
}

// ParseFlags with new --interactive flag
func ParseFlags() (*Config, error) {
	// ... existing code ...

	fs.BoolVar(&cfg.Interactive, "interactive", false, "Interactive mode: prompt for generation parameters")

	// Add individual flags for non-interactive mode
	fs.IntVar(&cfg.GenParams.Developers, "developers", 0, "Number of developers (0 = use seed file count)")
	fs.IntVar(&cfg.GenParams.Days, "days", 90, "Days of history to generate")
	fs.IntVar(&cfg.GenParams.MaxCommits, "max-commits", 0, "Maximum commits to generate (0 = unlimited)")

	// ... parse and validate ...

	return cfg, nil
}
```

---

## Data Flow

### Input Validation Flow

```
User Input: "5"
     │
     ▼
┌─────────────────────┐
│ Regex Match?        │──NO──▶ "Invalid input. Expected: positive integer"
│ Pattern: ^\d+$      │
└─────────────────────┘
     │ YES
     ▼
┌─────────────────────┐
│ Convert to Int      │──FAIL─▶ "Invalid number"
└─────────────────────┘
     │ SUCCESS
     ▼
┌─────────────────────┐
│ Range Check?        │──NO──▶ "Value must be between 1 and 100"
│ (1 <= val <= 100)   │
└─────────────────────┘
     │ YES
     ▼
   ACCEPT
```

---

## Edge Cases

### Case 1: Max Commits Reached Before Period Ends

**Scenario**: User specifies 6 months but max 500 commits. Generation produces 500 commits in 2 months.

**Solution**: Log message: `Reached max commits (500) before end of period`

---

### Case 2: User Requests More Developers Than Seed Has

**Scenario**: Seed has 2 developers, user wants 10.

**Solution**: Replicate developers by cycling:
- Dev 0 (Alice) → alice, alice_clone1, alice_clone2, alice_clone3, alice_clone4
- Dev 1 (Bob) → bob, bob_clone1, bob_clone2, bob_clone3, bob_clone4

---

### Case 3: User Enters Invalid Input 3 Times

**Scenario**: User types "abc", "10x", "1.5" on prompts.

**Solution**: After 3 retries, use default value and log warning.

---

## Testing Strategy

### Unit Tests

| Test Case | Input | Expected Output |
|-----------|-------|-----------------|
| Valid int | "5" | 5 |
| Empty (default) | "" | Default value |
| Invalid chars | "abc" | Error, retry |
| Out of range | "200" (max 100) | Error, retry |
| Leading/trailing spaces | " 5 " | 5 (trimmed) |

### Integration Tests

| Test Case | Setup | Assertion |
|-----------|-------|-----------|
| Interactive mode | Pipe inputs: "3\n6\n1500\n" | 3 devs, 180 days, 1500 max commits |
| Non-interactive | CLI flags only | Uses flag values |
| Max commits reached | 2 devs, 90 days, max 10 commits | Exactly 10 commits generated |
| Developer replication | Seed: 2 devs, Request: 5 devs | 5 developers in output |

---

## Migration Plan

### Backward Compatibility

| Old Command | New Equivalent |
|-------------|----------------|
| `-days 90` | `-days 90` or interactive mode |
| (no dev control) | `-developers N` or interactive mode |
| (no commit limit) | `-max-commits N` or interactive mode |

**No breaking changes**: Existing CLI flags continue to work.

---

## Performance Impact

| Metric | Impact | Justification |
|--------|--------|---------------|
| Startup time | +0.5s | Interactive prompts |
| Generation time | No change | Same algorithm |
| Memory usage | No change | Same data structures |

---

## Alternative Approaches

### Alternative 1: Config File (YAML/JSON)

**Pros**: Reproducible, version-controllable
**Cons**: Extra file management, less discoverable
**Decision**: Deferred to Phase 5

### Alternative 2: Environment Variables Only

**Pros**: CI/CD friendly
**Cons**: Poor UX for local development
**Decision**: Keep env vars, add interactive mode

### Alternative 3: GUI (Terminal UI)

**Pros**: Rich UX, live validation
**Cons**: Complexity, dependencies (bubbletea, etc.)
**Decision**: Out of scope for MVP

---

## Security Considerations

1. **Input Injection**: Regex validation prevents command injection
2. **Resource Exhaustion**: Max values capped (100 devs, 24 months, 100k commits)
3. **File System Access**: No new file write operations

---

## Open Questions

1. **Q**: Should we persist the last-used values?
   **A**: Deferred to Phase 5 (nice-to-have)

2. **Q**: Should we support loading parameters from stdin (for scripting)?
   **A**: Yes, pipe support: `echo -e "5\n6\n1500\n" | ./cursor-sim -interactive`

---

## Related Files

| File | Change Type |
|------|-------------|
| `internal/config/interactive.go` | NEW |
| `internal/config/config.go` | MODIFY (add Interactive flag) |
| `internal/seed/replicator.go` | NEW |
| `internal/generator/commit_generator.go` | MODIFY (add max commits) |
| `cmd/simulator/main.go` | MODIFY (call interactive mode) |

---

**Next Steps**: Task breakdown → TDD implementation
