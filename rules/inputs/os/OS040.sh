#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS040
  exit 0
fi

plist="/Library/Application Support/CrashReporter/DiagnosticMessagesHistory.plist"
auto_submit=false
third_party_submit=false

if [ -r "$plist" ]; then
  value=$(defaults read "$plist" AutoSubmit 2>/dev/null)
  if [ "$value" = "1" ]; then
    auto_submit=true
  fi
  value=$(defaults read "$plist" ThirdPartyDataSubmit 2>/dev/null)
  if [ "$value" = "1" ]; then
    third_party_submit=true
  fi
fi

jq -cn --arg os "$os" \
  --argjson a "$auto_submit" \
  --argjson t "$third_party_submit" \
  '{os: $os, auto_submit: $a, third_party_submit: $t}' | emit_ok OS040
