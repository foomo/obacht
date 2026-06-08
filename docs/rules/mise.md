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

## MIS002: mise not_found_auto_install is enabled

**Severity:** warn

mise's `not_found_auto_install` defaults to `true`, meaning a typo'd or attacker-controlled tool name in `mise.toml` is silently installed on first shell activation. Disabling forces explicit `mise install` and surfaces supply-chain decisions to the user.

**What it checks:**
- Whether `mise` is on `PATH` (skipped if not installed)
- `mise settings get not_found_auto_install` — fails unless the value is `false`

**Remediation:**
```toml
# ~/.config/mise/config.toml
[settings]
not_found_auto_install = false
```
```bash
# or
mise settings set not_found_auto_install false
```

## MIS003: mise status.missing_tools is not "always"

**Severity:** info

mise's `status.missing_tools` controls when mise warns about tools listed in config but not installed. The default (`"if_other_versions_installed"`) can hide drift. Setting it to `"always"` surfaces missing tools in every shell, reducing the chance a check-in references something nobody else has.

**What it checks:**
- Whether `mise` is on `PATH` (skipped if not installed)
- `mise settings get status.missing_tools` — fails unless the value is `always`

**Remediation:**
```toml
# ~/.config/mise/config.toml
[settings.status]
missing_tools = "always"
```
```bash
# or
mise settings set status.missing_tools always
```

## MIS004: mise github.credential_command is not "gh auth token"

**Severity:** info

mise hits GitHub for release metadata when installing many tools. Without a credential, requests are anonymous and rate-limited per-IP. Setting `github.credential_command = "gh auth token"` reuses the local `gh` CLI session for authenticated requests, raising rate limits and avoiding unauthenticated traffic.

**What it checks:**
- Whether `mise` is on `PATH` (skipped if not installed)
- `mise settings get github.credential_command` — fails unless the value is exactly `gh auth token`

**Remediation:**
```toml
# ~/.config/mise/config.toml
[settings.github]
credential_command = "gh auth token"
```
```bash
# or
mise settings set github.credential_command "gh auth token"
```
