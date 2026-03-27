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

## GIT003: Git safe.directory set to wildcard

**Severity:** high

Setting `safe.directory` to `*` disables Git's ownership checks for all repositories. This defeats the protection against CVE-2022-24765 where a malicious repository in a shared directory could execute arbitrary commands via Git hooks.

**What it checks:**
- Whether `safe.directory` is set to `*` in global Git configuration

**Remediation:**
```bash
git config --global --unset-all safe.directory
```

## GIT004: Global gitignore does not exclude .env files

**Severity:** warn

Without a global gitignore rule for `.env` files, secrets stored in `.env` files can accidentally be committed to repositories. A global exclusion acts as a safety net alongside per-repo `.gitignore` files.

**What it checks:**
- Whether a global `core.excludesfile` is configured
- Whether that file contains a `.env` exclusion pattern

**Remediation:**
```bash
echo '.env' >> ~/.gitignore_global
git config --global core.excludesfile ~/.gitignore_global
```
