# APT Rules

## APT001: APT metadata is stale

**Severity:** warn

APT caches repository metadata locally. Stale metadata (older than 7 days) means newly published security patches in dependencies will not be discovered or installed when packages are added or upgraded. Refresh metadata regularly to ensure timely access to upstream fixes.

**What it checks:**
- mtime of `/var/lib/apt/periodic/update-success-stamp` or `/var/cache/apt/pkgcache.bin`
- Skipped if `apt-get` is not installed or no timestamp source is available

**Remediation:**
```bash
sudo apt update
```
