# Step 01: Initialize v2 Project Structure

## Objective

Create the v2 project structure by archiving v1 and setting up the new directory layout.

## Estimated Time: 1 hour
## Recommended Model: Haiku

## Prerequisites

- None (first step)

## Acceptance Criteria

- [ ] v1 code archived to `services/cursor-sim-v1/`
- [ ] New v2 directory created at `services/cursor-sim/`
- [ ] go.mod initialized with correct module path
- [ ] Directory structure matches design.md
- [ ] Makefile copied and updated
- [ ] `go build ./...` succeeds
- [ ] `golangci-lint run` passes

## Implementation Steps

### 1. Archive v1

```bash
mv services/cursor-sim services/cursor-sim-v1
```

### 2. Create v2 Directory Structure

```bash
mkdir -p services/cursor-sim/{cmd/simulator,internal/{config,seed,models,generator,storage,api/cursor},testdata}
```

### 3. Initialize Go Module

```bash
cd services/cursor-sim
go mod init github.com/cursor-analytics-platform/services/cursor-sim
```

### 4. Create Minimal main.go

```go
// cmd/simulator/main.go
package main

import "fmt"

const Version = "2.0.0"

func main() {
    fmt.Printf("cursor-sim v%s\n", Version)
}
```

### 5. Copy and Update Makefile

Copy from v1, update paths and targets.

### 6. Create .golangci.yml

Copy from v1 or create minimal config.

## Test Cases

```go
// No tests for this step - just build verification
// Run: go build ./...
// Run: golangci-lint run
```

## Files to Create

- services/cursor-sim/cmd/simulator/main.go
- services/cursor-sim/go.mod
- services/cursor-sim/Makefile
- services/cursor-sim/.golangci.yml

## Definition of Done

- v1 archived, v2 structure in place
- Build succeeds, linter passes
- Commit with message following template
