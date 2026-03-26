# PRD: Environment Security Checker CLI

## Product summary
Build a cross-platform CLI that inspects a coworker’s local development environment and recommends security best practices, misconfigurations, and risky defaults. The tool should be easy to run, produce clear actionable findings, and support policy-driven rules so checks can evolve without rewriting core logic. Go is a strong implementation choice, and OPA/Rego is a good fit for the policy layer. [openpolicyagent](https://openpolicyagent.org/docs/cli)

## Problem statement
Engineers often work with inconsistent local setups, and subtle issues like insecure file permissions, exposed credentials, risky Docker access, or weak Git/SSH settings go unnoticed. The team needs a lightweight, trusted tool that gives quick security feedback without requiring a heavyweight agent or full endpoint management platform. [wiz](https://www.wiz.io/academy/application-security/ci-cd-security-best-practices)

## Goals
- Identify common security issues in developer environments.
- Recommend concrete remediation steps, not just warnings.
- Run locally on macOS, Linux, and Windows with minimal setup.
- Support policy-as-code so security rules can be reviewed and versioned.
- Export results in human-readable and machine-readable formats. [docs.styra](https://docs.styra.com/enterprise-opa/reference/cli-reference)

## Non-goals
- It is not a vulnerability scanner for production servers or cloud accounts.
- It is not an EDR/MDM replacement.
- It should not collect secrets or upload sensitive local data by default.
- It should not enforce changes automatically in the first version. [help.hcl-software](https://help.hcl-software.com/appscan/Enterprise/10.9.1/topics/c_best_practices_production_scan.html)

## Target users
- Software engineers.
- DevOps/platform engineers.
- Security champions embedded in engineering teams.
- New hires who need an environment baseline check. [docs.aws.amazon](https://docs.aws.amazon.com/prescriptive-guidance/latest/internal-developer-platform/principles.html)

## Core use cases
1. A developer runs `envcheck scan` after onboarding and gets a prioritized list of issues.
2. A security lead distributes a policy pack with company-specific checks.
3. A team uses JSON output in CI to flag noncompliant developer images or bootstrap scripts.
4. A user runs `envcheck explain <rule>` to understand why a check matters. [openpolicyagent](https://openpolicyagent.org/docs/integration)

## Functional requirements
- Discover environment facts such as OS, shell, tool versions, config paths, permissions, and known risky settings.
- Evaluate those facts against policy rules.
- Classify findings by severity: info, warn, high, critical.
- Show evidence for each finding and a suggested fix.
- Support output formats: pretty, JSON, and optionally SARIF/CSV.
- Allow policy packs to be bundled locally.
- Allow per-user config for exclusions and custom checks.
- Provide exit codes suitable for automation.
- Include a `doctor` or `baseline` mode for quick health checks.
- Include a `fix` or `recommend` mode that prints remediation guidance only. [docs.styra](https://docs.styra.com/enterprise-opa/reference/cli-reference)

## Suggested checks
- World-writable or weakly permissioned SSH files.
- Git global config risks, credential helpers, and signing settings.
- Docker socket access and membership in privileged groups.
- Kubeconfig file permissions and current-context safety.
- Presence of exposed secrets in environment variables.
- Unsafe shell history behavior.
- Outdated or missing security-relevant tool versions.
- Insecure defaults in cloud CLI configs.
- Suspicious PATH entries or writable directories early in PATH.
- Basic OS patch freshness checks where available. [echo](https://www.echo.ai/blog/container-scanning-best-practices)

## Policy engine
Use OPA/Rego as the policy layer if you want rules to be portable and testable. The CLI should collect normalized facts into a stable JSON schema, then feed that schema to Rego policies that return structured findings. That keeps the scanner code small and makes policy changes reviewable like code. [pkg.go](https://pkg.go.dev/github.com/open-policy-agent/opa)

## Architecture
- `collector`: gathers host facts.
- `schema`: defines the JSON model for facts and findings.
- `engine`: loads and evaluates Rego policies.
- `reporter`: renders terminal, JSON, and file outputs.
- `packs`: versioned policy bundles by environment or team.
- `cli`: command parsing and flags.

This split lets you swap or extend collectors without touching policy logic, and it gives you room to add integrations later. [openpolicyagent](https://openpolicyagent.org/docs/integration)

## UX requirements
- Default command should be `scan` with sensible defaults.
- Output should be concise and prioritize actionable items.
- Findings should include: rule ID, severity, evidence, why it matters, and next step.
- Support `--format json` for automation and `--format pretty` for humans.
- Support `--profile` or `--pack` to select different policy sets.
- Provide `--explain` for rule details and examples. [openpolicyagent](https://openpolicyagent.org/docs/cli)

## Security requirements
- Never exfiltrate secrets by default.
- Minimize collected data and redact sensitive values.
- Bundle policies with signed releases or checksums.
- Avoid remote policy fetches unless explicitly enabled.
- Version the facts schema and policies together.
- Add tests for policies and sample inputs. [cncf](https://www.cncf.io/blog/2025/03/18/open-policy-agent-best-practices-for-a-secure-deployment/)

## Success metrics
- 80%+ of pilot users run the tool more than once.
- Median scan time under 5 seconds on a developer laptop.
- At least 90% of findings include an accepted remediation.
- Low false-positive rate reported by pilot users.
- Clear adoption into onboarding or security checklists. [learn.microsoft](https://learn.microsoft.com/en-us/nuget/concepts/security-best-practices)

## Milestones
### Phase 1
Build the MVP scanner, JSON schema, and 10–15 Rego checks.

### Phase 2
Add policy packs, tests, prettier output, and command-line ergonomics.

### Phase 3
Add CI integration, SARIF export, and team-specific baselines.

### Phase 4
Add guided remediation and optional enforcement modes. [docs.styra](https://docs.styra.com/enterprise-opa/reference/cli-reference)

## Risks
- Too many noisy findings will reduce trust.
- Over-collecting data may create privacy concerns.
- Policy complexity can outgrow the scanner if the schema is unstable.
- Remote policy distribution can become a supply-chain concern if not controlled carefully. [docs.styra](https://docs.styra.com/opa/errors/rego-type-error/unsafe-built-in-function-calls-in-expression-name)

## Open questions
- Should the first release be opinionated for your team’s stack or broadly generic?
- Do you want a purely local tool, or one that can also scan CI images and dev containers?
- Should policies be centrally managed or user-extensible?
- Do you want the CLI to auto-fix any issues, or only recommend fixes? [docs.aws.amazon](https://docs.aws.amazon.com/prescriptive-guidance/latest/internal-developer-platform/principles.html)

I can turn this into a one-page PRD template, or into a more engineering-ready spec with commands, data schema, and milestone estimates.
