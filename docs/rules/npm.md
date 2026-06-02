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

## NPM003: npm allow-git is not restricted to none or root

**Severity:** high

npm 11.10+ supports the `allow-git` setting with three values: `all`, `none`, and `root`. Git-based dependencies — direct or transitive — can ship a `.npmrc` that overrides the git executable path, enabling arbitrary code execution at install time that bypasses `ignore-scripts`. `allow-git=none` blocks all git dependencies; `allow-git=root` permits only those declared in the project's own `package.json`. The default (`all`) is expected to change to `none` in npm v12.

**What it checks:**
- Whether `npm` is on `PATH` (skipped if not installed)
- Whether `npm config get allow-git` returns `none` or `root`

**Remediation:**
```bash
# strictest — recommended
npm config set allow-git none --location=user

# or, if you need to install git deps declared in your own package.json:
npm config set allow-git root --location=user
```

## NPM004: npm registry is not https

**Severity:** high

npm's default registry should be reached over HTTPS so package tarballs, version metadata, and auth tokens are protected from network-level tampering and interception. An explicit downgrade to `http://` in `.npmrc` disables TLS for every install and every publish.

**What it checks:**
- Whether `npm` is on `PATH` (skipped if not installed)
- Whether `npm config get registry` returns a URL starting with `https://`

**Remediation:**
```bash
# default public registry
npm config set registry https://registry.npmjs.org/ --location=user

# or your private registry's HTTPS URL
npm config set registry https://npm.yourcompany.com/ --location=user
```
