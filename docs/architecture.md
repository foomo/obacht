# Architecture

## Data Flow

```
collectors (concurrent) → Facts JSON (schema v1.0)
  → write to temp dir alongside .rego files
  → single `opa eval` invocation
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
- Writes policies and facts to a secure temp directory
- Invokes external OPA binary
- Maps findings back to rules

## Internal Packages

### `internal/collector`
Nine concurrent collectors gathering environment facts:
- SSH, Git, Docker, Kubernetes, Environment, Shell, Tools, PATH, OS
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
