# Custom Rules

Bouncer supports loading custom rules from an external directory using the `--rules-dir` flag.

## Directory Structure

```
my-rules/
  custom.yaml    # Rule metadata
  custom.rego    # Rego policy
```

## Rule Metadata (YAML)

```yaml
rules:
  - id: CUSTOM001
    title: Custom security check
    severity: warn
    category: custom
    description: |
      Description of what this rule checks.
    remediation: "How to fix the issue"
```

## Rego Policy

```rego
package bouncer.custom

findings[f] {
    # Your policy logic here
    some_condition
    f := {
        "rule_id": "CUSTOM001",
        "evidence": "Description of what was found"
    }
}
```

## Override Behavior

External rules with the same ID as built-in rules will override them. This allows you to customize severity levels or detection logic for existing checks.

## Usage

```bash
bouncer scan --rules-dir ./my-rules
```
