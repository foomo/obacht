# Rules Overview

obacht includes 98 built-in rules across 12 categories.

| ID     | Category    | Title                                                | Severity |
|--------|-------------|------------------------------------------------------|----------|
| CLD001 | Claude      | Global gitignore does not exclude Claude Code local settings | warn     |
| CLD002 | Claude      | Claude Code autoCompactEnabled is not disabled       | warn     |
| CLD003 | Claude      | Claude Code prStatusFooterEnabled is not disabled    | warn     |
| CLD004 | Claude      | Claude Code claudeInChromeDefaultEnabled is not disabled | warn     |
| CLD005 | Claude      | Claude Code sandbox.failIfUnavailable is not enabled | high     |
| CLD006 | Claude      | Claude Code DISABLE_COMPACT is not set in settings.json | warn     |
| CLD007 | Claude      | Claude Code DISABLE_TELEMETRY is not set in settings.json | warn     |
| CLD008 | Claude      | Claude Code DISABLE_BUG_COMMAND is not set in settings.json | warn     |
| CLD009 | Claude      | Claude Code DISABLE_AUTO_COMPACT is not set in settings.json | warn     |
| CLD010 | Claude      | Claude Code DISABLE_LOGIN_COMMAND is not set in settings.json | warn     |
| CLD011 | Claude      | Claude Code DISABLE_LOGOUT_COMMAND is not set in settings.json | warn     |
| CLD012 | Claude      | Claude Code DISABLE_ERROR_REPORTING is not set in settings.json | warn     |
| CLD013 | Claude      | Claude Code DISABLE_UPGRADE_COMMAND is not set in settings.json | warn     |
| CLD014 | Claude      | Claude Code DISABLE_FEEDBACK_COMMAND is not set in settings.json | warn     |
| CLD015 | Claude      | Claude Code DISABLE_EXTRA_USAGE_COMMAND is not set in settings.json | warn     |
| CLD016 | Claude      | Claude Code CLAUDE_CODE_DISABLE_FAST_MODE is not set in settings.json | warn     |
| CLD017 | Claude      | Claude Code DISABLE_INSTALL_GITHUB_APP_COMMAND is not set in settings.json | warn     |
| CLD018 | Claude      | Claude Code CLAUDE_CODE_DISABLE_CRON is not set in settings.json | warn     |
| CLD019 | Claude      | Claude Code CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY is not set in settings.json | warn     |
| CLD020 | Claude      | Claude Code CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING is not set in settings.json | warn     |
| CLD021 | Claude      | Claude Code CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS is not set in settings.json | warn     |
| CLD022 | Claude      | Claude Code FORCE_AUTOUPDATE_PLUGINS is not set in settings.json | warn     |
| CLD023 | Claude      | Claude Code IS_DEMO is not set in settings.json      | warn     |
| CLD024 | Claude      | Claude Code disableAutoMode is not set to "disable" in settings.json | warn     |
| CLD025 | Claude      | Claude Code disableDeepLinkRegistration is not set to "disable" in settings.json | warn     |
| CLD026 | Claude      | Claude Code attribution.commit and attribution.pr are not empty in settings.json | warn     |
| CLD027 | Claude      | Claude Code respectGitignore is not enabled in settings.json | warn     |
| CLD028 | Claude      | Claude Code skipWebFetchPreflight is not enabled in settings.json | warn     |
| CLD029 | Claude      | Claude Code autoMemoryDirectory is not set to ".claude/memory" in settings.json | warn     |
| CLD030 | Claude      | Claude Code plansDirectory is not set to ".claude/plans" in settings.json | warn     |
| CLD031 | Claude      | Claude Code sandbox.enabled is not true in settings.json | warn     |
| CLD032 | Claude      | Claude Code sandbox.autoAllowBashIfSandboxed is not false in settings.json | warn     |
| CLD033 | Claude      | Claude Code sandbox.allowUnsandboxedCommands is not false in settings.json | warn     |
| CLD034 | Claude      | Claude Code sandbox.network.allowManagedDomainsOnly is not true in settings.json | warn     |
| CLD035 | Claude      | Claude Code sandbox.network.allowedDomains does not include "github.com" | warn     |
| CLD036 | Claude      | Claude Code sandbox.network.deniedDomains does not include "uploads.github.com" | warn     |
| CLD037 | Claude      | Claude Code sandbox.filesystem.allowWrite does not include "~/.cache/npm/logs" | warn     |
| CLD038 | Claude      | Claude Code sandbox.filesystem.allowWrite does not include "~/.config/claude/debug" | warn     |
| CLD039 | Claude      | Claude Code permissions.disableBypassPermissionsMode is not "disable" in settings.json | high     |
| CLD040 | Claude      | Claude Code permissions.deny missing network/exfiltration tool blocks | high     |
| CLD041 | Claude      | Claude Code permissions.deny missing destructive filesystem command blocks | high     |
| CLD042 | Claude      | Claude Code permissions.deny missing destructive git command blocks | warn     |
| CLD043 | Claude      | Claude Code permissions.deny missing home credential directory blocks | high     |
| CLD044 | Claude      | Claude Code permissions.deny missing project secret file blocks | high     |
| CRD001 | Credentials | AWS credentials file has weak permissions            | high     |
| CRD002 | Credentials | .netrc file has weak permissions                     | high     |
| CRD003 | Credentials | GCP credentials file has weak permissions            | high     |
| CRD004 | Credentials | .npmrc with auth token has weak permissions          | high     |
| DOC001 | Docker      | Docker socket has overly permissive access           | high     |
| DOC002 | Docker      | User is in the docker group                          | warn     |
| ENV001 | Environment | Sensitive credentials found in environment variables | high     |
| GIT001 | Git         | Git credential helper stores passwords in plaintext  | high     |
| GIT002 | Git         | Git commit signing is not enabled                    | warn     |
| GIT003 | Git         | Git safe.directory set to wildcard                   | high     |
| GIT004 | Git         | Global gitignore does not exclude .env files         | warn     |
| KUB001 | Kubernetes  | Kubeconfig has weak permissions                      | high     |
| KUB002 | Kubernetes  | Production Kubernetes context is active              | warn     |
| OS001  | OS          | System Integrity Protection is disabled              | critical |
| OS002  | OS          | FileVault disk encryption is disabled                | critical |
| OS003  | OS          | Application Firewall is disabled                     | high     |
| OS004  | OS          | Stealth Mode is disabled                             | high     |
| OS005  | OS          | Gatekeeper is disabled                               | critical |
| OS006  | OS          | Automatic login is enabled                           | high     |
| OS007  | OS          | Guest account is enabled                             | high     |
| OS008  | OS          | Screen lock timeout exceeds 5 minutes                | warn     |
| OS009  | OS          | Automatic OS updates are disabled                    | high     |
| OS010  | OS          | Automatic App Store updates are disabled             | warn     |
| OS011  | OS          | Rapid Security Responses are disabled                | high     |
| OS013  | OS          | Screen Sharing is enabled                            | high     |
| OS014  | OS          | Internet Sharing is enabled                          | high     |
| OS015  | OS          | Printer Sharing is enabled                           | warn     |
| OS016  | OS          | Remote Apple Events are enabled                      | high     |
| OS017  | OS          | AirDrop is set to Everyone                           | high     |
| OS018  | OS          | No EDR agent deployed                                | warn     |
| OS019  | OS          | Legacy kernel extensions are not blocked             | warn     |
| OS020  | OS          | Device is not enrolled in MDM                        | high     |
| OS021  | OS          | Rosetta 2 is installed                               | info     |
| OS022  | OS          | AirDrop is not fully disabled                        | info     |
| OS023  | OS          | Time Machine backup is disabled                      | warn     |
| OS024  | OS          | Remote Login (SSH server) is enabled                 | high     |
| OS025  | OS          | Remote Management is enabled                         | high     |
| OS026  | OS          | Bluetooth Sharing is enabled                         | warn     |
| OS027  | OS          | Media Sharing is enabled                             | warn     |
| OS028  | OS          | File Sharing (SMB) is enabled                        | warn     |
| OS029  | OS          | Content Caching is enabled                           | warn     |
| PTH001 | PATH        | World-writable directory in PATH                     | high     |
| PTH002 | PATH        | Relative path entry in PATH                          | warn     |
| PRV001 | Privacy     | No password manager application detected             | warn     |
| PRV002 | Privacy     | No VPN configuration detected                        | info     |
| PRV003 | Privacy     | Encrypted DNS is not configured                      | warn     |
| PRV004 | Privacy     | Untrusted DNS resolver is configured                 | warn     |
| SHL001 | Shell       | Shell history file has weak permissions              | warn     |
| SSH001 | SSH         | SSH private key has weak permissions                 | high     |
| SSH002 | SSH         | SSH directory has weak permissions                   | high     |
| SSH003 | SSH         | SSH StrictHostKeyChecking is disabled                | high     |
| SSH004 | SSH         | SSH agent forwarding is enabled globally             | warn     |
| TOL001 | Tools       | Security-relevant tool is missing                    | info     |
| TOL002 | Tools       | Homebrew auto-update is disabled                     | warn     |
