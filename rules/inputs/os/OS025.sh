#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS025
  exit 0
fi

# Only treat ARD as enabled when an actual PID is attached.
remote_management_disabled=true
ard_pid=$(launchctl list 2>/dev/null | awk '$3 == "com.apple.RemoteDesktop.agent" {print $1; exit}')
if [ -n "$ard_pid" ] && [ "$ard_pid" != "-" ]; then
  remote_management_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$remote_management_disabled" '{os: $os, remote_management_disabled: $e}' | emit_ok OS025
