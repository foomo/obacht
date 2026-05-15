#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS014
  exit 0
fi

internet_sharing_disabled=true
if defaults read /Library/Preferences/SystemConfiguration/com.apple.nat Enabled 2>/dev/null | grep -q 1; then
  internet_sharing_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$internet_sharing_disabled" '{os: $os, internet_sharing_disabled: $e}' | emit_ok OS014
