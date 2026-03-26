# Rules Overview

Bouncer includes 13 built-in rules across 9 categories.

| ID | Category | Title | Severity |
|----|----------|-------|----------|
| SSH001 | SSH | SSH private key has weak permissions | high |
| SSH002 | SSH | SSH directory has weak permissions | high |
| GIT001 | Git | Git credential helper stores passwords in plaintext | high |
| GIT002 | Git | Git commit signing is not enabled | warn |
| DOC001 | Docker | Docker socket has overly permissive access | high |
| DOC002 | Docker | User is in the docker group | warn |
| KUB001 | Kubernetes | Kubeconfig has weak permissions | high |
| KUB002 | Kubernetes | Production Kubernetes context is active | warn |
| ENV001 | Environment | Sensitive credentials found in environment variables | high |
| SHL001 | Shell | Shell history file has weak permissions | warn |
| TOL001 | Tools | Security-relevant tool is missing | info |
| PTH001 | PATH | World-writable directory in PATH | high |
| PTH002 | PATH | Relative path entry in PATH | warn |
| OS001 | OS | Operating system may have pending security updates | info |
