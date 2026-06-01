# Bumblebee Rules

The `bumblebee` category surfaces known-compromised packages, extensions, and MCP servers found on disk. It wraps [bumblebee](https://github.com/perplexityai/bumblebee), Perplexity's read-only inventory collector, and matches its NDJSON output against an exposure catalog bundled with obacht (vendored from bumblebee's upstream `threat_intel/` directory).

All rules in this category are skipped automatically when:

- The `bumblebee` binary is not on `PATH` — install with `go install github.com/perplexityai/bumblebee/cmd/bumblebee@latest` or via `mise`.
- The embedded exposure catalog cannot be materialised to a temp directory.

When bumblebee runs, the `baseline` profile is used: global package roots, language toolchains, editor and browser extensions, MCP host configs. Project/deep profiles are not invoked from obacht.

## BUM001: Compromised npm package present on disk

**Severity:** critical

An npm package on disk matches an entry in the bundled exposure catalog of known-compromised releases (e.g. shai-hulud worm, typosquats, credential stealers). The match is exact on `(ecosystem, name, version)`. Presence on disk is not proof of execution, but indicates exposure that responders should triage.

**Remediation:** Remove or downgrade the affected package; rotate any secrets exposed by the package's lifecycle scripts; see the catalog entry's reporting for scope.

## BUM002: Compromised PyPI package present on disk

**Severity:** critical

A PyPI package on disk matches an entry in the bundled exposure catalog of known-malicious releases. Presence indicates exposure even if the package was not imported.

**Remediation:** `pip uninstall <name>`; rotate any secrets the package could have read; review post-install hooks in `*.dist-info/INSTALLER`.

## BUM003: Compromised Go module present on disk

**Severity:** critical

A Go module recorded in `go.sum` or `go.mod` matches an entry in the bundled exposure catalog. Module presence does not imply execution (modules run only if compiled in), but warrants investigation of dependent binaries.

**Remediation:** Pin to a clean version, run `go mod tidy`, rebuild and redeploy affected binaries.

## BUM004: Compromised RubyGem present on disk

**Severity:** critical

An installed RubyGem matches an entry in the bundled exposure catalog. Gems can run install-time hooks; treat presence as potential compromise pending review.

**Remediation:** `gem uninstall <name> -v <version>`; rotate developer credentials; audit `~/.gem` for unexpected artifacts.

## BUM005: Compromised Composer package present on disk

**Severity:** critical

A Composer (Packagist) package recorded in `composer.lock` or `vendor/composer/installed.json` matches an entry in the bundled exposure catalog.

**Remediation:** Remove the affected package version, regenerate `composer.lock`, redeploy affected applications.

## BUM006: Compromised MCP server configured

**Severity:** critical

An MCP server referenced in a host config (Claude Desktop, Cline, Gemini CLI, etc.) matches an entry in the bundled exposure catalog. MCP servers run with the host process' privileges and have broad tool access.

**Remediation:** Remove the server entry from the relevant `mcp.json` / `claude_desktop_config.json` / equivalent; restart the host application; rotate credentials the server could read.

## BUM007: Compromised editor extension installed

**Severity:** critical

A VS Code / Cursor / Windsurf / VSCodium extension on disk matches an entry in the bundled exposure catalog (e.g. the nx-console campaign). Editor extensions run with the developer's local privileges.

**Remediation:** Uninstall the extension from the editor; remove any cached extension directory; rotate developer credentials; review the editor's extension audit log.

## BUM008: Compromised browser extension installed

**Severity:** critical

A Chromium-family or Firefox extension manifest on disk matches an entry in the bundled exposure catalog. Browser extensions typically have access to all visited pages and stored credentials.

**Remediation:** Uninstall the extension from each browser profile; force a password rotation for high-value accounts the affected browser was signed into.

## Updating the bundled catalog

Catalogs live in `rules/catalogs/bumblebee/`. To refresh:

1. Pick a new upstream tag at <https://github.com/perplexityai/bumblebee/releases>.
2. Update `rules/catalogs/bumblebee/VENDOR.md` with the new tag.
3. Replace each `*.json` from the matching `threat_intel/` directory of that tag.
4. Run `make check`.
