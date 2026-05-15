#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS003
  exit 0
fi

firewall_enabled=false
if /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate 2>&1 | grep -q "enabled"; then
  firewall_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$firewall_enabled" '{os: $os, firewall_enabled: $e}' | emit_ok OS003
