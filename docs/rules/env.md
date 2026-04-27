# Environment Rules

## ENV001: Sensitive credentials found in environment variables

**Severity:** high

Environment variables are visible to all child processes and may be logged or leaked through error reports. Storing secrets such as API keys, tokens, or passwords in environment variables increases the risk of credential exposure.

**What it checks:**

Environment variable names against an exact-match list and a set of suffix patterns. The check is name-only — values are never read or logged.

Suffix patterns (any var ending in one of these is flagged):

`*_PASSWORD`, `*_SECRET`, `*_SECRET_KEY`, `*_API_KEY`, `*_PRIVATE_KEY`, `*_ACCESS_KEY`, `*_LICENSE_KEY`, `*_TOKEN`, `*_AUTH_TOKEN`, `*_ACCESS_TOKEN`, `*_AUTH`, `*_DSN`, `*_CREDENTIAL`, `*_CREDENTIALS`

Exact matches (representative sample — see `rules/inputs/env.sh` for the full list):

- AWS: `AWS_SECRET_ACCESS_KEY`, `AWS_ACCESS_KEY_ID`
- Source control: `GITHUB_TOKEN`, `GITHUB_PAT`, `GITLAB_TOKEN`
- Package registries: `NPM_TOKEN`
- Container/CI: `DOCKER_PASSWORD`
- Comms: `SLACK_TOKEN`, `SLACK_WEBHOOK_URL`
- Connection strings: `DATABASE_URL`, `MYSQL_PASSWORD`, `POSTGRES_PASSWORD`, `REDIS_URL`, `MONGODB_URI`, `MONGO_URL`, `AMQP_URL`, `RABBITMQ_URL`, `CELERY_BROKER_URL`
- Cloud: `GCP_SA_KEY`
- Observability: `DD_APP_KEY`
- Secrets backends: `VAULT_DEV_ROOT_TOKEN_ID`
- SaaS: `SUPABASE_SERVICE_ROLE_KEY`

The suffix patterns auto-cover most cloud, CI/CD, observability, and SaaS tokens (e.g. `OPENAI_API_KEY`, `STRIPE_SECRET_KEY`, `SENTRY_DSN`, `TWILIO_AUTH_TOKEN`, `VAULT_TOKEN`, `VERCEL_TOKEN`, `CLOUDFLARE_API_TOKEN`, `HEROKU_API_KEY`, `MAILGUN_API_KEY`, `SENDGRID_API_KEY`, `NEW_RELIC_LICENSE_KEY`, `HUGGINGFACE_TOKEN`, `BUGSNAG_API_KEY`, `ANSIBLE_VAULT_PASSWORD`, `BITBUCKET_TOKEN`, `GITLAB_PRIVATE_TOKEN`).

**Remediation:**

Use a secrets manager or encrypted configuration file instead of environment variables:

```bash
# Use a .env file with restricted permissions (not checked into git)
chmod 600 .env

# Or use a secrets manager
aws secretsmanager get-secret-value --secret-id my-secret
vault kv get secret/my-app
```
