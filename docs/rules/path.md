# PATH Rules

## PTH001: World-writable directory in PATH

**Severity:** high

A world-writable directory in your PATH allows any user on the system to place malicious executables that could be run instead of legitimate commands. This is a common privilege escalation and code injection vector.

**What it checks:**
- Each directory listed in the `$PATH` environment variable
- Whether any directory has world-writable permissions

**Remediation:**
```bash
# Remove world-writable permission from the directory
chmod o-w /path/to/directory

# Or remove the directory from PATH
export PATH=$(echo "$PATH" | tr ':' '\n' | grep -v '/unsafe/dir' | tr '\n' ':')
```

## PTH002: Relative path entry in PATH

**Severity:** warn

Relative paths (e.g., `.` or `bin`) in PATH mean that command resolution depends on the current working directory. An attacker who can write files to a directory you visit could execute arbitrary code.

**What it checks:**
- Each entry in the `$PATH` environment variable
- Whether any entry is a relative path (does not start with `/`)

**Remediation:**
```bash
# Convert relative paths to absolute paths or remove them
# Edit your shell profile (~/.bashrc, ~/.zshrc) and ensure all PATH entries are absolute
export PATH="/usr/local/bin:/usr/bin:/bin"
```
