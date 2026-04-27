# Claude Rules

Rules covering [Claude Code](https://docs.claude.com/en/docs/claude-code) CLI configuration. These checks read
`~/.claude.json` (or `$CLAUDE_CONFIG_DIR/.claude.json`) and the global gitignore. The recommended hardening turns off
features that increase attack surface, leak metadata, or silently rewrite working state.

## CLD001: Global gitignore does not exclude Claude Code local settings

**Severity:** warn

Claude Code stores user-specific settings in `.claude/settings.local.json` which may contain personal configuration. A
global gitignore entry for `**/.claude/settings.local.json` prevents accidental commits.

**What it checks:**

- Whether the user has a global gitignore configured (`git config --global core.excludesfile`)
- Whether that file contains `**/.claude/settings.local.json`

**Remediation:**

```bash
echo '**/.claude/settings.local.json' >> ~/.gitignore_global
git config --global core.excludesfile ~/.gitignore_global
```

## CLD002: Claude Code autoCompactEnabled is not disabled

**Severity:** warn

Claude Code's auto-compact feature silently rewrites conversation context when token budgets are reached, which can
drop critical details. Disabling it gives the user explicit control over when compaction happens.

**What it checks:**

- `autoCompactEnabled` in `~/.claude.json` (or `$CLAUDE_CONFIG_DIR/.claude.json`) is set to `false`

**Remediation:**

```json
{
  "autoCompactEnabled": false
}
```

## CLD003: Claude Code prStatusFooterEnabled is not disabled

**Severity:** warn

The PR status footer makes outbound calls to GitHub to display PR status in the CLI footer. Disabling it reduces
network chatter and avoids leaking working-directory metadata to remote services when not needed.

**What it checks:**

- `prStatusFooterEnabled` in `~/.claude.json` (or `$CLAUDE_CONFIG_DIR/.claude.json`) is set to `false`

**Remediation:**

```json
{
  "prStatusFooterEnabled": false
}
```

## CLD004: Claude Code claudeInChromeDefaultEnabled is not disabled

**Severity:** warn

The Claude-in-Chrome integration injects Claude into the default browser. Disabling the default-enable flag prevents
unintended browser instrumentation and keeps Claude Code scoped to the terminal until explicitly opted in.

**What it checks:**

- `claudeInChromeDefaultEnabled` in `~/.claude.json` (or `$CLAUDE_CONFIG_DIR/.claude.json`) is set to `false`

**Remediation:**

```json
{
  "claudeInChromeDefaultEnabled": false
}
```

## CLD005: Claude Code sandbox.failIfUnavailable is not enabled

**Severity:** high

When `sandbox.failIfUnavailable` is false (or unset), Claude Code will silently run without the platform sandbox if it
cannot be initialized, removing a key isolation boundary. Setting it to true forces Claude Code to fail closed when
the sandbox is unavailable, ensuring tool execution always runs under the expected confinement.

**What it checks:**

- `sandbox.failIfUnavailable` in `~/.claude.json` (or `$CLAUDE_CONFIG_DIR/.claude.json`) is set to `true`

**Remediation:**

```json
{
  "sandbox": {
    "failIfUnavailable": true
  }
}
```

## CLD006-CLD023: Claude Code env hardening in settings.json

The following rules each check that a single environment variable is set to `"1"` in the `env` block of
`~/.claude/settings.json` (or `$CLAUDE_CONFIG_DIR/settings.json`). Together they harden Claude Code against
data-leak, accidental command exposure, and silent state mutation. Each rule passes when its variable equals `"1"`
and fails otherwise (including when the variable is unset).

A complete remediation block looks like:

```json
{
  "env": {
    "DISABLE_COMPACT": "1",
    "DISABLE_TELEMETRY": "1",
    "DISABLE_BUG_COMMAND": "1",
    "DISABLE_AUTO_COMPACT": "1",
    "DISABLE_LOGIN_COMMAND": "1",
    "DISABLE_LOGOUT_COMMAND": "1",
    "DISABLE_ERROR_REPORTING": "1",
    "DISABLE_UPGRADE_COMMAND": "1",
    "DISABLE_FEEDBACK_COMMAND": "1",
    "DISABLE_EXTRA_USAGE_COMMAND": "1",
    "CLAUDE_CODE_DISABLE_FAST_MODE": "1",
    "DISABLE_INSTALL_GITHUB_APP_COMMAND": "1",
    "CLAUDE_CODE_DISABLE_CRON": "1",
    "CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY": "1",
    "CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING": "1",
    "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS": "1",
    "FORCE_AUTOUPDATE_PLUGINS": "1",
    "IS_DEMO": "1"
  }
}
```

### CLD006: DISABLE_COMPACT

**Severity:** warn — Turns off the manual `/compact` command, removing a class of silent context rewrites.

### CLD007: DISABLE_TELEMETRY

**Severity:** warn — Opts out of anonymized usage telemetry to Anthropic.

### CLD008: DISABLE_BUG_COMMAND

**Severity:** warn — Hides the `/bug` command, preventing accidental upload of working session content.

### CLD009: DISABLE_AUTO_COMPACT

**Severity:** warn — Disables automatic conversation compaction triggered when context fills up.

### CLD010: DISABLE_LOGIN_COMMAND

**Severity:** warn — Hides the `/login` command. Prevents accidental account switching from inside an active session.

### CLD011: DISABLE_LOGOUT_COMMAND

**Severity:** warn — Hides the `/logout` command. Pairs with `DISABLE_LOGIN_COMMAND`.

### CLD012: DISABLE_ERROR_REPORTING

**Severity:** warn — Stops Claude Code from uploading error stack traces (which can contain working-directory metadata).

### CLD013: DISABLE_UPGRADE_COMMAND

**Severity:** warn — Hides the in-app upgrade command so Claude Code is upgraded only via the system package manager.

### CLD014: DISABLE_FEEDBACK_COMMAND

**Severity:** warn — Hides the `/feedback` command, preventing accidental upload of free-form context to Anthropic.

### CLD015: DISABLE_EXTRA_USAGE_COMMAND

**Severity:** warn — Hides the extra usage reporting command.

### CLD016: CLAUDE_CODE_DISABLE_FAST_MODE

**Severity:** warn — Turns off the `/fast` toggle so the active model cannot be flipped mid-session.

### CLD017: DISABLE_INSTALL_GITHUB_APP_COMMAND

**Severity:** warn — Prevents Claude Code from installing GitHub Apps from inside a session.

### CLD018: CLAUDE_CODE_DISABLE_CRON

**Severity:** warn — Turns off scheduled/cron-style autonomous task execution.

### CLD019: CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY

**Severity:** warn — Suppresses periodic in-session feedback prompts.

### CLD020: CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING

**Severity:** warn — Disables file-checkpointing snapshots, reducing disk/IO and removing stale-snapshot risk.

### CLD021: CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS

**Severity:** warn — Opts out of experimental beta features for reproducibility and auditability.

### CLD022: FORCE_AUTOUPDATE_PLUGINS

**Severity:** warn — Forces plugin auto-update so installed plugins do not drift behind upstream fixes.

### CLD023: IS_DEMO

**Severity:** warn — Enables demo-mode affordances. Reflects a deliberate user-policy choice; remove the rule (or
override the env value) if your environment should not run in demo mode.

## CLD024-CLD030: Claude Code top-level settings.json hardening

The following rules check top-level keys (outside the `env` block) in `~/.claude/settings.json` (or
`$CLAUDE_CONFIG_DIR/settings.json`). They harden Claude Code against autonomous-mode escalation, deep-link
registration, attribution leakage, and silent reads of gitignored files.

A complete remediation block looks like:

```json
{
  "disableAutoMode": "disable",
  "disableDeepLinkRegistration": "disable",
  "attribution": {
    "commit": "",
    "pr": ""
  },
  "respectGitignore": true,
  "skipWebFetchPreflight": true,
  "autoMemoryDirectory": ".claude/memory",
  "plansDirectory": ".claude/plans"
}
```

### CLD024: disableAutoMode

**Severity:** warn — Prevents Claude Code from automatically switching into autonomous modes. Forces explicit user
opt-in per session.

### CLD025: disableDeepLinkRegistration

**Severity:** warn — Stops Claude Code from registering itself as the handler for `claude://` deep links. Removes a
class of inbound URL handlers other applications could otherwise invoke.

### CLD026: attribution.commit / attribution.pr empty

**Severity:** warn — Removes the attribution trailer Claude Code appends to commits and PRs. Keeps authored artefacts
free of tool-specific markers. Fails when either subkey is missing or non-empty.

### CLD027: respectGitignore

**Severity:** warn — Honours `.gitignore` when reading files, preventing ignored content (credentials, build outputs,
vendored data) from being pulled into the conversation.

### CLD028: skipWebFetchPreflight

**Severity:** warn — Skips the OPTIONS preflight before WebFetch, removing an extra outbound request and avoiding
working-host metadata leakage.

### CLD029: autoMemoryDirectory

**Severity:** warn — Pins Claude Code's automatic memory store to `.claude/memory` so memory files live alongside the
rest of the `.claude` directory.

### CLD030: plansDirectory

**Severity:** warn — Pins where Claude Code writes implementation plans to `.claude/plans` for unified review and
management.

## CLD031-CLD038: Claude Code sandbox configuration in settings.json

The following rules check the `sandbox` block in `~/.claude/settings.json` (or `$CLAUDE_CONFIG_DIR/settings.json`).
Together they enforce: the sandbox is on, it is not auto-bypassed for Bash, unsandboxed commands are not allowed,
network access is restricted to a managed allowlist that includes `github.com` while explicitly denying
`uploads.github.com`, and the filesystem allowlist contains the npm-log and Claude debug paths needed for normal
tool operation.

Array checks use must-contain semantics — adding extra domains or write paths does not fail the rule, only the
absence of the listed value does.

A complete remediation block looks like:

```json
{
  "sandbox": {
    "enabled": true,
    "autoAllowBashIfSandboxed": false,
    "allowUnsandboxedCommands": false,
    "network": {
      "allowManagedDomainsOnly": true,
      "allowedDomains": ["github.com"],
      "deniedDomains": ["uploads.github.com"]
    },
    "filesystem": {
      "allowWrite": [
        "~/.cache/npm/logs",
        "~/.config/claude/debug"
      ]
    }
  }
}
```

### CLD031: sandbox.enabled

**Severity:** warn — Sandbox must be on. Without it, tool calls run with the user's full permissions.

### CLD032: sandbox.autoAllowBashIfSandboxed

**Severity:** warn — Bash calls must stay in the explicit-permission flow even when sandboxed; auto-approval bypasses
per-command oversight.

### CLD033: sandbox.allowUnsandboxedCommands

**Severity:** warn — The sandbox boundary must be enforced; unsandboxed escape hatches make isolation optional.

### CLD034: sandbox.network.allowManagedDomainsOnly

**Severity:** warn — Outbound traffic must be limited to the allowlist, otherwise the allowlist is decorative.

### CLD035: sandbox.network.allowedDomains contains "github.com"

**Severity:** warn — `github.com` is required for clone, `gh`, and API calls under the network allowlist.

### CLD036: sandbox.network.deniedDomains contains "uploads.github.com"

**Severity:** warn — Blocking the upload subdomain prevents using GitHub's allowlisted access to write data outbound.

### CLD037: sandbox.filesystem.allowWrite contains "~/.cache/npm/logs"

**Severity:** warn — npm operations need to write log files; without this entry npm-driven tools fail under the
sandbox.

### CLD038: sandbox.filesystem.allowWrite contains "~/.config/claude/debug"

**Severity:** warn — Claude Code's debug artefacts need a writable location; without this entry diagnostics cannot be
persisted.

## CLD039-CLD044: Claude Code permissions deny-list hardening in settings.json

The following rules check the `permissions` block in `~/.claude/settings.json` (or
`$CLAUDE_CONFIG_DIR/settings.json`). Together they harden Claude Code by forcing every privileged tool call through
the explicit permission flow (no bypass mode) and by blocking categories of commands and read paths that the
assistant has no legitimate reason to invoke: network exfiltration tools, destructive filesystem commands,
publishing/destructive git commands, home credential directories, and project-local secret files.

Array checks use must-contain semantics — adding extra deny entries does not fail the rule, only the absence of the
listed value does.

A complete remediation block looks like:

```json
{
  "permissions": {
    "disableBypassPermissionsMode": "disable",
    "deny": [
      "Bash(nc:*)",
      "Bash(netcat:*)",
      "Bash(socat:*)",
      "Bash(ssh:*)",
      "Bash(scp:*)",
      "Bash(rsync:*)",

      "Bash(chmod 777:*)",
      "Bash(chown:*)",
      "Bash(rm -rf /:*)",
      "Bash(rm -rf ~:*)",
      "Bash(dd:*)",
      "Bash(mkfs:*)",

      "Bash(git push:*)",
      "Bash(git tag:*)",
      "Bash(git reset --hard:*)",

      "Read(~/.ssh/**)",
      "Read(~/.aws/**)",
      "Read(~/.gnupg/**)",
      "Read(~/.config/gh/**)",
      "Read(~/.kube/**)",
      "Read(~/.docker/config.json)",

      "Read(./.env)",
      "Read(./.env.*)",
      "Read(./*.pem)",
      "Read(./*.key)",
      "Read(./**/.env)",
      "Read(./**/.env.*)",
      "Read(./**/*.pem)",
      "Read(./**/*.key)",
      "Read(./**/id_rsa*)",
      "Read(./**/id_ed25519*)",
      "Read(./**/credentials*)"
    ]
  }
}
```

### CLD039: permissions.disableBypassPermissionsMode

**Severity:** high — Bypass mode turns off the permission-prompt flow and lets tool calls run without user
confirmation. Setting it to `"disable"` forces every privileged action through normal allow/deny review.

### CLD040: permissions.deny network/exfiltration tool blocks

**Severity:** high — Required entries: `Bash(nc:*)`, `Bash(netcat:*)`, `Bash(socat:*)`, `Bash(ssh:*)`,
`Bash(scp:*)`, `Bash(rsync:*)`. These move data off-host or open inbound connections from a sandbox-escaped tool call.
Claude Code's sanctioned network channels (WebFetch, `gh`) cover the legitimate use cases.

### CLD041: permissions.deny destructive filesystem command blocks

**Severity:** high — Required entries: `Bash(chmod 777:*)`, `Bash(chown:*)`, `Bash(rm -rf /:*)`, `Bash(rm -rf ~:*)`,
`Bash(dd:*)`, `Bash(mkfs:*)`. Blocks mass deletion, ownership changes, raw block writes, and filesystem creation —
all irreversible operations that should never come from an autonomous tool call.

### CLD042: permissions.deny destructive git command blocks

**Severity:** warn — Required entries: `Bash(git push:*)`, `Bash(git tag:*)`, `Bash(git reset --hard:*)`. Forces the
user to publish history or discard work out-of-band rather than from inside a tool call.

### CLD043: permissions.deny home credential directory blocks

**Severity:** high — Required entries: `Read(~/.ssh/**)`, `Read(~/.aws/**)`, `Read(~/.gnupg/**)`,
`Read(~/.config/gh/**)`, `Read(~/.kube/**)`, `Read(~/.docker/config.json)`. Keeps long-lived credentials (SSH keys,
IAM access keys, PGP keys, GitHub CLI tokens, cluster credentials, registry tokens) out of transcripts.

### CLD044: permissions.deny project secret file blocks

**Severity:** high — Required entries: `Read(./.env)`, `Read(./.env.*)`, `Read(./*.pem)`, `Read(./*.key)`,
`Read(./**/.env)`, `Read(./**/.env.*)`, `Read(./**/*.pem)`, `Read(./**/*.key)`, `Read(./**/id_rsa*)`,
`Read(./**/id_ed25519*)`, `Read(./**/credentials*)`. Both top-level and recursive (`**` glob) variants are
required so secrets in nested directories cannot be read either.
