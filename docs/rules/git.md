# Git Rules

## GIT001: Git credential helper stores passwords in plaintext

**Severity:** high

Using the `store` credential helper saves passwords as plaintext in `~/.git-credentials`. This file can be read by any process running as your user, exposing credentials for all configured Git remotes.

**What it checks:**
- Global Git configuration for `credential.helper`
- Flags the `store` helper as insecure

**Remediation:**
```bash
# Use the OS keychain instead
git config --global credential.helper osxkeychain   # macOS
git config --global credential.helper libsecret      # Linux
```

## GIT002: Git commit signing is not enabled

**Severity:** warn

Unsigned commits can be trivially spoofed by setting `user.email` to any value. Enabling commit signing provides cryptographic proof of authorship.

**What it checks:**
- Global Git configuration for `commit.gpgsign`
- Whether a signing key is configured

**Remediation:**
```bash
git config --global commit.gpgsign true
git config --global user.signingkey <your-key-id>
```
