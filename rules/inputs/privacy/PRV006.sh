#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok PRV006
  exit 0
fi

# User-domain trust settings (readable without sudo).
# Non-empty list = user has explicitly added trust overrides for one or more
# certificates -- a common indicator of attacker- or MDM-installed root CAs.
user_count=0
admin_count=0

if command -v security >/dev/null 2>&1; then
  user_dump=$(security dump-trust-settings 2>/dev/null || true)
  user_count=$(printf '%s\n' "$user_dump" | grep -cE '^Cert ' || true)

  admin_dump=$(security dump-trust-settings -d 2>/dev/null || true)
  admin_count=$(printf '%s\n' "$admin_dump" | grep -cE '^Cert ' || true)
else
  emit_skip PRV006 "security command not available"
  exit 0
fi

jq -cn --arg os "$os" --argjson u "$user_count" --argjson a "$admin_count" \
  '{os: $os, user_trust_overrides: $u, admin_trust_overrides: $a}' | emit_ok PRV006
