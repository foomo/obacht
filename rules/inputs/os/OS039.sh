#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS039
  exit 0
fi

lockdown_enabled=false
value=$(defaults read .GlobalPreferences LDMGlobalEnabled 2>/dev/null)
if [ "$value" = "1" ]; then
  lockdown_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$lockdown_enabled" '{os: $os, lockdown_enabled: $e}' | emit_ok OS039
