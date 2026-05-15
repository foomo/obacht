#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS031
  exit 0
fi

# Modern macOS (Ventura+) uses sysadminctl; pre-Ventura uses legacy defaults keys.
# Sentinels for password_lock_delay_seconds:
#   -1 = unknown (rule skipped)
#   -2 = screen lock disabled (rule fires with different evidence)
#    0 = immediate (rule passes)
#   >0 = delay in seconds (rule fires)
password_required_on_lock=false
password_lock_delay_seconds=-1
sl_out=$(sysadminctl -screenLock status 2>&1 || true)
if printf '%s' "$sl_out" | grep -qi "immediate"; then
  password_required_on_lock=true
  password_lock_delay_seconds=0
elif printf '%s' "$sl_out" | grep -qiE 'set to [0-9]+'; then
  n=$(printf '%s' "$sl_out" | grep -oE 'set to [0-9]+' | grep -oE '[0-9]+' | head -1)
  [ -n "$n" ] && password_lock_delay_seconds=$n
elif printf '%s' "$sl_out" | grep -qiE 'screenLock is OFF|disabled'; then
  password_lock_delay_seconds=-2
fi
# Fallback to legacy keys (pre-Ventura) when sysadminctl gave no answer.
if [ "$password_lock_delay_seconds" = "-1" ]; then
  ask_pw=$(defaults read com.apple.screensaver askForPassword 2>/dev/null | tr -dc '0-9')
  ask_pw_delay=$(defaults read com.apple.screensaver askForPasswordDelay 2>/dev/null | tr -dc '0-9')
  if [ -n "$ask_pw" ] && [ -n "$ask_pw_delay" ]; then
    password_lock_delay_seconds=$ask_pw_delay
    if [ "$ask_pw" = "1" ] && [ "$ask_pw_delay" = "0" ]; then
      password_required_on_lock=true
    fi
  fi
fi

jq -cn \
  --arg os "$os" \
  --argjson req "$password_required_on_lock" \
  --argjson delay "$password_lock_delay_seconds" \
  '{os: $os, password_required_on_lock: $req, password_lock_delay_seconds: $delay}' | emit_ok OS031
