# Homebrew Rules

## BRW001: Homebrew auto-update is disabled

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

## BRW002: Homebrew metadata is stale

**Severity:** warn

Homebrew caches formula and cask metadata locally. Stale metadata (older than 7 days) means newly published security patches in dependencies will not be discovered or installed when packages are added or upgraded. Refresh metadata regularly to ensure timely access to upstream fixes.

**What it checks:**
- mtime of `$(brew --cache)/api/formula.jws.json`, `cask.jws.json`, or `$(brew --repository)/.git/FETCH_HEAD`
- Skipped if `brew` is installed but no timestamp source is available (e.g. metadata never refreshed)

**Remediation:**
```bash
brew update
```

## BRW003: Homebrew analytics is not disabled

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
