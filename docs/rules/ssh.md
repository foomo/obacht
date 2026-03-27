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

## SSH003: SSH StrictHostKeyChecking is disabled

**Severity:** high

Setting `StrictHostKeyChecking` to `no` in SSH config disables host key verification, making connections vulnerable to man-in-the-middle attacks. An attacker could intercept connections and impersonate the remote server.

**What it checks:**
- Whether `StrictHostKeyChecking no` appears in `~/.ssh/config`

**Remediation:**
```bash
# Remove or change the setting in ~/.ssh/config
# Replace: StrictHostKeyChecking no
# With:    StrictHostKeyChecking ask
```

## SSH004: SSH agent forwarding is enabled globally

**Severity:** warn

Enabling `ForwardAgent` globally (under `Host *`) allows any remote server to use your local SSH agent. A compromised server could use your forwarded keys to access other systems. Only enable agent forwarding for specific trusted hosts.

**What it checks:**
- Whether `ForwardAgent yes` is set in the `Host *` section of `~/.ssh/config`

**Remediation:**
```bash
# Remove ForwardAgent yes from Host * section in ~/.ssh/config
# Add it only to specific trusted hosts:
Host trusted-server
    ForwardAgent yes
```
