# Tools Rules

## TOL001: Security-relevant tool is missing

**Severity:** info

Certain security tools are recommended for a secure development environment. Missing tools may indicate gaps in your security workflow.

**What it checks:**
- Presence of recommended security tools on the system PATH
- Checks for `git`, `opa`, `gpg`, and `ssh-agent`

**Remediation:**
```bash
# Install missing tools via Homebrew
brew install opa gnupg

# Or via mise
mise install opa
```

## TOL002: Homebrew auto-update is disabled

**Severity:** warn

Homebrew auto-update ensures that formulae and cask definitions are refreshed before installing or upgrading packages. Disabling it via `HOMEBREW_NO_AUTO_UPDATE` means security patches in dependencies may not be applied promptly.

**What it checks:**
- Whether Homebrew is installed
- Whether the `HOMEBREW_NO_AUTO_UPDATE` environment variable is set

**Remediation:**
```bash
# Remove or unset the environment variable
unset HOMEBREW_NO_AUTO_UPDATE

# Remove from shell profile (~/.bashrc, ~/.zshrc)
```

## TOL003: Package manager metadata is stale

**Severity:** warn

Package managers cache repository metadata locally. Stale metadata (older than 7 days) means newly published security patches in dependencies will not be discovered or installed when packages are added or upgraded. Refresh metadata regularly to ensure timely access to upstream fixes.

**What it checks:**
- Homebrew (`brew`): mtime of `$(brew --cache)/api/formula.jws.json`, `cask.jws.json`, or `$(brew --repository)/.git/FETCH_HEAD`
- APT (`apt-get`): mtime of `/var/lib/apt/periodic/update-success-stamp` or `/var/cache/apt/pkgcache.bin`
- Skipped if the package manager is installed but no timestamp source is available (e.g. metadata never refreshed)

**Remediation:**
```bash
# Homebrew
brew update

# APT (Debian/Ubuntu)
sudo apt update
```

## TOL004: Homebrew analytics is not disabled

**Severity:** warn

Homebrew collects anonymized analytics by default, including installed formulae, OS version, and CPU architecture. Setting `HOMEBREW_NO_ANALYTICS=1` opts out of this telemetry collection. Skipped if `brew` is not installed.

**What it checks:**
- Whether `brew` is on `PATH`
- Whether the `HOMEBREW_NO_ANALYTICS` environment variable is set to `1`

**Remediation:**
```bash
# Add to ~/.zshrc, ~/.bashrc, or equivalent shell profile
export HOMEBREW_NO_ANALYTICS=1

# Or persist via brew itself
brew analytics off
```

## TOL005: Go telemetry is not disabled

**Severity:** warn

The Go toolchain (1.23+) collects local telemetry counters by default (mode `local`), and can optionally upload them to the Go team (mode `on`). Setting telemetry mode to `off` disables all collection and upload. Skipped if `go` is not installed.

**What it checks:**
- Whether `go` is on `PATH`
- Whether `go env GOTELEMETRY` returns `off`

**Remediation:**
```bash
go telemetry off
```
