# Claude Code Rules

All rules in this category are skipped automatically when the `claude` CLI is not installed. Configuration is read from `$XDG_CONFIG_HOME/claude/.claude.json` and `$XDG_CONFIG_HOME/claude/settings.json` when present, otherwise from `~/.claude.json` and `~/.claude/settings.json`.

## CLD001: Global gitignore does not exclude Claude Code local settings

**Severity:** warn

Claude Code stores user-specific settings in `.claude/settings.local.json` which may contain personal configuration. A global gitignore entry for `**/.claude/settings.local.json` prevents accidental commits.

**Remediation:**
```bash
echo '**/.claude/settings.local.json' >> ~/.gitignore_global
git config --global core.excludesfile ~/.gitignore_global
```

## CLD002: autoCompactEnabled is not disabled

**Severity:** warn

Auto-compact silently rewrites conversation context when token budgets are reached, dropping details. Disabling gives the user explicit control.

**Remediation:** Set `"autoCompactEnabled": false` in `~/.claude.json`.

## CLD003: prStatusFooterEnabled is not disabled

**Severity:** info

The PR status footer makes outbound calls to GitHub and leaks working-directory metadata to remote services.

**Remediation:** Set `"prStatusFooterEnabled": false` in `~/.claude.json`.

## CLD004: claudeInChromeDefaultEnabled is not disabled

**Severity:** warn

Claude-in-Chrome injects Claude into the default browser. Disabling keeps Claude Code scoped to the terminal.

**Remediation:** Set `"claudeInChromeDefaultEnabled": false` in `~/.claude.json`.

## CLD005: sandbox.failIfUnavailable is not enabled

**Severity:** high

When false (or unset), Claude Code silently runs without the platform sandbox if it cannot be initialized. Setting it to `true` forces fail-closed behavior.

**Remediation:** Set `"sandbox": {"failIfUnavailable": true}` in `~/.claude.json`.

## CLD006: DISABLE_COMPACT is not set in settings.json

**Severity:** warn

`DISABLE_COMPACT="1"` turns off the manual `/compact` command which rewrites conversation history.

**Remediation:** Set `"env": {"DISABLE_COMPACT": "1"}` in `~/.claude/settings.json`.

## CLD007: DISABLE_TELEMETRY is not set in settings.json

**Severity:** high

`DISABLE_TELEMETRY="1"` opts out of anonymized usage telemetry sent to Anthropic.

**Remediation:** Set `"env": {"DISABLE_TELEMETRY": "1"}` in `~/.claude/settings.json`.

## CLD008: DISABLE_BUG_COMMAND is not set in settings.json

**Severity:** warn

`DISABLE_BUG_COMMAND="1"` hides the `/bug` command, preventing accidental upload of transcript context.

**Remediation:** Set `"env": {"DISABLE_BUG_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD009: DISABLE_AUTO_COMPACT is not set in settings.json

**Severity:** warn

`DISABLE_AUTO_COMPACT="1"` disables automatic conversation compaction.

**Remediation:** Set `"env": {"DISABLE_AUTO_COMPACT": "1"}` in `~/.claude/settings.json`.

## CLD010: DISABLE_LOGIN_COMMAND is not set in settings.json

**Severity:** info

`DISABLE_LOGIN_COMMAND="1"` hides `/login`, preventing accidental account switching from an active session.

**Remediation:** Set `"env": {"DISABLE_LOGIN_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD011: DISABLE_LOGOUT_COMMAND is not set in settings.json

**Severity:** info

`DISABLE_LOGOUT_COMMAND="1"` hides `/logout`; pairs with DISABLE_LOGIN_COMMAND.

**Remediation:** Set `"env": {"DISABLE_LOGOUT_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD012: DISABLE_ERROR_REPORTING is not set in settings.json

**Severity:** high

`DISABLE_ERROR_REPORTING="1"` stops Claude Code from uploading error stack traces (file paths, command args) to Anthropic.

**Remediation:** Set `"env": {"DISABLE_ERROR_REPORTING": "1"}` in `~/.claude/settings.json`.

## CLD013: DISABLE_UPGRADE_COMMAND is not set in settings.json

**Severity:** warn

`DISABLE_UPGRADE_COMMAND="1"` hides the in-app upgrade command so Claude Code is upgraded only via the system package manager.

**Remediation:** Set `"env": {"DISABLE_UPGRADE_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD014: DISABLE_FEEDBACK_COMMAND is not set in settings.json

**Severity:** warn

`DISABLE_FEEDBACK_COMMAND="1"` hides `/feedback`, which submits free-form messages (and may include transcript context) to Anthropic.

**Remediation:** Set `"env": {"DISABLE_FEEDBACK_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD015: DISABLE_EXTRA_USAGE_COMMAND is not set in settings.json

**Severity:** warn

`DISABLE_EXTRA_USAGE_COMMAND="1"` hides the extra usage reporting command.

**Remediation:** Set `"env": {"DISABLE_EXTRA_USAGE_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD016: CLAUDE_CODE_DISABLE_FAST_MODE is not set in settings.json

**Severity:** warn

`CLAUDE_CODE_DISABLE_FAST_MODE="1"` turns off the `/fast` toggle so the runtime model cannot be flipped mid-session.

**Remediation:** Set `"env": {"CLAUDE_CODE_DISABLE_FAST_MODE": "1"}` in `~/.claude/settings.json`.

## CLD017: DISABLE_INSTALL_GITHUB_APP_COMMAND is not set in settings.json

**Severity:** warn

`DISABLE_INSTALL_GITHUB_APP_COMMAND="1"` prevents Claude Code from installing GitHub Apps (persistent org-level permissions) from inside a session.

**Remediation:** Set `"env": {"DISABLE_INSTALL_GITHUB_APP_COMMAND": "1"}` in `~/.claude/settings.json`.

## CLD018: CLAUDE_CODE_DISABLE_CRON is not set in settings.json

**Severity:** info

`CLAUDE_CODE_DISABLE_CRON="1"` turns off scheduled/cron-style autonomous task execution.

**Remediation:** Set `"env": {"CLAUDE_CODE_DISABLE_CRON": "1"}` in `~/.claude/settings.json`.

## CLD019: CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY is not set in settings.json

**Severity:** warn

`CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY="1"` suppresses periodic in-session feedback prompts.

**Remediation:** Set `"env": {"CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY": "1"}` in `~/.claude/settings.json`.

## CLD020: CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING is not set in settings.json

**Severity:** info

`CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING="1"` disables file snapshots, reducing disk/IO and stale-snapshot risk.

**Remediation:** Set `"env": {"CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING": "1"}` in `~/.claude/settings.json`.

## CLD021: CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS is not set in settings.json

**Severity:** info

`CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS="1"` opts out of experimental beta features.

**Remediation:** Set `"env": {"CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS": "1"}` in `~/.claude/settings.json`.

## CLD022: FORCE_AUTOUPDATE_PLUGINS is not set in settings.json

**Severity:** warn

`FORCE_AUTOUPDATE_PLUGINS="1"` forces plugins to auto-update so installed plugins do not drift behind upstream fixes.

**Remediation:** Set `"env": {"FORCE_AUTOUPDATE_PLUGINS": "1"}` in `~/.claude/settings.json`.

## CLD023: IS_DEMO is not set in settings.json

**Severity:** warn

`IS_DEMO="1"` enables demo-mode affordances. Reflects a deliberate user-policy choice — remove the rule if your environment should not run in demo mode.

**Remediation:** Set `"env": {"IS_DEMO": "1"}` in `~/.claude/settings.json`.

## CLD024: disableAutoMode is not set to "disable"

**Severity:** warn

`disableAutoMode="disable"` prevents Claude Code from automatically switching into autonomous modes.

**Remediation:** Set `"disableAutoMode": "disable"` in `~/.claude/settings.json`.

## CLD025: disableDeepLinkRegistration is not set to "disable"

**Severity:** high

`disableDeepLinkRegistration="disable"` stops Claude Code from registering itself as the handler for `claude://` deep links.

**Remediation:** Set `"disableDeepLinkRegistration": "disable"` in `~/.claude/settings.json`.

## CLD026: attribution.commit and attribution.pr are not empty

**Severity:** info

Claude Code appends an attribution trailer to commit messages and PR bodies. Empty strings remove these trailers, keeping authored artefacts free of tool-specific markers.

**Remediation:** Set `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json`.

## CLD027: respectGitignore is not enabled

**Severity:** warn

`respectGitignore=true` makes Claude Code honour `.gitignore` when reading files, so ignored files (credentials, vendored data) are not pulled into the conversation.

**Remediation:** Set `"respectGitignore": true` in `~/.claude/settings.json`.

## CLD028: skipWebFetchPreflight is not enabled

**Severity:** warn

`skipWebFetchPreflight=true` skips the OPTIONS preflight before WebFetch, removing an extra outbound request per fetch.

**Remediation:** Set `"skipWebFetchPreflight": true` in `~/.claude/settings.json`.

## CLD029: autoMemoryDirectory is not set to ".claude/memory"

**Severity:** info

Pins Claude Code's automatic memory store to a known location alongside the rest of `.claude`.

**Remediation:** Set `"autoMemoryDirectory": ".claude/memory"` in `~/.claude/settings.json`.

## CLD030: plansDirectory is not set to ".claude/plans"

**Severity:** info

Pins where Claude Code writes implementation plans.

**Remediation:** Set `"plansDirectory": ".claude/plans"` in `~/.claude/settings.json`.

## CLD031: sandbox.enabled is not true

**Severity:** high

`sandbox.enabled=true` turns on the platform sandbox that isolates tool execution from the rest of the system.

**Remediation:** Set `"sandbox": {"enabled": true, ...}` in `~/.claude/settings.json`.

## CLD032: sandbox.autoAllowBashIfSandboxed is not false

**Severity:** high

`false` keeps Bash calls in the explicit-permission flow even when sandboxed.

**Remediation:** Set `"sandbox": {"autoAllowBashIfSandboxed": false, ...}` in `~/.claude/settings.json`.

## CLD033: sandbox.allowUnsandboxedCommands is not false

**Severity:** high

`false` prevents tool calls from running outside the sandbox; isolation becomes mandatory rather than optional.

**Remediation:** Set `"sandbox": {"allowUnsandboxedCommands": false, ...}` in `~/.claude/settings.json`.

## CLD034: sandbox.network.allowManagedDomainsOnly is not true

**Severity:** warn

`true` restricts outbound network calls to the explicitly allowed domain list.

**Remediation:** Set `"sandbox": {"network": {"allowManagedDomainsOnly": true, ...}}` in `~/.claude/settings.json`.

## CLD035: permissions.disableBypassPermissionsMode is not "disable"

**Severity:** high

`"disable"` prevents Claude Code from entering bypass mode (which otherwise lets tool calls run without user confirmation).

**Remediation:** Set `"permissions": {"disableBypassPermissionsMode": "disable"}` in `~/.claude/settings.json`.

## CLD036: permissions.deny missing network/exfiltration tool blocks

**Severity:** high

`permissions.deny` should block `nc`, `netcat`, `socat`, `ssh`, `scp`, `rsync`. These tools can move data off-host or open inbound connections from a sandbox-escaped tool call.

**Remediation:** Add to `permissions.deny` in `~/.claude/settings.json`:
```json
"Bash(nc:*)", "Bash(netcat:*)", "Bash(socat:*)",
"Bash(ssh:*)", "Bash(scp:*)", "Bash(rsync:*)"
```

## CLD037: permissions.deny missing destructive filesystem command blocks

**Severity:** high

Blocks `chmod 777`, `chown`, `rm -rf /`, `rm -rf ~`, `dd`, `mkfs` from being invoked by tool calls.

**Remediation:** Add to `permissions.deny`:
```json
"Bash(chmod 777:*)", "Bash(chown:*)",
"Bash(rm -rf /:*)", "Bash(rm -rf ~:*)",
"Bash(dd:*)", "Bash(mkfs:*)"
```

## CLD038: permissions.deny missing destructive git command blocks

**Severity:** warn

Blocks `git push`, `git tag`, `git reset --hard` so publishing and destructive resets happen out-of-band.

**Remediation:** Add to `permissions.deny`:
```json
"Bash(git push:*)", "Bash(git tag:*)", "Bash(git reset --hard:*)"
```

## CLD039: permissions.deny missing home credential directory blocks

**Severity:** high

Blocks reads of `~/.ssh`, `~/.aws`, `~/.gnupg`, `~/.config/gh`, `~/.kube`, `~/.docker/config.json` so credential material stays out of transcripts.

**Remediation:** Add to `permissions.deny`:
```json
"Read(~/.ssh/**)", "Read(~/.aws/**)", "Read(~/.gnupg/**)",
"Read(~/.config/gh/**)", "Read(~/.kube/**)", "Read(~/.docker/config.json)"
```

## CLD040: permissions.deny missing project secret file blocks

**Severity:** high

Blocks reads of `.env*`, `*.pem`, `*.key`, `id_rsa*`, `id_ed25519*`, `credentials*` at both top-level and `**`-glob depth.

**Remediation:** Add to `permissions.deny`:
```json
"Read(./.env)", "Read(./.env.*)", "Read(./*.pem)", "Read(./*.key)",
"Read(./**/.env)", "Read(./**/.env.*)", "Read(./**/*.pem)", "Read(./**/*.key)",
"Read(./**/id_rsa*)", "Read(./**/id_ed25519*)", "Read(./**/credentials*)"
```

## CLD041: Claude Desktop installed Native Messaging manifests for Chromium browsers

**Severity:** high

macOS only. Claude Desktop silently installs a Native Messaging manifest at `com.anthropic.claude_browser_extension.json` under `~/Library/Application Support` for several Chromium-based browsers — including browsers that are not even installed. The manifest pre-authorizes a small set of Chrome extension IDs to invoke a Claude Desktop helper binary outside the browser sandbox. The file is rewritten on every Claude Desktop launch. Mitigate by truncating each manifest to an empty file and setting the user-immutable (`uchg`) flag.

**Remediation:**
```bash
find ~/Library/Application\ Support -type f \
  -name com.anthropic.claude_browser_extension.json \
  -exec sh -c ': > "$1" && chflags uchg "$1"' _ {} \;
```
