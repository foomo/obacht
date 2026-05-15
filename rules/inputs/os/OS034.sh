#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS034
  exit 0
fi

airplay_receiver_enabled=false
airplay_val=$(defaults -currentHost read com.apple.controlcenter AirplayRecieverEnabled 2>/dev/null || echo "")
[ "$airplay_val" = "1" ] && airplay_receiver_enabled=true

jq -cn --arg os "$os" --argjson e "$airplay_receiver_enabled" '{os: $os, airplay_receiver_enabled: $e}' | emit_ok OS034
