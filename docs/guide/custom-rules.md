# Custom Rules

obacht supports loading custom rules from an external directory using the `--rules-dir` flag. Each rule is defined in a YAML file with metadata and a Rego policy. Input scripts can be defined inline or in separate files.

## Directory Structure

The `--rules-dir` flag expects a directory with the following layout:

```
my-rules/
  policies/            # YAML rule files
    filesystem.yaml
    custom_ssh.yaml
  inputs/              # Optional shell scripts for collecting facts
    filesystem.sh      # Used by filesystem.yaml (matched by name)
    custom_ssh.sh      # Used by custom_ssh.yaml
```

- **`policies/`** — Contains YAML rule files and optional `.rego` policy files
- **`inputs/`** — Contains shell scripts that collect system facts. Scripts are matched to YAML files by name (e.g., `inputs/foo.sh` provides input for `policies/foo.yaml`)

## Rule File Format

A rule file is a YAML file with:

- **`input`** (file-level, optional) — An inline shell script that collects facts and outputs JSON to stdout, shared by all rules in the file. If omitted, obacht looks for a matching script in the `inputs/` directory.
- **`rules`** — A list of rules, each with metadata and its own Rego policy

Each rule contains:

- **`id`**, **`title`**, **`severity`**, **`category`** — Rule metadata
- **`description`**, **`remediation`** — Human-readable details
- **`policy`** — Rego policy for this specific rule (inline or file reference)
- **`input`** (optional) — Rule-specific input script, overrides file-level input

> **Auto-prefix:** The Rego `package obacht.<category>` declaration and `import rego.v1` are automatically prepended based on the rule's `category` field. You only need to write the policy body (e.g. `findings contains f if { ... }`). If your policy already contains a `package` declaration, it is used as-is.

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

### Example: Input from file

Instead of inlining the input script, place it in `inputs/<name>.sh`:

```
my-rules/
  policies/
    ssh_check.yaml     # No inline input: field
  inputs/
    ssh_check.sh       # Automatically used as input
```

The script in `inputs/ssh_check.sh` is used as the file-level input for `policies/ssh_check.yaml`. If the YAML also has an inline `input:` field, the inline value takes precedence.

### Policy file reference

Instead of inline Rego, you can reference a `.rego` file in the `policies/` directory:

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
If no input script is defined (neither inline nor in `inputs/`), the rule is marked as `skip`.

### Tips

- Use `stat -f '%Lp'` (macOS) or `stat -c '%a'` (Linux) for file permissions
- Use `command -v` to check if a tool is installed
- Never output sensitive values (passwords, tokens) — only names and metadata
- Scripts run with a 30-second timeout

## Replacement Behavior

`--rules-dir` replaces the embedded rule set entirely. When the flag is set, only rules
from the given directory run; built-in rules are not loaded. To customise a built-in rule,
copy it into your own rules directory and edit it there.

## Skipping a rule from Rego

When a rule cannot evaluate because runtime state needed for the check is unavailable — a disconnected disk, an unreachable network, a daemon that isn't running — emit a `skips` entry instead of a `findings` entry. The result shows as `skip` rather than a misleading `pass` or `fail`.

```rego
skips contains s if {
  input.os == "darwin"
  input.timemachine_enabled
  not input.timemachine_destination_connected
  s := {
    "rule_id": "OS033",
    "evidence": "Time Machine destination not connected — cannot evaluate backup recency",
  }
}

findings contains f if {
  input.os == "darwin"
  input.timemachine_enabled
  input.timemachine_destination_connected
  not input.timemachine_recent_backup
  f := {"rule_id": "OS033", "evidence": "Time Machine has no backup within the last 14 days"}
}
```

The `skips` collection has the same shape as `findings`: a `rule_id` string and an `evidence` string. The engine looks up both collections by rule ID. If a rule emits both a fail and a skip in the same evaluation, fail wins — the rule is reported as `fail`.

**Skip vs pass vs fail:**

- **`skip`** — the check could not run because required state is unavailable. Tell the user *why*.
- **`pass`** — the check ran and the system is in the desired state.
- **`fail`** — the check ran and the system is in an undesired state.

A "feature off" condition is *not* a skip. If a rule has nothing meaningful to say when a feature is disabled (e.g. Time Machine off, covered by a separate rule), the rule simply passes silently.

## Usage

```bash
obacht scan --rules-dir ./my-rules
```
