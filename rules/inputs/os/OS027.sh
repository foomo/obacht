#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS027
  exit 0
fi

media_sharing_disabled=true
media_sharing=$(defaults read com.apple.amp.mediasharingd home-sharing-enabled 2>/dev/null || echo "0")
if [ "$media_sharing" = "1" ]; then
  media_sharing_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$media_sharing_disabled" '{os: $os, media_sharing_disabled: $e}' | emit_ok OS027
