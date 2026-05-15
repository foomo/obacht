#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS030
  exit 0
fi

user_is_admin=false
if groups 2>/dev/null | tr ' ' '\n' | grep -qx admin; then
  user_is_admin=true
fi

jq -cn --arg os "$os" --argjson e "$user_is_admin" '{os: $os, user_is_admin: $e}' | emit_ok OS030
