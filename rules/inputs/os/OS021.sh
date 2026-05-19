#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS021
  exit 0
fi

rosetta_installed=false
if pgrep -q oahd 2>/dev/null; then
  rosetta_installed=true
fi

jq -cn --arg os "$os" --argjson e "$rosetta_installed" '{os: $os, rosetta_installed: $e}' | emit_ok OS021
