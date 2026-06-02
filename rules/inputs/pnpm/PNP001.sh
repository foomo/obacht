#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# pnpm package-manager release-age gating.
#
# Reads pnpm's `minimumReleaseAge` setting (minutes), introduced in
# pnpm 10.16 and defaulted to 1440 in pnpm 11.

if ! command -v pnpm >/dev/null 2>&1; then
  emit_skip PNP001 "pnpm not installed"
  exit 0
fi

minutes=$(pnpm config get minimumReleaseAge 2>/dev/null | tr -d '[:space:]')

configured=false
value_seconds=0
case "$minutes" in
  ''|undefined|null|false) ;;
  *[!0-9]*) ;;
  *)
    if [ "$minutes" -gt 0 ]; then
      configured=true
      value_seconds=$((minutes * 60))
    fi
    ;;
esac

jq -cn \
  --argjson configured "$configured" \
  --argjson value_seconds "$value_seconds" \
  '{configured: $configured, value_seconds: $value_seconds}' | emit_ok PNP001
