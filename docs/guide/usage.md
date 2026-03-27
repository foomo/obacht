# Usage

## Commands

### `bouncer scan`

Scan the local development environment for security issues.

```bash
bouncer scan [flags]
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

### `bouncer explain <rule-id>`

Show detailed information about a specific rule.

```bash
bouncer explain SSH001
```

### `bouncer doctor`

Check bouncer dependencies and configuration.

```bash
bouncer doctor
```
