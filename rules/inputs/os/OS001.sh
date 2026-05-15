#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS001
  exit 0
fi

sip_enabled=false
if csrutil status 2>&1 | grep -q "enabled"; then
  sip_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$sip_enabled" '{os: $os, sip_enabled: $e}' | emit_ok OS001
