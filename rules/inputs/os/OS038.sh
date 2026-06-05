#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS038
  exit 0
fi

socketfilterfw=/usr/libexec/ApplicationFirewall/socketfilterfw
if [ ! -x "$socketfilterfw" ]; then
  emit_skip OS038 "socketfilterfw not available"
  exit 0
fi

allow_signed_enabled=false
allow_signed_app_enabled=false
if "$socketfilterfw" --getallowsigned 2>/dev/null | grep -qi "Automatically allow built-in software.*enabled"; then
  allow_signed_enabled=true
fi
if "$socketfilterfw" --getallowsigned 2>/dev/null | grep -qi "Automatically allow downloaded signed software.*enabled"; then
  allow_signed_app_enabled=true
fi

jq -cn --arg os "$os" \
  --argjson b "$allow_signed_enabled" \
  --argjson a "$allow_signed_app_enabled" \
  '{os: $os, allow_signed_builtin: $b, allow_signed_downloaded: $a}' | emit_ok OS038
