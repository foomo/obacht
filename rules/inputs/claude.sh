#!/bin/sh
# Check if Claude Code CLI is installed.
if ! command -v claude >/dev/null 2>&1; then
  printf '{"installed": false, "gitignore_excludes_settings": false}'
  exit 0
fi

# Check if global gitignore excludes Claude Code local settings.
gitignore_excludes_settings=false
excludes_file=$(git config --global core.excludesfile 2>/dev/null || true)
if [ -n "$excludes_file" ]; then
  # Expand ~ in path.
  excludes_file=$(eval echo "$excludes_file")
  if [ -f "$excludes_file" ] && grep -Fxq '**/.claude/settings.local.json' "$excludes_file" 2>/dev/null; then
    gitignore_excludes_settings=true
  fi
fi

printf '{"installed": true, "gitignore_excludes_settings": %s}' "$gitignore_excludes_settings"
