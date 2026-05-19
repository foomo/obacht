#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS020
  exit 0
fi

mdm_enrolled=false
if profiles status -type enrollment 2>/dev/null | grep -q "MDM enrollment: Yes"; then
  mdm_enrolled=true
fi

jq -cn --arg os "$os" --argjson e "$mdm_enrolled" '{os: $os, mdm_enrolled: $e}' | emit_ok OS020
