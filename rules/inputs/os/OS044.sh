#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS044
  exit 0
fi

show_all_extensions=false
value=$(defaults read NSGlobalDomain AppleShowAllExtensions 2>/dev/null)
if [ "$value" = "1" ]; then
  show_all_extensions=true
fi

jq -cn --arg os "$os" --argjson e "$show_all_extensions" '{os: $os, show_all_extensions: $e}' | emit_ok OS044
