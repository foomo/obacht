#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS029
  exit 0
fi

content_caching_disabled=true
content_caching=$(defaults read /Library/Preferences/com.apple.AssetCache.plist Activated 2>/dev/null || echo "0")
if [ "$content_caching" = "1" ]; then
  content_caching_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$content_caching_disabled" '{os: $os, content_caching_disabled: $e}' | emit_ok OS029
