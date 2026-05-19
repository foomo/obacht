#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Go telemetry opt-out check.
if ! command -v go >/dev/null 2>&1; then
  emit_skip TOL005 "go CLI not installed"
  exit 0
fi

mode="$(go env GOTELEMETRY 2>/dev/null)"
if [ "$mode" = "off" ]; then
  telemetry_disabled=true
else
  telemetry_disabled=false
fi

jq -cn --argjson disabled "$telemetry_disabled" --arg mode "$mode" \
  '{telemetry_disabled: $disabled, mode: $mode}' | emit_ok TOL005
