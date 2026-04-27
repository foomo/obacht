# obacht CLI — Phase 1 (MVP) Implementation Plan

## Context

Greenfield Go CLI that inspects developer environments for security misconfigurations using OPA/Rego policies. No code exists yet. Implementing Phase 1 (MVP) **incrementally, step by step**, with security checks and sub-agents at each milestone.

## Key Decisions

| Decision | Choice |
|---|---|
| CLI name | `obacht` |
| Go module | `github.com/foomo/obacht`, Go 1.22+ |
| CLI framework | `github.com/spf13/cobra` |
| Policy engine | External OPA binary (>= v1.0.0), invoked via `opa eval` |
| Policy embedding | `embed.FS` for built-in policies, `--rules-dir` for external |
| Config loading | `github.com/knadh/koanf/v2` for YAML merging |
| Terminal styling | `github.com/charmbracelet/lipgloss` |
| TUI components | `github.com/charmbracelet/bubbletea` + `bubbles` (spinner, table) |
| Output style | Grouped by category, severity-sorted within group, all checks shown (pass/fail/skip/error) |
| TTY behavior | TTY + pretty → Bubble Tea (spinner + table); non-TTY or JSON → plain output |
| Tooling | `mise` for Go, OPA, golangci-lint, bun |
| Linting | `golangci-lint` with gosec |
| Docs | VitePress + Bun, GH Pages on `v*.*.*` tags (foomo/go pattern) |
| Testing | `testify` for assertions, `foomo/go` testing tags for integration tests |
| Schema versioning | `"1.0"` field in Facts struct from day one |

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/charmbracelet/lipgloss` | Terminal styling |
| `github.com/charmbracelet/bubbletea` | TUI runtime (spinner during scan) |
| `github.com/charmbracelet/bubbles` | Spinner, table components |
| `github.com/knadh/koanf/v2` | Config loading + multi-file YAML merging |
| `golang.org/x/sync` | `errgroup` for concurrent collectors |
| `github.com/foomo/go` | Test tags for integration tests (test dep) |
| `github.com/stretchr/testify` | Assertions & test suites (test dep) |

## Directory Layout

```
cmd/obacht/obacht.go
pkg/
  schema/       facts.go, findings.go        (public API)
  engine/       engine.go                     (public API)
internal/
  cli/          root.go, scan.go, explain.go, doctor.go, exitcodes.go
  collector/    collector.go, ssh.go, git.go, docker.go, kube.go, env.go, shell.go, tools.go, path.go, os.go
  reporter/     reporter.go, pretty.go, json.go
  preflight/    opa.go                        (OPA binary check)
policies/
  embed.go
  rules/        ssh.yaml, git.yaml, docker.yaml, kube.yaml, env.yaml, shell.yaml, tools.yaml, path.yaml, os.yaml
  rego/         ssh.rego, git.rego, docker.rego, kube.rego, env.rego, shell.rego, tools.rego, path.rego, os.rego
docs/           VitePress site
.mise.toml
.golangci.yml
.github/workflows/  release.yml (GoReleaser + docs), ci.yml (lint + test)
Makefile
.goreleaser.yaml
```

## Architecture

### Data Flow
```
collectors (concurrent, errgroup) → Facts JSON (schema v1.0)
  → write to temp dir (0700) alongside .rego files
  → single `opa eval -d <dir> -i <facts.json> 'data.obacht[_].findings[_]'`
  → parse OPA JSON output → diff against full rule ID list from rules.yaml
  → CheckResult per rule: pass (not in findings), fail (in findings), skip/error (from collector status)
  → reporter (lipgloss/bubbletea for TTY pretty, plain JSON otherwise)
```

### Collector Three-State Status
Each collector returns a status alongside its facts:
- `ok` → evaluate rules, infer pass/fail from Rego output
- `skipped` → mark all rules in that category as skip (e.g., Docker not installed)
- `error` → mark all rules in that category as error (e.g., collector timed out)

### Rule Metadata (YAML)
One file per category (e.g., `ssh.yaml`):
```yaml
rules:
  - id: SSH001
    title: SSH private key has weak permissions
    severity: high
    category: ssh
    description: |
      SSH private keys should be readable only by the owner.
      Keys with permissions wider than 0600 can be read by
      other users on the system.
    remediation: "Run: chmod 600 ~/.ssh/id_rsa"
```

### External Rules (`--rules-dir`)
- Loads all `*.rego` + `*.yaml` from the specified directory
- Merged with built-in rules via koanf
- External rules with the same ID **override** built-ins
- Multiple YAML files supported per directory

### OPA Invocation
- Single `opa eval` call with all policies (built-in + external) in one temp dir
- Temp dir is `0700`, facts JSON is `0600`, cleaned up via `defer os.RemoveAll`
- Preflight check: verify `opa` in PATH and version >= 1.0.0, suggest `brew install opa` if missing

### Env Var Detection
- Curated **allowlist** of known-dangerous env var names (e.g., `AWS_SECRET_ACCESS_KEY`, `GITHUB_TOKEN`)
- Tight suffix patterns (e.g., `*_PASSWORD`, `*_SECRET_KEY` — not `*_TOKEN` broadly)
- List defined in `env.yaml`, extensible via `--rules-dir`
- **Never** store env var values — only names + matched pattern

## Output Format

### Pretty (TTY)
Bubble Tea spinner during scan, then lipgloss-styled results:
```
SSH
  ✓ SSH001: SSH private key permissions
  ✗ SSH002: SSH directory permissions
            Evidence: ~/.ssh has mode 0755
            Fix: Run: chmod 700 ~/.ssh

Git
  ✓ GIT001: Credential helper safety
  ✓ GIT002: Commit signing enabled

Docker
  - DOC001: Docker socket permissions (skipped — Docker not installed)
  - DOC002: Docker group membership (skipped — Docker not installed)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Summary: 1 failed, 2 passed, 2 skipped (0 critical, 0 high, 1 warn, 0 info)
```

### JSON (non-TTY or `--format json`)
Full `ScanResult` with all checks, each having a `status` field (pass/fail/skip/error). No Bubble Tea runtime, plain stdout.

## Implementation Steps (incremental, one at a time)

Each step is a standalone milestone. Run security checks (golangci-lint with gosec, `go vet`, tests) after each step. Use sub-agents to parallelize independent work within each step.

---

### Step 1: Project Scaffolding & Tooling
**Goal**: Buildable binary with `--help`, all tooling configured.

- `.mise.toml` — Go 1.22+, golangci-lint, opa (>= 1.0.0), bun
- `go.mod` with module `github.com/foomo/obacht`
- `go get` all dependencies
- `cmd/obacht/obacht.go` → calls `cli.Execute()`
- `internal/cli/root.go` — cobra root with `--format`, `--verbose`, `--rules-dir` flags
- `internal/cli/exitcodes.go` — OK=0, Findings=1, Error=2
- `internal/preflight/opa.go` — check OPA binary in PATH, version check, install hint
- `.golangci.yml` — enable gosec, govet, staticcheck, errcheck, revive
- `Makefile` — build, test, lint, run, test-rego targets
- **Sub-agents**: scaffolding code + mise/lint config in parallel
- **Verify**: `mise install && make build && ./bin/obacht --help && make lint`

---

### Step 2: Facts Schema
**Goal**: Type-safe contract between collectors and policies.

- `pkg/schema/facts.go` — `Facts` struct with `SchemaVersion` field ("1.0") and sub-structs: OS, SSH, Git, Docker, Kube, Env, Shell, Tools, Path
- `pkg/schema/findings.go` — `CheckResult` (rule_id, title, severity, category, status, evidence, remediation), `ScanResult`, `Summary`, `Severity` enum, `Status` enum (pass/fail/skip/error)
- Key: `EnvFacts.SuspiciousVars` stores only var name + matched pattern, **never** the value
- `pkg/schema/facts_test.go` — round-trip marshal/unmarshal
- **Verify**: `make test && make lint`

---

### Step 3: Collector Interface + SSH & Git Collectors
**Goal**: First real data collection, testable independently.

- `internal/collector/collector.go` — `Collector` interface (`Name()`, `Collect(ctx)` returning facts + status), `CollectAll()` using `errgroup.Group` with per-collector context timeouts
- `internal/collector/ssh.go` — stat `~/.ssh/*` files, record modes. Parser tested separately.
- `internal/collector/git.go` — `git config --global --list`, parse credential helper + signing. Parser tested separately.
- Collector three-state: ok/skipped/error
- **Sub-agents**: SSH collector + Git collector in parallel
- **Verify**: `make test && make lint`

---

### Step 4: OPA Engine
**Goal**: Evaluate Rego policies against facts via external OPA, extract findings.

- `pkg/engine/engine.go` — write policies + facts to temp dir (0700), invoke `opa eval`, parse JSON result, diff findings against full rule list to infer pass/skip/error
- Rego convention: every package exposes `findings[f] { ... }` returning Finding-shaped objects
- Preflight check integrated: fail fast if OPA not installed
- `pkg/engine/engine_test.go` — hardcoded facts + test policy → expected findings
- **Verify**: `make test && make lint`

---

### Step 5: Rule Metadata + First Rego Policies (SSH + Git) — End-to-End Slice
**Goal**: First working scan with real output.

- `policies/embed.go` — `//go:embed rego/*.rego rules/*.yaml`
- `policies/rules/ssh.yaml` — SSH001, SSH002 metadata
- `policies/rules/git.yaml` — GIT001, GIT002 metadata
- `policies/rego/ssh.rego` — SSH001 (private key perms > 0600), SSH002 (.ssh dir perms > 0700)
- `policies/rego/git.rego` — GIT001 (credential helper `store`), GIT002 (signing not enabled)
- Rule loading via koanf: glob all `*.yaml` in policies/rules/ and optional `--rules-dir`
- Override semantics: external rule with same ID replaces built-in
- Rego tests with fixture JSON
- CI test: every rule ID in YAML must exist in Rego, and vice versa
- Wire up `scan` command minimally (collect → evaluate → print JSON for now)
- **Verify**: `make test && make test-rego && make lint && ./bin/obacht scan --format json`

---

### Step 6: Reporters (Pretty + JSON)
**Goal**: Polished output for both humans and machines.

- `internal/reporter/reporter.go` — `Reporter` interface, `ForFormat()` factory
- `internal/reporter/json.go` — `json.MarshalIndent` of `ScanResult` (all checks with status field)
- `internal/reporter/pretty.go` — Bubble Tea program for TTY: spinner during scan, then lipgloss-styled results grouped by category with ✓/✗/- marks. Plain lipgloss output for non-TTY pretty.
- TTY detection: TTY + pretty → Bubble Tea; non-TTY or JSON → plain output
- Bubbles table for doctor command output
- **Sub-agents**: pretty reporter + JSON reporter in parallel
- **Verify**: `make test && make lint && ./bin/obacht scan` (pretty) and `./bin/obacht scan --format json`

---

### Step 7: Remaining Collectors + Rego Rules
**Goal**: Full set of 13 rules across 9 categories.

**Collectors** (sub-agents for parallel work):
- `docker.go` — stat socket, check group membership
- `kube.go` — stat + parse kubeconfig YAML for contexts
- `env.go` — scan `os.Environ()` against allowlist from env.yaml (redact values)
- `shell.go` — history file perms, HISTCONTROL settings
- `tools.go` — `exec.LookPath` + version check
- `path.go` — split PATH, stat dirs for write permission
- `os.go` — runtime info, platform-specific patch check

**Rule metadata** (YAML per category):
- `docker.yaml`, `kube.yaml`, `env.yaml`, `shell.yaml`, `tools.yaml`, `path.yaml`, `os.yaml`

**Rego policies**:
| File | Rules | Checks |
|---|---|---|
| `docker.rego` | DOC001, DOC002 | Socket world-readable; user in docker group |
| `kube.rego` | KUB001, KUB002 | Config perms > 0600; production context |
| `env.rego` | ENV001 | Secret patterns in env vars (allowlist) |
| `shell.rego` | SHL001 | History file perms |
| `tools.rego` | TOL001 | Missing/outdated tools |
| `path.rego` | PTH001, PTH002 | Writable PATH dir; relative PATH entry |
| `os.rego` | OS001 | Stale OS updates |

- **Verify**: `make test && make test-rego && make lint && ./bin/obacht scan`

---

### Step 8: Explain & Doctor Commands
**Goal**: Complete CLI command set.

- `internal/cli/explain.go` — takes rule ID, looks up metadata from rules YAML, prints title, severity, description, remediation
- `internal/cli/doctor.go` — full diagnostic: OPA binary + version, policy validation (YAML/Rego parse + ID sync), collector health (run each, report ok/skipped/error), system info (OS, arch, shell, obacht version). Output via bubbles table.
- `scan` command gets `--category` flag to filter
- **Verify**: `./bin/obacht explain SSH001 && ./bin/obacht doctor && make lint`

---

### Step 9: VitePress Documentation + GitHub Pages
**Goal**: Documentation site deployed on git tags.

Structure:
```
docs/
  .vitepress/
    config.mts
    theme/
      custom.css
      index.ts
  public/
    logo.png
  package.json (bun)
  index.md                 — landing page (what is obacht)
  guide/
    getting-started.md     — install (mise, brew), first scan
    usage.md               — commands, flags, output formats
    custom-rules.md        — --rules-dir, writing rego + yaml, override behavior
  rules/
    index.md               — overview of all 13 built-in rules
    ssh.md                 — SSH001, SSH002
    git.md                 — GIT001, GIT002
    docker.md              — DOC001, DOC002
    kube.md                — KUB001, KUB002
    env.md                 — ENV001
    shell.md               — SHL001
    tools.md               — TOL001
    path.md                — PTH001, PTH002
    os.md                  — OS001
  architecture.md          — collector → engine → reporter, pkg/ public API
```

GitHub Actions (mirroring foomo/go pattern):
- `.github/workflows/release.yml` — triggered on `v*.*.*` tags: GoReleaser job, then docs job (bun install, vitepress build, deploy-pages)
- `.github/workflows/ci.yml` — triggered on PRs: `make lint && make test && make test-rego`

- **Sub-agents**: VitePress setup + GH Actions in parallel
- **Verify**: `cd docs && bun install && bun run build`

---

### Step 10: Build Tooling & Release
**Goal**: Cross-platform release pipeline.

- `.goreleaser.yaml` — darwin/amd64+arm64, linux/amd64+arm64, windows/amd64
- Update Makefile with `release` target
- **Verify**: `goreleaser check && make build`

## Security Checks (applied at every step)

- `golangci-lint run` with gosec enabled
- `go vet ./...`
- Never store secret values in facts — only names + patterns
- No network calls unless explicitly enabled
- Rego policies bundled in binary, no remote fetch
- Temp dir `0700`, facts JSON `0600`, cleaned up via `defer os.RemoveAll`
- CI test enforcing YAML/Rego rule ID parity

## Testing Strategy

| Layer | Approach | Tool |
|---|---|---|
| Schema | Marshal/unmarshal round-trip | `go test` + testify |
| Collector parsers | Unit tests with known input → expected output | `go test` + testify |
| Collector integration | Run against real system state | `go test` + foomo/go test tags |
| Rego policies | OPA test framework with fixture JSON | `opa test` |
| Rule ID sync | Every YAML rule ID exists in Rego and vice versa | `go test` |
| Engine | Known facts in → expected CheckResults out | `go test` + testify |
| Reporters | Golden file comparison | `go test` + testify |
| CLI (e2e) | Build binary, run commands, validate output | `go test` + testify |

## CLI Commands

- `obacht scan` — collect facts, evaluate policies, report all checks
- `obacht scan --format json` — machine-readable output
- `obacht scan --category ssh,git` — filter by category
- `obacht scan --rules-dir ./my-rules` — load external rules (override on ID collision)
- `obacht explain SSH001` — detailed rule explanation from YAML metadata
- `obacht doctor` — full diagnostic (OPA, policies, collectors, system info)

## Exit Codes

- `0` — no failed checks
- `1` — one or more failed checks
- `2` — runtime error (OPA not found, collector crash, etc.)

## Deferred to Phase 2+
- Policy packs / `--profile` / `--pack`
- SARIF/CSV output
- Per-user config for exclusions (koanf ready for this)
- Remote policy fetching
- Auto-fix / enforcement
- Signed policy bundles
- Windows-specific collectors
