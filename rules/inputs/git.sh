#!/bin/sh
# Check if git is installed.
if ! command -v git >/dev/null 2>&1; then
  printf '{"installed": false, "version": "", "credential_helper": "", "signing_enabled": false, "signing_format": "", "safe_directory_wildcard": false, "gitignore_excludes_env": false}'
  exit 0
fi

version=$(git --version 2>/dev/null | tr -d '\n')
config=$(git config --global --list 2>/dev/null || true)

credential_helper=""
signing_enabled=false
signing_format=""
safe_directory_wildcard=false
excludes_file=""

IFS='
'
for line in $config; do
  key=$(printf '%s' "$line" | cut -d= -f1 | tr '[:upper:]' '[:lower:]')
  value=$(printf '%s' "$line" | cut -d= -f2-)
  case "$key" in
    credential.helper) credential_helper="$value" ;;
    commit.gpgsign)
      case "$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]')" in
        true) signing_enabled=true ;;
      esac ;;
    gpg.format) signing_format="$value" ;;
    safe.directory)
      if [ "$value" = "*" ]; then safe_directory_wildcard=true; fi ;;
    core.excludesfile) excludes_file="$value" ;;
  esac
done

# Check if global gitignore excludes .env files.
gitignore_excludes_env=false
if [ -n "$excludes_file" ]; then
  # Expand ~ in path.
  excludes_file=$(eval echo "$excludes_file")
  if [ -f "$excludes_file" ] && grep -q '\.env' "$excludes_file" 2>/dev/null; then
    gitignore_excludes_env=true
  fi
fi

printf '{"installed": true, "version": "%s", "credential_helper": "%s", "signing_enabled": %s, "signing_format": "%s", "safe_directory_wildcard": %s, "gitignore_excludes_env": %s}' \
  "$version" "$credential_helper" "$signing_enabled" "$signing_format" "$safe_directory_wildcard" "$gitignore_excludes_env"
