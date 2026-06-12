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

## BUN003: bun auth token stored in plaintext in bunfig.toml

**Severity:** critical

Bun stores registry credentials in `bunfig.toml` under `[install.registry]`, `[install.scopes."<scope>"]`, or as inline tables beneath `[install.scopes]`. Literal values for `token`, `username`, or `password` left on disk are recoverable by any process that can read the file — malware, compromised dependencies, sync clients, or backups. Env-var references such as `$NPM_TOKEN` are accepted as safe indirection and pass.

**What it checks:**
- Whether `bun` is on `PATH` (skipped if not installed)
- Resolves bunfig from `$BUN_CONFIG_FILE`, `$XDG_CONFIG_HOME/.bunfig.toml`, `$HOME/.config/.bunfig.toml`, then `$HOME/.bunfig.toml`
- Walks each `[install.registry]` and `[install.scopes…]` section (including inline-table values) and flags any literal `token`, `username`, or `password` assignment

**Remediation:**
```toml
# ~/.bunfig.toml — env-var indirection, safe
[install.scopes]
"@myorg" = { url = "https://npm.example.com/", token = "$NPM_TOKEN" }
```

Inject `NPM_TOKEN` from a keychain, password manager, or CI secret store at session start. Rotate any token that has been on disk.
