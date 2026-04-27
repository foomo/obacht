#!/bin/sh
# Check if Claude Code CLI is installed.
if ! command -v claude >/dev/null 2>&1; then
  printf '{"installed": false, "gitignore_excludes_settings": false, "config_present": false, "auto_compact_enabled": "unset", "pr_status_footer_enabled": "unset", "claude_in_chrome_default_enabled": "unset", "sandbox_fail_if_unavailable": "unset"}'
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

# Locate Claude config file.
config_dir="${CLAUDE_CONFIG_DIR:-$HOME}"
config_file="$config_dir/.claude.json"

config_present=false
auto_compact_enabled="unset"
pr_status_footer_enabled="unset"
claude_in_chrome_default_enabled="unset"
sandbox_fail_if_unavailable="unset"

# Extract a top-level boolean value for the given key (true|false|unset).
extract_bool() {
  key=$1
  file=$2
  value=$(grep -o "\"$key\"[[:space:]]*:[[:space:]]*\(true\|false\)" "$file" 2>/dev/null | head -1 | sed 's/.*:[[:space:]]*//')
  if [ -z "$value" ]; then
    printf 'unset'
  else
    printf '%s' "$value"
  fi
}

if [ -f "$config_file" ]; then
  config_present=true
  auto_compact_enabled=$(extract_bool "autoCompactEnabled" "$config_file")
  pr_status_footer_enabled=$(extract_bool "prStatusFooterEnabled" "$config_file")
  claude_in_chrome_default_enabled=$(extract_bool "claudeInChromeDefaultEnabled" "$config_file")
  sandbox_fail_if_unavailable=$(extract_bool "failIfUnavailable" "$config_file")
fi

printf '{"installed": true, "gitignore_excludes_settings": %s, "config_present": %s, "auto_compact_enabled": "%s", "pr_status_footer_enabled": "%s", "claude_in_chrome_default_enabled": "%s", "sandbox_fail_if_unavailable": "%s"}' \
  "$gitignore_excludes_settings" "$config_present" \
  "$auto_compact_enabled" "$pr_status_footer_enabled" \
  "$claude_in_chrome_default_enabled" "$sandbox_fail_if_unavailable"
