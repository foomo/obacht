#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

installed=false
signing_enabled=false
if command -v git >/dev/null 2>&1; then
  installed=true
  value=$(git config --global commit.gpgsign 2>/dev/null || true)
  case "$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]')" in
    true) signing_enabled=true ;;
  esac
fi

jq -cn \
  --argjson installed "$installed" \
  --argjson signing_enabled "$signing_enabled" \
  '{installed: $installed, signing_enabled: $signing_enabled}' | emit_ok GIT002
