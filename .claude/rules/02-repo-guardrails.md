---
description: Repository safety and file protection guardrails
---

# Repository Guardrails

Protect the repository from accidental damage and maintain integrity.

---

## NEVER

- **Modify files outside project root**: All operations within project directory
- **Delete files** without explicit user request
- **Overwrite user changes** without confirmation
- **Use `git push --force`** to main or master branch (requires force push, never do silently)
- **Skip hooks**: Don't use `--no-verify`, `--no-gpg-sign` without user request
- **Bulk modify files** without user confirmation
- **Change file encoding** from UTF-8 (except where explicitly needed)
- **Stage unrelated files**: Only stage files related to current task
- **Revert commits** that aren't yours or don't have explicit user approval
- **Modify .gitignore** without understanding impact

---

## ALWAYS

- **Use absolute paths** for clarity (no relative paths like `../`)
- **Confirm before bulk operations**: Ask user before modifying 3+ files
- **Preserve file encoding**: UTF-8 for all text files
- **Stage only task-related files**: Don't include unrelated changes
- **Verify git status** before committing
- **Review diffs** before staging
- **Document file changes** in commit messages
- **Maintain .gitignore** accuracy
- **Use git add** not `git add -A` without review
- **Commit atomic changes**: One logical unit per commit

---

## Write Scope

### Protected Directories
- `.git/` - Repository metadata
- `node_modules/` - Dependencies
- `go.mod` - Go module file (use `go get` instead)
- `package-lock.json` - Use npm commands instead

### Restricted Operations
- Don't delete directories (user confirms first)
- Don't move large files (ask about refactoring)
- Don't create large new files without warning

### Safe Modifications
- Code files (.go, .ts, .tsx, .js)
- Test files (*_test.go, *.test.ts, *.spec.ts)
- Configuration files (with user confirmation)
- Documentation files (.md)
- Spec files (SPEC.md with proper updates)

---

## Staging and Committing

### Before Staging
```bash
git status          # See all changes
git diff HEAD       # Review changes
```

### Staging Rule
```bash
git add <specific-files>  # Stage only related files
```

### Commit Format
```bash
git commit -m "type(scope): description

Detailed explanation if needed.

Files changed: N
Time: Xh / Yh est"
```

---

## Conflict Resolution

- If file conflict detected: Stop and ask user for guidance
- If merge conflict: Don't auto-resolve without user confirmation
- If branch diverged: Ask user before rebasing

---

## See Also

- Security rules in `01-security.md`
- Coding standards in `03-coding-standards.md`
- SDD process in `04-sdd-process.md`
