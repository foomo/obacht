# Custom Rules

Bouncer supports loading custom rules from an external directory using the `--rules-dir` flag. Each rule is a self-contained YAML file with metadata, an input script, and a Rego policy.

## Rule File Format

A rule file is a YAML file with three sections:

- **`input`** — A shell script that collects facts and outputs JSON to stdout
- **`policy`** — A Rego policy that evaluates the input and produces findings
- **`rules`** — Rule metadata (ID, title, severity, category, description, remediation)

### Self-contained rule (single rule per file)

```yaml
rules:
  - id: DL001
    title: Downloads directory has secure permissions
    severity: warn
    category: filesystem
    description: |
      ~/Downloads should have restricted permissions to prevent
      unauthorized access to downloaded files.
    remediation: "Run: chmod 700 ~/Downloads"
    input: |
      #!/bin/sh
      mode=$(stat -f '%Lp' ~/Downloads 2>/dev/null || stat -c '%a' ~/Downloads 2>/dev/null)
      printf '{"mode": "%s"}' "$mode"
    policy: |
      package bouncer.filesystem

      import rego.v1

      findings contains f if {
        input.mode != "700"
        f := {
          "rule_id": "DL001",
          "evidence": sprintf("~/Downloads has mode %s (expected 700)", [input.mode]),
        }
      }
```

### Shared input/policy (multiple rules per file)

When multiple rules share the same input data, define `input` and `policy` at the file level:

```yaml
input: |
  #!/bin/sh
  dir_mode=$(stat -f '%Lp' ~/.ssh 2>/dev/null || echo "")
  printf '{"directory_mode": "0%s"}' "$dir_mode"

policy: |
  package bouncer.custom_ssh

  import rego.v1

  findings contains f if {
    input.directory_mode != "0700"
    f := {
      "rule_id": "CSSH001",
      "evidence": sprintf("~/.ssh has mode %s", [input.directory_mode]),
    }
  }

rules:
  - id: CSSH001
    title: SSH directory has weak permissions
    severity: high
    category: ssh
    description: The ~/.ssh directory should only be accessible by the owner.
    remediation: "Run: chmod 700 ~/.ssh"
```

### Policy file reference

Instead of inline Rego, you can reference a `.rego` file in the same directory:

```yaml
rules:
  - id: CUSTOM001
    title: Custom check
    severity: warn
    category: custom
    input: |
      printf '{"value": true}'
    policy: custom.rego
```

## Input Scripts

Input scripts are shell commands that:

1. Collect system facts (file permissions, command output, config values, etc.)
2. Output valid JSON to stdout
3. Exit with code 0 on success

If the script fails (non-zero exit), the rule is marked as `error`.
If no input script is defined, the rule is marked as `skip`.

### Tips

- Use `stat -f '%Lp'` (macOS) or `stat -c '%a'` (Linux) for file permissions
- Use `command -v` to check if a tool is installed
- Never output sensitive values (passwords, tokens) — only names and metadata
- Scripts run with a 30-second timeout

## Override Behavior

External rules with the same ID as built-in rules will override them. This allows you to customize severity levels or detection logic for existing checks.

## Usage

```bash
bouncer scan --rules-dir ./my-rules
```

## Directory Structure

```
my-rules/
  filesystem.yaml    # Self-contained rule file
  custom_ssh.yaml    # Shared input/policy rule file
  custom.rego        # Optional external rego file
```
