# Architecture

## Data Flow

```
collectors (concurrent) → Facts JSON (schema v1.0)
  → embedded OPA evaluation with Rego policies
  → parse findings → diff against rule list
  → CheckResult per rule → reporter
```

## Public API (`pkg/`)

### `pkg/schema`
Type-safe data contracts:
- `Facts` — collected environment data
- `ScanResult` — evaluation results with summary
- `Rule` — rule metadata

### `pkg/engine`
OPA evaluation engine:
- Evaluates policies in-process via embedded OPA library
- Maps findings back to rules

## Internal Packages

### `internal/collector`
Twelve concurrent collectors gathering environment facts:
- Claude Code, Credentials, Docker, Environment, Git, Kubernetes, OS, PATH, Privacy, Shell, SSH, Tools
- Three-state status: ok, skipped, error

### `internal/reporter`
Output formatting:
- Pretty (lipgloss-styled terminal output)
- JSON (machine-readable)

### `internal/cli`
Cobra command definitions:
- `scan` — main security scan
- `explain` — rule details
- `doctor` — setup diagnostics
