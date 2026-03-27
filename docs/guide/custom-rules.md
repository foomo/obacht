# Custom Rules

Bouncer supports loading custom rules from an external directory using the `--rules-dir` flag. Each rule is defined in a YAML file with metadata, an input script, and a Rego policy.

## Rule File Format

A rule file is a YAML file with:

- **`input`** (file-level) — A shell script that collects facts and outputs JSON to stdout, shared by all rules in the file
- **`rules`** — A list of rules, each with metadata and its own Rego policy

Each rule contains:

- **`id`**, **`title`**, **`severity`**, **`category`** — Rule metadata
- **`description`**, **`remediation`** — Human-readable details
- **`policy`** — Rego policy for this specific rule (inline or file reference)
- **`input`** (optional) — Rule-specific input script, overrides file-level input

> **Auto-prefix:** The Rego `package bouncer.<category>` declaration and `import rego.v1` are automatically prepended based on the rule's `category` field. You only need to write the policy body (e.g. `findings contains f if { ... }`). If your policy already contains a `package` declaration, it is used as-is.

### Example: Single rule

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
      findings contains f if {
        input.mode != "700"
        f := {
          "rule_id": "DL001",
          "evidence": sprintf("~/Downloads has mode %s (expected 700)", [input.mode]),
        }
      }
```

### Example: Multiple rules with shared input

When rules share the same collected data, define `input` at file-level and `policy` per-rule:

```yaml
input: |
  #!/bin/sh
  dir_mode=$(stat -f '%Lp' ~/.ssh 2>/dev/null || echo "")
  config_exists=false
  [ -f ~/.ssh/config ] && config_exists=true
  printf '{"directory_mode": "0%s", "config_exists": %s}' "$dir_mode" "$config_exists"

rules:
  - id: CSSH001
    title: SSH directory has weak permissions
    severity: high
    category: ssh
    description: The ~/.ssh directory should only be accessible by the owner.
    remediation: "Run: chmod 700 ~/.ssh"
    policy: |
      findings contains f if {
        input.directory_mode != "0700"
        f := {
          "rule_id": "CSSH001",
          "evidence": sprintf("~/.ssh has mode %s", [input.directory_mode]),
        }
      }
  - id: CSSH002
    title: SSH config file missing
    severity: info
    category: ssh
    description: An SSH config file helps manage connections securely.
    remediation: "Create ~/.ssh/config"
    policy: |
      findings contains f if {
        not input.config_exists
        f := {
          "rule_id": "CSSH002",
          "evidence": "~/.ssh/config does not exist",
        }
      }
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
  custom_ssh.yaml    # Shared input with per-rule policy
  custom.rego        # Optional external rego file
```
