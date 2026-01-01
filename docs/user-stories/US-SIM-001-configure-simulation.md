# User Story: US-SIM-001

## Configure Simulation Parameters

**Story ID:** US-SIM-001  
**Feature:** [F001 - Simulator Core](../features/F001-simulator-core.md)  
**Priority:** High  
**Story Points:** 3

---

## Story

**As a** developer setting up the analytics platform,  
**I want to** configure the simulator with specific parameters at startup,  
**So that** I can control the characteristics of the generated data to match different testing scenarios.

---

## Description

The simulator needs to accept command-line arguments that control the simulation behavior. This includes specifying the number of developers to simulate, the rate of event generation, and a random seed for reproducibility. Without this capability, the simulator would produce the same data every time, limiting its usefulness for testing edge cases.

The primary use cases include simulating a small team of 10 developers for quick testing, a medium team of 50 developers for realistic development, and a large team of 200 developers for stress testing the aggregator. The velocity parameter allows testing both quiet periods (low activity) and busy sprints (high activity).

---

## Acceptance Criteria

### AC1: Help Flag Displays Usage Information
**Given** the user runs `./cursor-sim --help`  
**When** the command executes  
**Then** a help message displays showing all available flags with their descriptions and default values

### AC2: Port Configuration
**Given** the user runs `./cursor-sim --port=9090`  
**When** the server starts  
**Then** it listens on port 9090 instead of the default 8080

### AC3: Developer Count Configuration
**Given** the user runs `./cursor-sim --developers=100`  
**When** the simulator initializes  
**Then** exactly 100 developer profiles are generated

### AC4: Velocity Configuration
**Given** the user runs `./cursor-sim --velocity=low`  
**When** events are generated  
**Then** the event rate is approximately 10 events per developer per hour

### AC5: Seed Reproducibility
**Given** the user runs `./cursor-sim --seed=12345` twice  
**When** both instances generate data  
**Then** both instances produce identical developer profiles and event patterns

### AC6: Invalid Input Handling
**Given** the user runs `./cursor-sim --developers=-5`  
**When** the command is parsed  
**Then** an error message explains the invalid value and the program exits with code 1

### AC7: Startup Time
**Given** the user runs `./cursor-sim --developers=100`  
**When** measuring startup time  
**Then** the server is ready to accept requests within 5 seconds

---

## Technical Notes

The CLI should use Cobra for argument parsing, which is the standard in Go projects. Configuration should also be loadable from environment variables as a fallback, following the 12-factor app methodology. The priority order should be command-line flags, then environment variables, then default values.

Default values are port 8080, developers 50, velocity "medium", fluctuation 0.2, and a random seed. The history parameter should default to 30 days.

---

## Definition of Done

- [ ] All acceptance criteria pass
- [ ] Unit tests cover flag parsing logic
- [ ] Integration test verifies server starts on configured port
- [ ] Help text is clear and comprehensive
- [ ] Code reviewed and merged to main branch

---

## Related Tasks

- [TASK-002](../tasks/TASK-002-sim-cli.md): Implement CLI with Cobra
