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

## PNP002: pnpm blockExoticSubdeps is not enabled

**Severity:** high

pnpm 10.26+ supports `blockExoticSubdeps`. When `true`, transitive dependencies cannot pull code from git repositories or raw tarball URLs — sources that bypass registry security scanning. Direct dependencies declared in `package.json` may still use exotic specifiers; only their further subdeps are blocked. Defaults to `true` in pnpm v11; enabling it early on 10.26+ mitigates a known supply-chain vector.

**What it checks:**
- Whether `pnpm` is on `PATH` (skipped if not installed)
- Whether `pnpm config get blockExoticSubdeps` returns `true`

**Remediation:**
```yaml
# pnpm-workspace.yaml
blockExoticSubdeps: true
```
```bash
# or user-level
pnpm config set blockExoticSubdeps true --location=user
```

## PNP003: pnpm trustPolicy is not set to no-downgrade

**Severity:** high

pnpm 10.21+ supports `trustPolicy`. With `no-downgrade`, pnpm refuses to install a package version whose trust evidence (signed publisher / provenance) has regressed compared to any earlier release — defeating compromised-publisher attacks where the attacker drops back to unsigned releases after taking over maintainer credentials. Default is unset (no trust check performed).

**What it checks:**
- Whether `pnpm` is on `PATH` (skipped if not installed)
- Whether `pnpm config get trustPolicy` returns `no-downgrade`

**Remediation:**
```yaml
# pnpm-workspace.yaml
trustPolicy: no-downgrade
```
```bash
# or user-level
pnpm config set trustPolicy no-downgrade --location=user
```

Note: known friction with dedupe and `@types/*` wildcards on some pnpm 10.2x releases — see `trustPolicyExclude` for per-package opt-outs.

## PNP004: pnpm strictDepBuilds is not enabled

**Severity:** high

pnpm 10.3+ supports `strictDepBuilds`. When `true`, pnpm exits with a non-zero code if any dependency tries to run a lifecycle script that has not been explicitly allowed via `allowBuilds` (or the legacy `onlyBuiltDependencies`). Without it, unallowed builds become silent warnings — a transitive package can quietly slip past review. Enabling this turns warnings into CI-blocking failures.

**What it checks:**
- Whether `pnpm` is on `PATH` (skipped if not installed)
- Whether `pnpm config get strictDepBuilds` returns `true`

**Remediation:**
```yaml
# pnpm-workspace.yaml
strictDepBuilds: true

# pair with an explicit allow-list of packages that genuinely need build scripts:
allowBuilds:
  esbuild: true
  rolldown: true
```
```bash
# or user-level
pnpm config set strictDepBuilds true --location=user
```

## PNP005: pnpm auth token stored in plaintext in user .npmrc

**Severity:** critical

pnpm reuses npm's `.npmrc` for registry authentication. Literal values for `_authToken`, `_auth`, `_password`, and `_authIdent` left at rest can be exfiltrated by any local process, backup, or compromised dependency, granting publish rights to the developer's packages. Env-var references such as `${NPM_TOKEN}` are accepted as safe indirection and pass. This rule overlaps NPM005 by design — pnpm and npm read the same file.

**What it checks:**
- Whether `pnpm` is on `PATH` (skipped if not installed)
- Reads `${NPM_CONFIG_USERCONFIG:-$HOME/.npmrc}` if it exists
- Flags any literal (non-empty, non-env-var) value assigned to `_authToken`, `_auth`, `_password`, or `_authIdent`

**Remediation:**
```ini
# ~/.npmrc — env-var indirection, safe
//registry.npmjs.org/:_authToken=${NPM_TOKEN}
```

Inject `NPM_TOKEN` from a keychain, password manager, or CI secret store at session start. Rotate any token that has been on disk.
