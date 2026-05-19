#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS002
  exit 0
fi

filevault_enabled=false
if fdesetup status 2>&1 | grep -q "On"; then
  filevault_enabled=true
fi

jq -cn --arg os "$os" --argjson e "$filevault_enabled" '{os: $os, filevault_enabled: $e}' | emit_ok OS002
