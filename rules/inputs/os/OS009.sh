#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS009
  exit 0
fi

# Preserve legacy semantics: defaults read success => enabled reported true.
os_auto_update_enabled=false
if defaults read /Library/Preferences/com.apple.SoftwareUpdate AutomaticallyInstallMacOSUpdates >/dev/null 2>&1; then
  os_auto_update_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$os_auto_update_enabled" '{os: $os, os_auto_update_enabled: $e}' | emit_ok OS009
