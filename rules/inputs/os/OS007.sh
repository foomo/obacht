#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS007
  exit 0
fi

# Guest account is disabled when the GuestEnabled key is absent OR set to 0.
guest_account_disabled=true
value=$(defaults read /Library/Preferences/com.apple.loginwindow GuestEnabled 2>/dev/null)
if [ -n "$value" ] && [ "$value" != "0" ]; then
  guest_account_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$guest_account_disabled" '{os: $os, guest_account_disabled: $e}' | emit_ok OS007
