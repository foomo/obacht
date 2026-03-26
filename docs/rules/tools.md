# Tools Rules

## TOL001: Security-relevant tool is missing

**Severity:** info

Certain security tools are recommended for a secure development environment. Missing tools may indicate gaps in your security workflow.

**What it checks:**
- Presence of recommended security tools on the system PATH
- Tools such as `opa`, `cosign`, `trivy`, `gitleaks`, and similar utilities

**Remediation:**
```bash
# Install missing tools via Homebrew
brew install opa cosign trivy gitleaks

# Or via mise
mise install opa
```
