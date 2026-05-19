#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS015
  exit 0
fi

# Printer sharing is disabled when cupsctl reports _share_printers=0.
# If cupsctl fails (CUPS not running), treat as disabled.
printer_sharing_disabled=true
if cupsctl 2>/dev/null | grep -q '^_share_printers=1'; then
  printer_sharing_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$printer_sharing_disabled" '{os: $os, printer_sharing_disabled: $e}' | emit_ok OS015
