#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# pnpm exotic-subdependency gating.
#
# pnpm 10.26+ supports `blockExoticSubdeps`. When true, transitive
# dependencies cannot pull code from git repositories or raw tarball
# URLs — sources that bypass registry security scanning. Direct
# dependencies declared in package.json can still use exotic
# specifiers; only their further subdeps are blocked.
# Defaults to true in pnpm v11.

if ! command -v pnpm >/dev/null 2>&1; then
  emit_skip PNP002 "pnpm not installed"
  exit 0
fi

value=$(pnpm config get blockExoticSubdeps 2>/dev/null | tr -d '[:space:]')

blocked=false
case "$value" in
  true|True|TRUE) blocked=true ;;
esac

jq -cn \
  --argjson blocked "$blocked" \
  --arg raw "$value" \
  '{blocked: $blocked, raw: $raw}' | emit_ok PNP002
