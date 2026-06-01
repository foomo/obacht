#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

# Global gitignore excludes Claude Code local settings?
claude_skip_if_missing CLD001

excluded=false
excludes_file=$(git config --global core.excludesfile 2>/dev/null || true)
if [ -n "$excludes_file" ]; then
  # Expand ~ if present.
  excludes_file=$(eval echo "$excludes_file")
  if [ -f "$excludes_file" ] && grep -Fxq '**/.claude/settings.local.json' "$excludes_file" 2>/dev/null; then
    excluded=true
  fi
fi

jq -cn --argjson v "$excluded" '{gitignore_excludes_settings: $v}' | emit_ok CLD001
