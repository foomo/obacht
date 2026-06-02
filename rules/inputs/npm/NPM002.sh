#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# npm install-script execution gating.
#
# With `ignore-scripts=true`, npm refuses to run lifecycle scripts
# (preinstall, install, postinstall, ...) of dependencies. This is the
# single most effective mitigation against compromised npm packages.

if ! command -v npm >/dev/null 2>&1; then
  emit_skip NPM002 "npm not installed"
  exit 0
fi

value=$(npm config get ignore-scripts 2>/dev/null | tr -d '[:space:]')

ignore_scripts=false
case "$value" in
  true|"True"|"TRUE") ignore_scripts=true ;;
esac

jq -cn \
  --argjson ignore_scripts "$ignore_scripts" \
  --arg raw "$value" \
  '{ignore_scripts: $ignore_scripts, raw: $raw}' | emit_ok NPM002
