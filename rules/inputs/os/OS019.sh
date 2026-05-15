#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS019
  exit 0
fi

legacy_kexts_blocked=true
if kmutil showloaded --list-only 2>/dev/null | grep -q "com.apple"; then
  legacy_kexts_blocked=false
fi

jq -cn --arg os "$os" --argjson e "$legacy_kexts_blocked" '{os: $os, legacy_kexts_blocked: $e}' | emit_ok OS019
