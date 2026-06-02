#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# npm git-dependency gating.
#
# npm 11.10+ supports the `allow-git` setting (values: all|none|root).
# Default is `all`, which lets transitive git URLs ship their own
# .npmrc overriding the git executable — a known supply-chain RCE
# vector that bypasses `ignore-scripts`. `none` blocks all git deps;
# `root` allows only those declared in the project's package.json.
# Both are acceptable hardening levels; `all` (or unset) is not.

if ! command -v npm >/dev/null 2>&1; then
  emit_skip NPM003 "npm not installed"
  exit 0
fi

value=$(npm config get allow-git 2>/dev/null | tr -d '[:space:]')

restricted=false
case "$value" in
  none|None|NONE|root|Root|ROOT) restricted=true ;;
esac

jq -cn \
  --argjson restricted "$restricted" \
  --arg raw "$value" \
  '{restricted: $restricted, raw: $raw}' | emit_ok NPM003
