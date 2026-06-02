# Go Rules

## GO001: Go telemetry is not disabled

**Severity:** warn

The Go toolchain (1.23+) collects local telemetry counters by default (mode `local`), and can optionally upload them to the Go team (mode `on`). Setting telemetry mode to `off` disables all collection and upload. Skipped if `go` is not installed.

**What it checks:**
- Whether `go` is on `PATH`
- Whether `go env GOTELEMETRY` returns `off`

**Remediation:**
```bash
go telemetry off
```
