# Claude Code Hooks

This directory contains hooks that extend Claude Code's capabilities for spec-driven development in the Cursor Analytics Platform project.

## Overview

Hooks are automated checks and context providers that help Claude follow project conventions and maintain quality standards. They integrate into the development workflow at key points to enforce TDD practices and spec-driven development.

## Available Hooks

### pre-implementation.md

A checklist Claude should follow before implementing any feature. This ensures specifications are read, tests are written first, and dependencies are verified.

### post-test.sh

Runs after tests to analyze coverage and verify thresholds are met. Supports all three services with their respective testing tools.

### validate-spec.sh

Validates that all specification files exist and contain required sections. Helps maintain documentation quality.

## Usage Guidelines

When Claude is asked to implement a feature, it should mentally run through the pre-implementation checklist before writing any code. This ensures the TDD workflow is followed and specifications are consulted.

After running tests, the post-test hook can verify coverage meets the 80% threshold required for this project.

Before any significant documentation changes, the validate-spec hook ensures all specs remain complete and well-formed.
