# npm Rules

## NPM001: npm minimum release age is below 7 days

**Severity:** high

npm 11.10+ supports the `min-release-age` setting, which refuses to install package versions newer than the configured number of days. Setting this to at least 7 days mitigates supply-chain attacks where malicious releases are typically yanked within hours of publication.

**What it checks:**
- Whether `npm` is on `PATH` (skipped if not installed)
- Whether `npm config get min-release-age` (or the legacy `minimum-release-age`) returns a value `>= 7` days

**Remediation:**
```bash
# user-level config
npm config set min-release-age 7 --location=user
```

## NPM002: npm install scripts are not disabled

**Severity:** high

npm runs lifecycle scripts (`preinstall`, `install`, `postinstall`) declared by dependencies during `npm install`. These scripts are the primary execution vector for compromised packages — most malicious npm packages deliver their payload via a `postinstall` script. Setting `ignore-scripts=true` in the user/global npm config disables this behaviour for all installs.

**What it checks:**
- Whether `npm` is on `PATH` (skipped if not installed)
- Whether `npm config get ignore-scripts` returns `true`

**Remediation:**
```bash
npm config set ignore-scripts true --location=user
```

Note: with `ignore-scripts=true`, packages that genuinely require a build step (`esbuild`, `sharp`, `node-gyp` users) will fail to install correctly until you allow-list them per project (e.g. `npm rebuild <pkg>` after install, or use `--ignore-scripts=false` for the one install).
