#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS011
  exit 0
fi

# Preserve legacy semantics: defaults read success => enabled reported true.
rsr_enabled=false
if defaults read /Library/Preferences/com.apple.SoftwareUpdate ConfigDataInstall >/dev/null 2>&1; then
  rsr_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$rsr_enabled" '{os: $os, rsr_enabled: $e}' | emit_ok OS011
