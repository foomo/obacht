#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS023
  exit 0
fi

timemachine_enabled=false
if tmutil destinationinfo 2>/dev/null | grep -q "Name"; then
  if defaults read /Library/Preferences/com.apple.TimeMachine AutoBackup 2>/dev/null | grep -q "1"; then
    timemachine_enabled=true
  fi
fi

jq -cn --arg os "$os" --argjson e "$timemachine_enabled" '{os: $os, timemachine_enabled: $e}' | emit_ok OS023
