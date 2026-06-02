# mise Rules

## MIS001: mise minimum_release_age is below 7 days

**Severity:** high

mise supports the `minimum_release_age` setting, which filters tool versions by release date. A relative duration like `"7d"` or an absolute cutoff date (e.g. `"2024-06-01"`) are both accepted. A value below 7 days offers limited supply-chain protection. Unrecognized formats (e.g. mise's `Nm` for months) are reported as a skip rather than a false pass.

**What it checks:**
- Whether `mise` is on `PATH` (skipped if not installed)
- `mise settings get minimum_release_age`, parsed as a duration string or absolute date (age computed against `now`)

**Remediation:**
```toml
# ~/.config/mise/config.toml
[settings]
minimum_release_age = "7d"
```
```bash
# or
mise settings set minimum_release_age 7d
```
