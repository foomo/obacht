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
