# Kubernetes Rules

## KUB001: Kubeconfig has weak permissions

**Severity:** high

The kubeconfig file (`~/.kube/config`) contains cluster credentials and access tokens. Weak file permissions can expose these secrets to other users on the system.

**What it checks:**
- File permissions on `~/.kube/config`
- Ensures permissions are `0600` or stricter

**Remediation:**
```bash
chmod 600 ~/.kube/config
```

## KUB002: Production Kubernetes context is active

**Severity:** warn

Having a production cluster set as the active kubectl context increases the risk of accidentally running destructive commands against production infrastructure.

**What it checks:**
- The current kubectl context name
- Flags contexts containing `prod`, `production`, or similar patterns

**Remediation:**
```bash
# Switch to a non-production context
kubectl config use-context dev-cluster
```
