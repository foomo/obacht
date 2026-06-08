# GitHub CLI Rules

## GH001: GitHub CLI telemetry is not disabled

**Severity:** warn

The GitHub CLI (`gh`) sends telemetry by default. The rule accepts any of three opt-out signals as sufficient to pass:

1. `telemetry: false` in the gh config file. The config dir is resolved per gh's own precedence: `$GH_CONFIG_DIR` → `$XDG_CONFIG_HOME/gh` → `$HOME/.config/gh`.
2. `GH_TELEMETRY` env var set to `0`, `false`, `off`, or `no` (case-insensitive).
3. `DO_NOT_TRACK` env var set to `1` or `true` (cross-tool W3C-style opt-out).

**What it checks:**
- Whether `gh` is installed (skipped otherwise)
- The resolved `<config-dir>/config.yml` for a `telemetry: false` line
- The `GH_TELEMETRY` and `DO_NOT_TRACK` environment variables

**Remediation:**
```bash
# Option A — config file (uses gh's resolved config dir)
printf 'telemetry: false\n' >> "${GH_CONFIG_DIR:-${XDG_CONFIG_HOME:-$HOME/.config}/gh}/config.yml"

# Option B — env var in your shell profile
export GH_TELEMETRY=0
# or
export DO_NOT_TRACK=1
```

## GH002: Git is not configured to use `gh` as credential helper

**Severity:** warn

When the GitHub CLI is installed, the recommended setup delegates git's credential helper to `!gh auth git-credential` for both `https://github.com` and `https://gist.github.com`. Without this, tokens may fall back to plaintext storage or an inconsistent credential cache.

**What it checks:**
- Whether `gh` and `git` are installed (skipped otherwise)
- Whether `credential.https://github.com.helper` includes `!gh auth git-credential`
- Whether `credential.https://gist.github.com.helper` includes `!gh auth git-credential`

**Remediation:**
```bash
gh auth setup-git
```
