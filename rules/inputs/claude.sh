#!/bin/sh
# Detect Claude Desktop Native Messaging manifests installed under
# ~/Library/Application Support for Chromium-based browsers. List paths
# that are not yet mitigated. Mitigated = empty file with the user
# immutable (uchg) flag set, so Claude Desktop cannot rewrite it.
claude_desktop_native_messaging_manifests="[]"
manifest_search_dir="$HOME/Library/Application Support"
if [ -d "$manifest_search_dir" ]; then
  manifest_files=$(find "$manifest_search_dir" -type f -name "com.anthropic.claude_browser_extension.json" 2>/dev/null)
  if [ -n "$manifest_files" ]; then
    saved_ifs=$IFS
    IFS='
'
    json="["
    first=1
    for mf in $manifest_files; do
      mf_size=$(stat -f "%z" "$mf" 2>/dev/null || echo 0)
      mf_flags=$(stat -f "%Sf" "$mf" 2>/dev/null || echo "")
      # Skip mitigated manifests: locked with uchg AND too small to
      # contain a functional Native Messaging manifest (tolerates both
      # `: > f` (0 bytes) and `echo "" > f` (1 byte newline) forms).
      if [ "$mf_size" -le 1 ] && [ "$mf_flags" = "uchg" ]; then
        continue
      fi
      mf_escaped=$(printf '%s' "$mf" | sed 's/\\/\\\\/g; s/"/\\"/g')
      if [ "$first" = "1" ]; then
        json="${json}\"${mf_escaped}\""
        first=0
      else
        json="${json},\"${mf_escaped}\""
      fi
    done
    json="${json}]"
    IFS=$saved_ifs
    claude_desktop_native_messaging_manifests=$json
  fi
fi

# Check if Claude Code CLI is installed.
if ! command -v claude >/dev/null 2>&1; then
  printf '{"installed": false, "gitignore_excludes_settings": false, "config_present": false, "auto_compact_enabled": "unset", "pr_status_footer_enabled": "unset", "claude_in_chrome_default_enabled": "unset", "sandbox_fail_if_unavailable": "unset", "settings_present": false, "env_disable_compact": "unset", "env_disable_telemetry": "unset", "env_disable_bug_command": "unset", "env_disable_auto_compact": "unset", "env_disable_login_command": "unset", "env_disable_logout_command": "unset", "env_disable_error_reporting": "unset", "env_disable_upgrade_command": "unset", "env_disable_feedback_command": "unset", "env_disable_extra_usage_command": "unset", "env_claude_code_disable_fast_mode": "unset", "env_disable_install_github_app_command": "unset", "env_claude_code_disable_cron": "unset", "env_claude_code_disable_feedback_survey": "unset", "env_claude_code_disable_file_checkpointing": "unset", "env_claude_code_disable_experimental_betas": "unset", "env_force_autoupdate_plugins": "unset", "env_is_demo": "unset", "settings_disable_auto_mode": "unset", "settings_disable_deep_link_registration": "unset", "settings_auto_memory_directory": "unset", "settings_plans_directory": "unset", "settings_respect_gitignore": "unset", "settings_skip_web_fetch_preflight": "unset", "settings_attribution_commit": "unset", "settings_attribution_pr": "unset", "sandbox_enabled": "unset", "sandbox_auto_allow_bash_if_sandboxed": "unset", "sandbox_allow_unsandboxed_commands": "unset", "sandbox_network_allow_managed_domains_only": "unset", "sandbox_network_allowed_domains_has_github": "false", "sandbox_network_denied_domains_has_uploads_github": "false", "sandbox_filesystem_allow_write_has_npm_logs": "false", "sandbox_filesystem_allow_write_has_claude_debug": "false", "permissions_present": false, "permissions_disable_bypass_mode": "unset", "permissions_deny_network_missing": "", "permissions_deny_destructive_fs_missing": "", "permissions_deny_git_missing": "", "permissions_deny_home_secrets_missing": "", "permissions_deny_project_secrets_missing": "", "claude_desktop_native_messaging_manifests": %s}' "$claude_desktop_native_messaging_manifests"
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

# Locate Claude Code user settings.json (separate from .claude.json above).
settings_dir="${CLAUDE_CONFIG_DIR:-$HOME/.claude}"
settings_file="$settings_dir/settings.json"

settings_present=false
env_disable_compact="unset"
env_disable_telemetry="unset"
env_disable_bug_command="unset"
env_disable_auto_compact="unset"
env_disable_login_command="unset"
env_disable_logout_command="unset"
env_disable_error_reporting="unset"
env_disable_upgrade_command="unset"
env_disable_feedback_command="unset"
env_disable_extra_usage_command="unset"
env_claude_code_disable_fast_mode="unset"
env_disable_install_github_app_command="unset"
env_claude_code_disable_cron="unset"
env_claude_code_disable_feedback_survey="unset"
env_claude_code_disable_file_checkpointing="unset"
env_claude_code_disable_experimental_betas="unset"
env_force_autoupdate_plugins="unset"
env_is_demo="unset"

settings_disable_auto_mode="unset"
settings_disable_deep_link_registration="unset"
settings_auto_memory_directory="unset"
settings_plans_directory="unset"
settings_respect_gitignore="unset"
settings_skip_web_fetch_preflight="unset"
settings_attribution_commit="unset"
settings_attribution_pr="unset"

sandbox_enabled="unset"
sandbox_auto_allow_bash_if_sandboxed="unset"
sandbox_allow_unsandboxed_commands="unset"
sandbox_network_allow_managed_domains_only="unset"
sandbox_network_allowed_domains_has_github="false"
sandbox_network_denied_domains_has_uploads_github="false"
sandbox_filesystem_allow_write_has_npm_logs="false"
sandbox_filesystem_allow_write_has_claude_debug="false"

permissions_present=false
permissions_disable_bypass_mode="unset"
permissions_deny_network_missing=""
permissions_deny_destructive_fs_missing=""
permissions_deny_git_missing=""
permissions_deny_home_secrets_missing=""
permissions_deny_project_secrets_missing=""

# Extract a string-typed key from the file.
# Matches "KEY": "VALUE" anywhere. Returns "<missing>" if not found, or the
# value (which may be empty when the file contains "KEY": "").
extract_string() {
  key=$1
  file=$2
  if grep -Eq "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$file" 2>/dev/null; then
    grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$file" 2>/dev/null | head -1 | sed 's/.*:[[:space:]]*"\(.*\)"[[:space:]]*$/\1/'
  else
    printf 'unset'
  fi
}

# Extract a string-typed key scoped to a content string (not a file).
extract_string_in_block() {
  key=$1
  content=$2
  if printf '%s' "$content" | grep -Eq "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\""; then
    printf '%s' "$content" | grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" | head -1 | sed 's/.*:[[:space:]]*"\(.*\)"[[:space:]]*$/\1/'
  else
    printf 'unset'
  fi
}

# Extract the {...} body of a top-level object setting (e.g. attribution).
# Returns the body or empty string if not found.
extract_object_block() {
  key=$1
  file=$2
  tr '\n' ' ' < "$file" 2>/dev/null \
    | grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\\{[^}]*\\}" \
    | head -1
}

# Extract the {...} body of an object key from a content string, handling
# arbitrarily nested braces via balanced-brace counting.
# Stdin: content; arg1: key. Outputs the matched object body or empty.
extract_nested_block_in() {
  key=$1
  awk -v key="$key" '
    function extract(content, key,    depth, c, start, i, pat, found_re, m) {
      pat = "\""key"\"[[:space:]]*:[[:space:]]*\\{"
      if (match(content, pat)) {
        start = RSTART + RLENGTH - 1
        depth = 1
        for (i = start + 1; i <= length(content); i++) {
          c = substr(content, i, 1)
          if (c == "{") depth++
          else if (c == "}") {
            depth--
            if (depth == 0) return substr(content, start, i - start + 1)
          }
        }
      }
      return ""
    }
    { data = data $0 "\n" }
    END { printf "%s", extract(data, key) }
  '
}

# Extract a top-level boolean from a content string (true|false|unset).
extract_bool_in() {
  key=$1
  content=$2
  value=$(printf '%s' "$content" | grep -o "\"$key\"[[:space:]]*:[[:space:]]*\(true\|false\)" 2>/dev/null | head -1 | sed 's/.*:[[:space:]]*//')
  if [ -z "$value" ]; then
    printf 'unset'
  else
    printf '%s' "$value"
  fi
}

# Returns "true" if the string array stored under <key> in <content> contains
# the literal string <needle>. Returns "false" otherwise (or if the array is
# missing). Handles only flat string arrays (no nested objects/arrays).
array_contains_string() {
  key=$1
  needle=$2
  content=$3
  arr=$(printf '%s' "$content" | tr '\n' ' ' | grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\\[[^]]*\\]" | head -1)
  if [ -z "$arr" ]; then
    printf 'false'
  elif printf '%s' "$arr" | grep -Fq "\"$needle\""; then
    printf 'true'
  else
    printf 'false'
  fi
}

# Returns space-separated list of needles missing from the string array stored
# under <key> in <content>. Empty output means all present (or array missing,
# in which case all needles are reported missing). Handles only flat string
# arrays (no nested objects/arrays).
array_missing_strings() {
  key=$1
  content=$2
  shift 2
  arr=$(printf '%s' "$content" | tr '\n' ' ' | grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\\[[^]]*\\]" | head -1)
  missing=""
  for needle in "$@"; do
    if [ -z "$arr" ] || ! printf '%s' "$arr" | grep -Fq "\"$needle\""; then
      if [ -z "$missing" ]; then
        missing="$needle"
      else
        missing="$missing $needle"
      fi
    fi
  done
  printf '%s' "$missing"
}

# extract_env keeps the prior behaviour for env-block keys: it does not
# distinguish unset from empty-value, since env values are always non-empty.
extract_env() {
  key=$1
  file=$2
  value=$(grep -Eo "\"$key\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$file" 2>/dev/null | head -1 | sed 's/.*:[[:space:]]*"\(.*\)"[[:space:]]*$/\1/')
  if [ -z "$value" ]; then
    printf 'unset'
  else
    printf '%s' "$value"
  fi
}

if [ -f "$settings_file" ]; then
  settings_present=true
  env_disable_compact=$(extract_env "DISABLE_COMPACT" "$settings_file")
  env_disable_telemetry=$(extract_env "DISABLE_TELEMETRY" "$settings_file")
  env_disable_bug_command=$(extract_env "DISABLE_BUG_COMMAND" "$settings_file")
  env_disable_auto_compact=$(extract_env "DISABLE_AUTO_COMPACT" "$settings_file")
  env_disable_login_command=$(extract_env "DISABLE_LOGIN_COMMAND" "$settings_file")
  env_disable_logout_command=$(extract_env "DISABLE_LOGOUT_COMMAND" "$settings_file")
  env_disable_error_reporting=$(extract_env "DISABLE_ERROR_REPORTING" "$settings_file")
  env_disable_upgrade_command=$(extract_env "DISABLE_UPGRADE_COMMAND" "$settings_file")
  env_disable_feedback_command=$(extract_env "DISABLE_FEEDBACK_COMMAND" "$settings_file")
  env_disable_extra_usage_command=$(extract_env "DISABLE_EXTRA_USAGE_COMMAND" "$settings_file")
  env_claude_code_disable_fast_mode=$(extract_env "CLAUDE_CODE_DISABLE_FAST_MODE" "$settings_file")
  env_disable_install_github_app_command=$(extract_env "DISABLE_INSTALL_GITHUB_APP_COMMAND" "$settings_file")
  env_claude_code_disable_cron=$(extract_env "CLAUDE_CODE_DISABLE_CRON" "$settings_file")
  env_claude_code_disable_feedback_survey=$(extract_env "CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY" "$settings_file")
  env_claude_code_disable_file_checkpointing=$(extract_env "CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING" "$settings_file")
  env_claude_code_disable_experimental_betas=$(extract_env "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS" "$settings_file")
  env_force_autoupdate_plugins=$(extract_env "FORCE_AUTOUPDATE_PLUGINS" "$settings_file")
  env_is_demo=$(extract_env "IS_DEMO" "$settings_file")

  settings_disable_auto_mode=$(extract_string "disableAutoMode" "$settings_file")
  settings_disable_deep_link_registration=$(extract_string "disableDeepLinkRegistration" "$settings_file")
  settings_auto_memory_directory=$(extract_string "autoMemoryDirectory" "$settings_file")
  settings_plans_directory=$(extract_string "plansDirectory" "$settings_file")
  settings_respect_gitignore=$(extract_bool "respectGitignore" "$settings_file")
  settings_skip_web_fetch_preflight=$(extract_bool "skipWebFetchPreflight" "$settings_file")

  attribution_block=$(extract_object_block "attribution" "$settings_file")
  if [ -n "$attribution_block" ]; then
    settings_attribution_commit=$(extract_string_in_block "commit" "$attribution_block")
    settings_attribution_pr=$(extract_string_in_block "pr" "$attribution_block")
  fi

  sandbox_block=$(extract_nested_block_in "sandbox" < "$settings_file")
  if [ -n "$sandbox_block" ]; then
    sandbox_enabled=$(extract_bool_in "enabled" "$sandbox_block")
    sandbox_auto_allow_bash_if_sandboxed=$(extract_bool_in "autoAllowBashIfSandboxed" "$sandbox_block")
    sandbox_allow_unsandboxed_commands=$(extract_bool_in "allowUnsandboxedCommands" "$sandbox_block")

    network_block=$(printf '%s' "$sandbox_block" | extract_nested_block_in "network")
    if [ -n "$network_block" ]; then
      sandbox_network_allow_managed_domains_only=$(extract_bool_in "allowManagedDomainsOnly" "$network_block")
      sandbox_network_allowed_domains_has_github=$(array_contains_string "allowedDomains" "github.com" "$network_block")
      sandbox_network_denied_domains_has_uploads_github=$(array_contains_string "deniedDomains" "uploads.github.com" "$network_block")
    fi

    filesystem_block=$(printf '%s' "$sandbox_block" | extract_nested_block_in "filesystem")
    if [ -n "$filesystem_block" ]; then
      sandbox_filesystem_allow_write_has_npm_logs=$(array_contains_string "allowWrite" "~/.cache/npm/logs" "$filesystem_block")
      sandbox_filesystem_allow_write_has_claude_debug=$(array_contains_string "allowWrite" "~/.config/claude/debug" "$filesystem_block")
    fi
  fi

  permissions_block=$(extract_nested_block_in "permissions" < "$settings_file")
  if [ -n "$permissions_block" ]; then
    permissions_present=true
    permissions_disable_bypass_mode=$(extract_string_in_block "disableBypassPermissionsMode" "$permissions_block")
    permissions_deny_network_missing=$(array_missing_strings "deny" "$permissions_block" \
      "Bash(nc:*)" "Bash(netcat:*)" "Bash(socat:*)" \
      "Bash(ssh:*)" "Bash(scp:*)" "Bash(rsync:*)")
    permissions_deny_destructive_fs_missing=$(array_missing_strings "deny" "$permissions_block" \
      "Bash(chmod 777:*)" "Bash(chown:*)" \
      "Bash(rm -rf /:*)" "Bash(rm -rf ~:*)" \
      "Bash(dd:*)" "Bash(mkfs:*)")
    permissions_deny_git_missing=$(array_missing_strings "deny" "$permissions_block" \
      "Bash(git push:*)" "Bash(git tag:*)" "Bash(git reset --hard:*)")
    permissions_deny_home_secrets_missing=$(array_missing_strings "deny" "$permissions_block" \
      "Read(~/.ssh/**)" "Read(~/.aws/**)" "Read(~/.gnupg/**)" \
      "Read(~/.config/gh/**)" "Read(~/.kube/**)" "Read(~/.docker/config.json)")
    permissions_deny_project_secrets_missing=$(array_missing_strings "deny" "$permissions_block" \
      "Read(./.env)" "Read(./.env.*)" "Read(./*.pem)" "Read(./*.key)" \
      "Read(./**/.env)" "Read(./**/.env.*)" "Read(./**/*.pem)" "Read(./**/*.key)" \
      "Read(./**/id_rsa*)" "Read(./**/id_ed25519*)" "Read(./**/credentials*)")
  fi
fi

printf '{"installed": true, "gitignore_excludes_settings": %s, "config_present": %s, "auto_compact_enabled": "%s", "pr_status_footer_enabled": "%s", "claude_in_chrome_default_enabled": "%s", "sandbox_fail_if_unavailable": "%s", "settings_present": %s, "env_disable_compact": "%s", "env_disable_telemetry": "%s", "env_disable_bug_command": "%s", "env_disable_auto_compact": "%s", "env_disable_login_command": "%s", "env_disable_logout_command": "%s", "env_disable_error_reporting": "%s", "env_disable_upgrade_command": "%s", "env_disable_feedback_command": "%s", "env_disable_extra_usage_command": "%s", "env_claude_code_disable_fast_mode": "%s", "env_disable_install_github_app_command": "%s", "env_claude_code_disable_cron": "%s", "env_claude_code_disable_feedback_survey": "%s", "env_claude_code_disable_file_checkpointing": "%s", "env_claude_code_disable_experimental_betas": "%s", "env_force_autoupdate_plugins": "%s", "env_is_demo": "%s", "settings_disable_auto_mode": "%s", "settings_disable_deep_link_registration": "%s", "settings_auto_memory_directory": "%s", "settings_plans_directory": "%s", "settings_respect_gitignore": "%s", "settings_skip_web_fetch_preflight": "%s", "settings_attribution_commit": "%s", "settings_attribution_pr": "%s", "sandbox_enabled": "%s", "sandbox_auto_allow_bash_if_sandboxed": "%s", "sandbox_allow_unsandboxed_commands": "%s", "sandbox_network_allow_managed_domains_only": "%s", "sandbox_network_allowed_domains_has_github": "%s", "sandbox_network_denied_domains_has_uploads_github": "%s", "sandbox_filesystem_allow_write_has_npm_logs": "%s", "sandbox_filesystem_allow_write_has_claude_debug": "%s", "permissions_present": %s, "permissions_disable_bypass_mode": "%s", "permissions_deny_network_missing": "%s", "permissions_deny_destructive_fs_missing": "%s", "permissions_deny_git_missing": "%s", "permissions_deny_home_secrets_missing": "%s", "permissions_deny_project_secrets_missing": "%s", "claude_desktop_native_messaging_manifests": %s}' \
  "$gitignore_excludes_settings" "$config_present" \
  "$auto_compact_enabled" "$pr_status_footer_enabled" \
  "$claude_in_chrome_default_enabled" "$sandbox_fail_if_unavailable" \
  "$settings_present" \
  "$env_disable_compact" "$env_disable_telemetry" "$env_disable_bug_command" \
  "$env_disable_auto_compact" "$env_disable_login_command" "$env_disable_logout_command" \
  "$env_disable_error_reporting" "$env_disable_upgrade_command" "$env_disable_feedback_command" \
  "$env_disable_extra_usage_command" "$env_claude_code_disable_fast_mode" "$env_disable_install_github_app_command" \
  "$env_claude_code_disable_cron" "$env_claude_code_disable_feedback_survey" "$env_claude_code_disable_file_checkpointing" \
  "$env_claude_code_disable_experimental_betas" "$env_force_autoupdate_plugins" "$env_is_demo" \
  "$settings_disable_auto_mode" "$settings_disable_deep_link_registration" \
  "$settings_auto_memory_directory" "$settings_plans_directory" \
  "$settings_respect_gitignore" "$settings_skip_web_fetch_preflight" \
  "$settings_attribution_commit" "$settings_attribution_pr" \
  "$sandbox_enabled" "$sandbox_auto_allow_bash_if_sandboxed" "$sandbox_allow_unsandboxed_commands" \
  "$sandbox_network_allow_managed_domains_only" "$sandbox_network_allowed_domains_has_github" "$sandbox_network_denied_domains_has_uploads_github" \
  "$sandbox_filesystem_allow_write_has_npm_logs" "$sandbox_filesystem_allow_write_has_claude_debug" \
  "$permissions_present" "$permissions_disable_bypass_mode" \
  "$permissions_deny_network_missing" "$permissions_deny_destructive_fs_missing" \
  "$permissions_deny_git_missing" "$permissions_deny_home_secrets_missing" \
  "$permissions_deny_project_secrets_missing" \
  "$claude_desktop_native_messaging_manifests"
