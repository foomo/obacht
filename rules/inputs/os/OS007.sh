#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS007
  exit 0
fi

# Preserve legacy semantics: defaults read success => "disabled" reported true.
guest_account_disabled=false
if defaults read /Library/Preferences/com.apple.loginwindow GuestEnabled >/dev/null 2>&1; then
  guest_account_disabled=true
fi

jq -cn --arg os "$os" --argjson e "$guest_account_disabled" '{os: $os, guest_account_disabled: $e}' | emit_ok OS007
