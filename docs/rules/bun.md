# Bun Rules

## BUN001: bun install.minimumReleaseAge is below 7 days

**Severity:** high

Bun 1.3+ supports `install.minimumReleaseAge` in `bunfig.toml` (value in seconds). Setting it to at least 604800 (7 days) blocks installation of fresh, potentially compromised package versions. The global bunfig is resolved from `$BUN_CONFIG_FILE`, `$XDG_CONFIG_HOME/.bunfig.toml`, `$HOME/.config/.bunfig.toml`, and finally `$HOME/.bunfig.toml`.

**What it checks:**
- Whether `bun` is on `PATH` (skipped if not installed)
- The `[install].minimumReleaseAge` value in the resolved bunfig (parsed via the shared duration helper)

**Remediation:**
```toml
# ~/.bunfig.toml
[install]
minimumReleaseAge = 604800
```

## BUN002: bun install scripts are not disabled

**Severity:** high

Bun runs lifecycle scripts (preinstall, install, postinstall, prepare) declared by dependencies during `bun install`. These scripts are the primary execution vector for compromised packages — most malicious npm releases deliver their payload via a postinstall script. Setting `[install].ignoreScripts = true` in the global bunfig disables this behaviour for all installs. Bun has no `config get` subcommand, so the global bunfig is resolved from `$BUN_CONFIG_FILE`, `$XDG_CONFIG_HOME/.bunfig.toml`, `$HOME/.config/.bunfig.toml`, and finally `$HOME/.bunfig.toml`.

**What it checks:**
- Whether `bun` is on `PATH` (skipped if not installed)
- The `[install].ignoreScripts` value in the resolved bunfig

**Remediation:**
```toml
# ~/.bunfig.toml
[install]
ignoreScripts = true
```
