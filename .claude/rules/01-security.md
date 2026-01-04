---
description: Security and compliance guardrails for all development
---

# Security Rules

Enterprise-grade security constraints that always apply.

---

## NEVER

- **Commit secrets**: .env files, API keys, credentials.json, private tokens
- **Expose credentials**: Don't include secrets in commit messages, logs, or output
- **Run destructive git commands**: `git push --force` to main/master, `git reset --hard` without explicit user request
- **Execute untrusted scripts**: Without sandbox mode or explicit user approval
- **Store PII**: In logs, comments, test data, or configuration files
- **Skip pre-commit hooks**: Don't use `git commit --no-verify` or `--no-gpg-sign`
- **Modify production files**: Without explicit user approval
- **Use hardcoded paths**: Always use environment variables or relative paths
- **Accept untrusted input**: Without validation at system boundaries

---

## ALWAYS

- **Validate file paths** before operations to prevent directory traversal
- **Use environment variables** for secrets (API keys, database URLs)
- **Prefer sandbox mode** for unknown scripts: `/sandbox`
- **Review hook code** before enabling in settings
- **Check file permissions** after creation (especially for config files)
- **Log security-sensitive operations** for audit trails
- **Rotate secrets** regularly and update documentation
- **Use HTTPS** for external service connections
- **Validate input types** at API boundaries
- **Fail securely**: Default to deny, explicitly allow

---

## Implementation

### Secrets Protection
- `.env` and `*.local.json` files are .gitignore'd (verified)
- Use `CLAUDE_PROJECT_DIR` and environment variables
- Document where secrets should be stored (user's home dir)

### Git Safety
- Prevent force push to main/master (hooks can enforce)
- Confirm before destructive operations
- Include signed commits where available

### Access Control
- Respect file permissions
- Don't modify outside project root
- Verify write permissions before editing

### Audit
- Log sensitive operations
- Track who made what changes
- Maintain audit trail for compliance

---

## See Also

- Hook configuration in `.claude/settings.local.json`
- Repository guardrails in `02-repo-guardrails.md`
- Coding standards in `03-coding-standards.md`
