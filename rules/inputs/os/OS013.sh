#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS013
  exit 0
fi

screen_sharing_disabled=true
if launchctl list com.apple.screensharing >/dev/null 2>&1; then
  screen_sharing_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$screen_sharing_disabled" '{os: $os, screen_sharing_disabled: $e}' | emit_ok OS013
