#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS011
  exit 0
fi

# Rapid Security Response is enabled when the ConfigDataInstall key is set to 1.
rsr_enabled=false
value=$(defaults read /Library/Preferences/com.apple.SoftwareUpdate ConfigDataInstall 2>/dev/null)
if [ "$value" = "1" ]; then
  rsr_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$rsr_enabled" '{os: $os, rsr_enabled: $e}' | emit_ok OS011
