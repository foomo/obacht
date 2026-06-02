#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# npm package-manager release-age gating.
#
# Reads npm's `min-release-age` setting (days), introduced in npm 11.10.0.
# A non-zero, >=7d configured value passes; missing/zero/<7d fails.

if ! command -v npm >/dev/null 2>&1; then
  emit_skip NPM001 "npm not installed"
  exit 0
fi

# Prefer the canonical key. Fall back to the alternate name some early
# npm 11.x prereleases used, then to undefined.
read_npm_config() {
  out=$(npm config get "$1" 2>/dev/null | tr -d '[:space:]')
  case "$out" in
    ''|undefined|null) return 1 ;;
  esac
  printf '%s' "$out"
}

days=""
if v=$(read_npm_config min-release-age); then
  days="$v"
elif v=$(read_npm_config minimum-release-age); then
  days="$v"
fi

configured=false
value_seconds=0
case "$days" in
  ''|*[!0-9]*) ;;
  *)
    if [ "$days" -gt 0 ]; then
      configured=true
      value_seconds=$((days * 86400))
    fi
    ;;
esac

jq -cn \
  --argjson configured "$configured" \
  --argjson value_seconds "$value_seconds" \
  '{configured: $configured, value_seconds: $value_seconds}' | emit_ok NPM001
