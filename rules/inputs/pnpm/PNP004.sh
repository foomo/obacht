#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# pnpm strict-dep-builds gating.
#
# pnpm 10.3+ supports `strictDepBuilds`. With true, pnpm exits with a
# non-zero code if any dependency tries to run a lifecycle script that
# has not been explicitly allowed (via `allowBuilds` /
# `onlyBuiltDependencies`). This turns silent warnings into
# CI-blocking failures and closes the loophole where a transitive
# package quietly requests a build script.

if ! command -v pnpm >/dev/null 2>&1; then
  emit_skip PNP004 "pnpm not installed"
  exit 0
fi

value=$(pnpm config get strictDepBuilds 2>/dev/null | tr -d '[:space:]')

strict=false
case "$value" in
  true|True|TRUE) strict=true ;;
esac

jq -cn \
  --argjson strict "$strict" \
  --arg raw "$value" \
  '{strict: $strict, raw: $raw}' | emit_ok PNP004
