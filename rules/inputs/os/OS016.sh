#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS016
  exit 0
fi

remote_apple_events_disabled=true
if launchctl list com.apple.AEServer >/dev/null 2>&1; then
  remote_apple_events_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$remote_apple_events_disabled" '{os: $os, remote_apple_events_disabled: $e}' | emit_ok OS016
