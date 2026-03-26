# Shell Rules

## SHL001: Shell history file has weak permissions

**Severity:** warn

Shell history files (e.g., `~/.bash_history`, `~/.zsh_history`) may contain sensitive commands including passwords, tokens, or connection strings typed on the command line. Weak permissions allow other users to read your command history.

**What it checks:**
- File permissions on common shell history files
- Ensures permissions are `0600` or stricter

**Remediation:**
```bash
chmod 600 ~/.bash_history
chmod 600 ~/.zsh_history
```
