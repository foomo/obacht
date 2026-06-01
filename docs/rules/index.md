# Rules Overview

obacht includes 118 built-in rules across 13 categories.

| ID     | Category    | Title                                                                     | Severity |
|--------|-------------|---------------------------------------------------------------------------|----------|
| CLD001 | Claude      | Global gitignore does not exclude Claude Code local settings              | warn     |
| CLD002 | Claude      | autoCompactEnabled is not disabled                                        | warn     |
| CLD003 | Claude      | prStatusFooterEnabled is not disabled                                     | info     |
| CLD004 | Claude      | claudeInChromeDefaultEnabled is not disabled                              | warn     |
| CLD005 | Claude      | sandbox.failIfUnavailable is not enabled                                  | high     |
| CLD006 | Claude      | DISABLE_COMPACT is not set                                                | warn     |
| CLD007 | Claude      | DISABLE_TELEMETRY is not set                                              | high     |
| CLD008 | Claude      | DISABLE_BUG_COMMAND is not set                                            | warn     |
| CLD009 | Claude      | DISABLE_AUTO_COMPACT is not set                                           | warn     |
| CLD010 | Claude      | DISABLE_LOGIN_COMMAND is not set                                          | info     |
| CLD011 | Claude      | DISABLE_LOGOUT_COMMAND is not set                                         | info     |
| CLD012 | Claude      | DISABLE_ERROR_REPORTING is not set                                        | high     |
| CLD013 | Claude      | DISABLE_UPGRADE_COMMAND is not set                                        | warn     |
| CLD014 | Claude      | DISABLE_FEEDBACK_COMMAND is not set                                       | warn     |
| CLD015 | Claude      | DISABLE_EXTRA_USAGE_COMMAND is not set                                    | warn     |
| CLD016 | Claude      | CLAUDE_CODE_DISABLE_FAST_MODE is not set                                  | warn     |
| CLD017 | Claude      | DISABLE_INSTALL_GITHUB_APP_COMMAND is not set                             | warn     |
| CLD018 | Claude      | CLAUDE_CODE_DISABLE_CRON is not set                                       | info     |
| CLD019 | Claude      | CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY is not set                            | warn     |
| CLD020 | Claude      | CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING is not set                         | info     |
| CLD021 | Claude      | CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS is not set                         | info     |
| CLD022 | Claude      | FORCE_AUTOUPDATE_PLUGINS is not set                                       | warn     |
| CLD023 | Claude      | IS_DEMO is not set                                                        | warn     |
| CLD024 | Claude      | disableAutoMode is not "disable"                                          | warn     |
| CLD025 | Claude      | disableDeepLinkRegistration is not "disable"                              | high     |
| CLD026 | Claude      | attribution.commit and attribution.pr are not empty                       | info     |
| CLD027 | Claude      | respectGitignore is not enabled                                           | warn     |
| CLD028 | Claude      | skipWebFetchPreflight is not enabled                                      | warn     |
| CLD029 | Claude      | autoMemoryDirectory is not ".claude/memory"                               | info     |
| CLD030 | Claude      | plansDirectory is not ".claude/plans"                                     | info     |
| CLD031 | Claude      | sandbox.enabled is not true                                               | high     |
| CLD032 | Claude      | sandbox.autoAllowBashIfSandboxed is not false                             | high     |
| CLD033 | Claude      | sandbox.allowUnsandboxedCommands is not false                             | high     |
| CLD034 | Claude      | sandbox.network.allowManagedDomainsOnly is not true                       | warn     |
| CLD035 | Claude      | permissions.disableBypassPermissionsMode is not "disable"                 | high     |
| CLD036 | Claude      | permissions.deny missing network/exfiltration tool blocks                 | high     |
| CLD037 | Claude      | permissions.deny missing destructive filesystem command blocks            | high     |
| CLD038 | Claude      | permissions.deny missing destructive git command blocks                   | warn     |
| CLD039 | Claude      | permissions.deny missing home credential directory blocks                 | high     |
| CLD040 | Claude      | permissions.deny missing project secret file blocks                       | high     |
| CLD041 | Claude      | Claude Desktop installed Native Messaging manifests for Chromium browsers | high     |
| CRD001 | Credentials | AWS credentials file has weak permissions                                 | high     |
| CRD002 | Credentials | .netrc file has weak permissions                                          | high     |
| CRD003 | Credentials | GCP credentials file has weak permissions                                 | high     |
| CRD004 | Credentials | .npmrc with auth token has weak permissions                               | high     |
| DOC001 | Docker      | Docker socket has overly permissive access                                | high     |
| DOC002 | Docker      | User is in the docker group                                               | warn     |
| ENV001 | Environment | Sensitive credentials found in environment variables                      | high     |
| GIT001 | Git         | Git credential helper stores passwords in plaintext                       | high     |
| GIT002 | Git         | Git commit signing is not enabled                                         | warn     |
| GIT003 | Git         | Git safe.directory set to wildcard                                        | high     |
| GIT004 | Git         | Global gitignore does not exclude .env files                              | warn     |
| KUB001 | Kubernetes  | Kubeconfig has weak permissions                                           | high     |
| KUB002 | Kubernetes  | Production Kubernetes context is active                                   | warn     |
| OS001  | OS          | System Integrity Protection is disabled                                   | critical |
| OS002  | OS          | FileVault disk encryption is disabled                                     | critical |
| OS003  | OS          | Application Firewall is disabled                                          | high     |
| OS004  | OS          | Stealth Mode is disabled                                                  | high     |
| OS005  | OS          | Gatekeeper is disabled                                                    | critical |
| OS006  | OS          | Automatic login is enabled                                                | high     |
| OS007  | OS          | Guest account is enabled                                                  | high     |
| OS008  | OS          | Screen lock timeout exceeds 5 minutes                                     | warn     |
| OS009  | OS          | Automatic OS updates are disabled                                         | high     |
| OS010  | OS          | Automatic App Store updates are disabled                                  | warn     |
| OS011  | OS          | Rapid Security Responses are disabled                                     | high     |
| OS013  | OS          | Screen Sharing is enabled                                                 | high     |
| OS014  | OS          | Internet Sharing is enabled                                               | high     |
| OS015  | OS          | Printer Sharing is enabled                                                | warn     |
| OS016  | OS          | Remote Apple Events are enabled                                           | high     |
| OS017  | OS          | AirDrop is set to Everyone                                                | high     |
| OS018  | OS          | No EDR agent deployed                                                     | info     |
| OS019  | OS          | Legacy kernel extensions are not blocked                                  | info     |
| OS020  | OS          | Device is not enrolled in MDM                                             | high     |
| OS021  | OS          | Rosetta 2 is installed                                                    | info     |
| OS022  | OS          | AirDrop is not fully disabled                                             | info     |
| OS023  | OS          | Time Machine backup is disabled                                           | warn     |
| OS024  | OS          | Remote Login (SSH server) is enabled                                      | high     |
| OS025  | OS          | Remote Management is enabled                                              | high     |
| OS026  | OS          | Bluetooth Sharing is enabled                                              | warn     |
| OS027  | OS          | Media Sharing is enabled                                                  | warn     |
| OS028  | OS          | File Sharing (SMB) is enabled                                             | warn     |
| OS029  | OS          | Content Caching is enabled                                                | warn     |
| OS030  | OS          | Current user has local admin privileges                                   | warn     |
| OS031  | OS          | Password not required immediately after screen lock                       | high     |
| OS032  | OS          | Time Machine destination is not encrypted                                 | warn     |
| OS033  | OS          | Time Machine has no recent backup                                         | warn     |
| OS034  | OS          | AirPlay Receiver is enabled                                               | warn     |
| OS035  | OS          | Automatic download of OS updates is disabled                              | warn     |
| OS036  | OS          | macOS major version is unsupported                                        | warn     |
| OS037  | OS          | App Store in-app review prompts are enabled                               | info     |
| PTH001 | PATH        | World-writable directory in PATH                                          | high     |
| PTH002 | PATH        | Relative path entry in PATH                                               | warn     |
| PRV001 | Privacy     | No password manager application detected                                  | warn     |
| PRV002 | Privacy     | No VPN configuration detected                                             | info     |
| PRV003 | Privacy     | Encrypted DNS is not configured                                           | warn     |
| PRV004 | Privacy     | Untrusted DNS resolver is configured                                      | warn     |
| PRV005 | Privacy     | DO_NOT_TRACK env var is not set                                           | info     |
| SHL001 | Shell       | Shell history file has weak permissions                                   | warn     |
| BUM001 | Bumblebee   | Compromised npm package present on disk                                   | critical |
| BUM002 | Bumblebee   | Compromised PyPI package present on disk                                  | critical |
| BUM003 | Bumblebee   | Compromised Go module present on disk                                     | critical |
| BUM004 | Bumblebee   | Compromised RubyGem present on disk                                       | critical |
| BUM005 | Bumblebee   | Compromised Composer package present on disk                              | critical |
| BUM006 | Bumblebee   | Compromised MCP server configured                                         | critical |
| BUM007 | Bumblebee   | Compromised editor extension installed                                    | critical |
| BUM008 | Bumblebee   | Compromised browser extension installed                                   | critical |
| SSH001 | SSH         | SSH private key has weak permissions                                      | high     |
| SSH002 | SSH         | SSH directory has weak permissions                                        | high     |
| SSH003 | SSH         | SSH StrictHostKeyChecking is disabled                                     | high     |
| SSH004 | SSH         | SSH agent forwarding is enabled globally                                  | warn     |
| SSH005 | SSH         | SSH key uses weak algorithm                                               | high     |
| TOL001 | Tools       | Security-relevant tool is missing                                         | info     |
| TOL002 | Tools       | Homebrew auto-update is disabled                                          | warn     |
| TOL003 | Tools       | Package manager metadata is stale                                         | warn     |
| TOL004 | Tools       | Homebrew analytics is not disabled                                        | warn     |
| TOL005 | Tools       | Go telemetry is not disabled                                              | warn     |
