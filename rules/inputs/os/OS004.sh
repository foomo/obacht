#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS004
  exit 0
fi

stealth_mode_enabled=false
if /usr/libexec/ApplicationFirewall/socketfilterfw --getstealthmode 2>&1 | grep -q "on"; then
  stealth_mode_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$stealth_mode_enabled" '{os: $os, stealth_mode_enabled: $e}' | emit_ok OS004
