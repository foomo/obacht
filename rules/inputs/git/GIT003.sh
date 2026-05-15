#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

installed=false
safe_directory_wildcard=false
if command -v git >/dev/null 2>&1; then
  installed=true
  # safe.directory is multi-valued; check all values for '*'.
  values=$(git config --global --get-all safe.directory 2>/dev/null || true)
  IFS='
'
  for v in $values; do
    if [ "$v" = "*" ]; then
      safe_directory_wildcard=true
      break
    fi
  done
  unset IFS
fi

jq -cn \
  --argjson installed "$installed" \
  --argjson safe_directory_wildcard "$safe_directory_wildcard" \
  '{installed: $installed, safe_directory_wildcard: $safe_directory_wildcard}' | emit_ok GIT003
