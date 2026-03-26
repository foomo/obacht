# Environment Rules

## ENV001: Sensitive credentials found in environment variables

**Severity:** high

Environment variables are visible to all child processes and may be logged or leaked through error reports. Storing secrets such as API keys, tokens, or passwords in environment variables increases the risk of credential exposure.

**What it checks:**
- Environment variable names matching sensitive patterns (e.g., `*_TOKEN`, `*_SECRET`, `*_PASSWORD`, `*_API_KEY`)
- Whether these variables contain non-empty values

**Remediation:**

Use a secrets manager or encrypted configuration file instead of environment variables:

```bash
# Use a .env file with restricted permissions (not checked into git)
chmod 600 .env

# Or use a secrets manager
aws secretsmanager get-secret-value --secret-id my-secret
vault kv get secret/my-app
```
