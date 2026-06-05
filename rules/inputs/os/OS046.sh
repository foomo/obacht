#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS046
  exit 0
fi

# defaults uses YES/NO for booleans on this key
multicast_disabled=false
value=$(defaults read /Library/Preferences/com.apple.mDNSResponder.plist NoMulticastAdvertisements 2>/dev/null)
if [ "$value" = "1" ] || [ "$value" = "YES" ]; then
  multicast_disabled=true
fi

jq -cn --arg os "$os" --argjson e "$multicast_disabled" '{os: $os, multicast_disabled: $e}' | emit_ok OS046
