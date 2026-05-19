#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# .npmrc file with auth token permissions.
npmrc="$HOME/.npmrc"
npmrc_has_token=false
npmrc_mode=""
if [ -f "$npmrc" ]; then
  npmrc_mode=$(stat -f '%Lp' "$npmrc" 2>/dev/null || stat -c '%a' "$npmrc" 2>/dev/null || echo "")
  npmrc_mode="0$npmrc_mode"
  if grep -q '_authToken' "$npmrc" 2>/dev/null; then
    npmrc_has_token=true
  fi
fi

printf '{"npmrc_has_token": %s, "npmrc_mode": "%s"}' "$npmrc_has_token" "$npmrc_mode" | emit_ok CRD004
