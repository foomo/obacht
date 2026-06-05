#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS042
  exit 0
fi

touchid_configured=false
for f in /etc/pam.d/sudo_local /etc/pam.d/sudo; do
  [ -r "$f" ] || continue
  if grep -E '^[^#]*pam_tid\.so' "$f" >/dev/null 2>&1; then
    touchid_configured=true
    break
  fi
done

jq -cn --arg os "$os" --argjson e "$touchid_configured" '{os: $os, touchid_configured: $e}' | emit_ok OS042
