# Rules Overview

obacht includes 59 built-in rules across 12 categories.

| ID     | Category    | Title                                                | Severity |
|--------|-------------|------------------------------------------------------|----------|
| CLD001 | Claude      | Global gitignore does not exclude Claude Code local settings | warn     |
| CLD002 | Claude      | Claude Code autoCompactEnabled is not disabled       | warn     |
| CLD003 | Claude      | Claude Code prStatusFooterEnabled is not disabled    | warn     |
| CLD004 | Claude      | Claude Code claudeInChromeDefaultEnabled is not disabled | warn     |
| CLD005 | Claude      | Claude Code sandbox.failIfUnavailable is not enabled | high     |
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
