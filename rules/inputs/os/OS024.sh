#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS024
  exit 0
fi

remote_login_disabled=true
if systemsetup -getremotelogin 2>/dev/null | grep -qi "on"; then
  remote_login_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$remote_login_disabled" '{os: $os, remote_login_disabled: $e}' | emit_ok OS024
