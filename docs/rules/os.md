# OS Rules

## OS001: Operating system may have pending security updates

**Severity:** info

Running an operating system with pending security updates leaves known vulnerabilities unpatched. Timely updates are a fundamental security practice.

**What it checks:**
- Whether the OS has pending software updates (macOS) or security patches (Linux)
- Uses platform-specific mechanisms to detect available updates

**Remediation:**
```bash
# macOS
softwareupdate -ia

# Debian/Ubuntu
sudo apt update && sudo apt upgrade -y

# Fedora/RHEL
sudo dnf upgrade --security -y
```
