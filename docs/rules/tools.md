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
