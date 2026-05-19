#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Homebrew analytics opt-out check.
if ! command -v brew >/dev/null 2>&1; then
  emit_skip TOL004 "brew CLI not installed"
  exit 0
fi

if [ "${HOMEBREW_NO_ANALYTICS:-}" = "1" ]; then
  analytics_disabled=true
else
  analytics_disabled=false
fi

jq -cn --argjson disabled "$analytics_disabled" \
  --arg value "${HOMEBREW_NO_ANALYTICS:-}" \
  '{analytics_disabled: $disabled, value: $value}' | emit_ok TOL004
