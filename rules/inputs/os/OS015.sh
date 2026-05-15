#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS015
  exit 0
fi

# Preserve legacy semantics: cupsctl success => disabled reported true.
printer_sharing_disabled=false
if cupsctl >/dev/null 2>&1; then
  printer_sharing_disabled=true
fi

jq -cn --arg os "$os" --argjson e "$printer_sharing_disabled" '{os: $os, printer_sharing_disabled: $e}' | emit_ok OS015
