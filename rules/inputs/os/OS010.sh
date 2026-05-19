#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS010
  exit 0
fi

# Enabled when the AutoUpdate key is set to 1.
app_auto_update_enabled=false
value=$(defaults read /Library/Preferences/com.apple.commerce AutoUpdate 2>/dev/null)
if [ "$value" = "1" ]; then
  app_auto_update_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$app_auto_update_enabled" '{os: $os, app_auto_update_enabled: $e}' | emit_ok OS010
