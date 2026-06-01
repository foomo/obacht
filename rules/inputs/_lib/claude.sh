#!/bin/sh
# shellcheck shell=sh
#
# Helpers for Claude Code (claude CLI) configuration probes.
#
# Resolution order for both config files:
#   1. $XDG_CONFIG_HOME/claude/{.claude.json,settings.json}
#   2. $HOME/.claude.json + $HOME/.claude/settings.json
#
# Usage:
#   # include: _lib/json.sh
#   # include: _lib/claude.sh
#   claude_skip_if_missing CLD001
#   v=$(claude_config_bool '.autoCompactEnabled')
#   jq -cn --arg v "$v" '{auto_compact_enabled: $v}' | emit_ok CLD001

# Echo absolute path to the active .claude.json, or empty if none found.
claude_config_path() {
  if [ -n "${XDG_CONFIG_HOME:-}" ] && [ -f "$XDG_CONFIG_HOME/claude/.claude.json" ]; then
    printf '%s' "$XDG_CONFIG_HOME/claude/.claude.json"
    return
  fi
  if [ -f "$HOME/.claude.json" ]; then
    printf '%s' "$HOME/.claude.json"
  fi
}

# Echo absolute path to the active settings.json, or empty if none found.
claude_settings_path() {
  if [ -n "${XDG_CONFIG_HOME:-}" ] && [ -f "$XDG_CONFIG_HOME/claude/settings.json" ]; then
    printf '%s' "$XDG_CONFIG_HOME/claude/settings.json"
    return
  fi
  if [ -f "$HOME/.claude/settings.json" ]; then
    printf '%s' "$HOME/.claude/settings.json"
  fi
}

# Emit a skip envelope and exit 0 when the `claude` CLI is not on PATH.
# Every CLD rule must call this first so the entire category is a no-op
# on hosts without Claude Code installed.
claude_skip_if_missing() {
  if ! command -v claude >/dev/null 2>&1; then
    emit_skip "$1" "claude CLI not installed"
    exit 0
  fi
}

# Internal: read a jq path from a JSON file. Echoes the value as a string
# ("true"/"false"/"<value>") or "unset" when the file is missing OR the
# path resolves to null (i.e. the key is absent). An *explicit* empty
# string is preserved as "" (needed by CLD026).
_claude_read_jq() {
  file="$1"
  expr="$2"
  if [ -z "$file" ]; then
    printf 'unset'
    return
  fi
  jq -r "$expr as \$v | if \$v == null then \"unset\" else (\$v | tostring) end" "$file" 2>/dev/null || printf 'unset'
}

# Read a value from .claude.json by jq path. Echoes the value or "unset".
claude_config_value() {
  _claude_read_jq "$(claude_config_path)" "$1"
}

# Read a value from settings.json by jq path. Echoes the value or "unset".
claude_settings_value() {
  _claude_read_jq "$(claude_settings_path)" "$1"
}
