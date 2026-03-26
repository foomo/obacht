# SSH Rules

## SSH001: SSH private key has weak permissions

**Severity:** high

SSH private keys should be readable only by the owner. Overly permissive file permissions allow other users on the system to read your private keys, potentially compromising authentication to remote servers.

**What it checks:**
- File permissions on all private key files in `~/.ssh/`
- Ensures permissions are `0600` or stricter

**Remediation:**
```bash
chmod 600 ~/.ssh/id_*
```

## SSH002: SSH directory has weak permissions

**Severity:** high

The `~/.ssh` directory should only be accessible by the owner. Weak directory permissions can expose SSH configuration, known hosts, and authorized keys to other users.

**What it checks:**
- Directory permissions on `~/.ssh/`
- Ensures permissions are `0700` or stricter

**Remediation:**
```bash
chmod 700 ~/.ssh
```
