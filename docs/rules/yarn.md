# Yarn Rules

## YRN001: yarn npmMinimalAgeGate is below 7 days

**Severity:** high

Yarn berry 4.10+ supports `npmMinimalAgeGate` (in minutes, or as a duration string like `"7d"`), which delays installation of freshly published packages. A value of at least 7 days mitigates supply-chain attacks via fast-yanked malicious releases. Yarn classic (1.x) does not support this setting and is skipped.

**What it checks:**
- Whether `yarn` is on `PATH` (skipped if not installed)
- Whether `yarn --version` reports `2.x+` (yarn 1.x is skipped)
- Whether `yarn config get npmMinimalAgeGate` resolves to `>= 7` days (numeric values are interpreted as minutes; duration strings are parsed via the shared duration helper)

**Remediation:**
```yaml
# ~/.yarnrc.yml
npmMinimalAgeGate: 10080
```
```bash
# or
yarn config set --home npmMinimalAgeGate 10080
```

## YRN002: yarn install scripts are not disabled

**Severity:** high

Yarn berry runs lifecycle scripts (`preinstall`, `install`, `postinstall`) declared by dependencies during install. Setting `enableScripts: false` in the user/global `.yarnrc.yml` disables this behaviour, removing the primary execution vector used by compromised npm packages. Yarn classic (1.x) lacks an equivalent global toggle and is skipped.

**What it checks:**
- Whether `yarn` is on `PATH` (skipped if not installed)
- Whether `yarn --version` reports `2.x+` (yarn 1.x is skipped)
- Whether `yarn config get enableScripts` returns `false`

**Remediation:**
```bash
yarn config set --home enableScripts false
```
```yaml
# ~/.yarnrc.yml
enableScripts: false
```

Note: with `enableScripts: false`, packages that need a native build (`esbuild`, `sharp`, `node-gyp` users) must be allow-listed via `npmPreapprovedPackages` in `.yarnrc.yml`.

## YRN003: yarn auth token stored in plaintext in .yarnrc.yml

**Severity:** critical

Yarn berry stores `npmAuthToken` and `npmAuthIdent` either at the top level of `~/.yarnrc.yml` or nested under `npmRegistries:` / `npmScopes:` blocks. Literal values left on disk are recoverable by any process that can read the file — malware, compromised dependencies, sync clients, or backups. Env-var references such as `${NPM_TOKEN}` are accepted as safe indirection and pass. Yarn classic (1.x) does not use this file format and is skipped.

**What it checks:**
- Whether `yarn` is on `PATH` (skipped if not installed)
- Whether `yarn --version` reports `2.x+` (yarn 1.x is skipped)
- Reads `~/.yarnrc.yml` if it exists
- Flags any literal (non-empty, non-env-var) value assigned to `npmAuthToken` or `npmAuthIdent`, at any indentation level

**Remediation:**
```yaml
# ~/.yarnrc.yml — env-var indirection, safe
npmAuthToken: "${NPM_TOKEN}"
```

Inject `NPM_TOKEN` from a keychain, password manager, or CI secret store at session start. Rotate any token that has been on disk.
