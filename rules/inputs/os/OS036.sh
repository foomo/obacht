#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS036
  exit 0
fi

macos_major=$(sw_vers -productVersion 2>/dev/null | cut -d. -f1)
macos_major=$(printf '%s' "$macos_major" | tr -dc '0-9')
[ -z "$macos_major" ] && macos_major=0

jq -cn --arg os "$os" --argjson v "$macos_major" '{os: $os, macos_major_version: $v}' | emit_ok OS036
