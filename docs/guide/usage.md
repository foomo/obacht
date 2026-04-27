# Usage

## Commands

### `obacht scan`

Scan the local development environment for security issues.

```bash
obacht scan [flags]
```

**Flags:**
- `--format <pretty|json>` — Output format (default: pretty)
- `--category <categories>` — Comma-separated list of categories to scan
- `--severity <severities>` — Comma-separated list of severities to include (critical, high, warn, info)
- `--rules-dir <path>` — Load additional rules from a directory (expects `policies/` and `inputs/` subdirectories)
- `--verbose` — Enable verbose output

**Exit Codes:**
- `0` — No failed checks
- `1` — One or more failed checks
- `2` — Runtime error

### `obacht explain <rule-id>`

Show detailed information about a specific rule.

```bash
obacht explain SSH001
```

### `obacht doctor`

Check obacht dependencies and configuration.

```bash
obacht doctor
```
