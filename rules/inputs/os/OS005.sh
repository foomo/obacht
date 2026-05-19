#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS005
  exit 0
fi

gatekeeper_enabled=false
if spctl --status 2>&1 | grep -q "enabled"; then
  gatekeeper_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$gatekeeper_enabled" '{os: $os, gatekeeper_enabled: $e}' | emit_ok OS005
