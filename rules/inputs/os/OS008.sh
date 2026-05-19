#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS008
  exit 0
fi

timeout_raw=$(defaults -currentHost read com.apple.screensaver idleTime 2>/dev/null || echo "0")
screen_lock_timeout_seconds=$(printf '%s' "$timeout_raw" | tr -dc '0-9')
[ -z "$screen_lock_timeout_seconds" ] && screen_lock_timeout_seconds=0

jq -cn --arg os "$os" --argjson t "$screen_lock_timeout_seconds" '{os: $os, screen_lock_timeout_seconds: $t}' | emit_ok OS008
