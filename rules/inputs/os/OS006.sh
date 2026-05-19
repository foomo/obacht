#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS006
  exit 0
fi

# AutoLogin: if defaults read succeeds, auto-login IS enabled (bad).
auto_login_disabled=true
if defaults read /Library/Preferences/com.apple.loginwindow autoLoginUser >/dev/null 2>&1; then
  auto_login_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$auto_login_disabled" '{os: $os, auto_login_disabled: $e}' | emit_ok OS006
