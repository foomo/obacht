# pnpm Rules

## PNP001: pnpm minimumReleaseAge is below 7 days

**Severity:** high

pnpm 10.16+ supports the `minimumReleaseAge` setting (in minutes), which blocks installation of package versions newer than the configured threshold. Setting it to at least 10080 minutes (7 days) reduces exposure to compromised packages.

**What it checks:**
- Whether `pnpm` is on `PATH` (skipped if not installed)
- Whether `pnpm config get minimumReleaseAge` returns a value `>= 10080`

**Remediation:**
```yaml
# pnpm-workspace.yaml
minimumReleaseAge: 10080
```
```bash
# or user-level
pnpm config set minimumReleaseAge 10080 --location=user
```
