#!/bin/sh
# Check if git is installed.
if ! command -v git >/dev/null 2>&1; then
  printf '{"installed": false, "version": "", "credential_helper": "", "signing_enabled": false, "signing_format": ""}'
  exit 0
fi

version=$(git --version 2>/dev/null | tr -d '\n')
config=$(git config --global --list 2>/dev/null || true)

credential_helper=""
signing_enabled=false
signing_format=""

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
  esac
done

printf '{"installed": true, "version": "%s", "credential_helper": "%s", "signing_enabled": %s, "signing_format": "%s"}' \
  "$version" "$credential_helper" "$signing_enabled" "$signing_format"
