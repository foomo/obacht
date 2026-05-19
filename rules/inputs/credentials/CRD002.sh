#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# .netrc file permissions.
netrc="$HOME/.netrc"
netrc_exists=false
netrc_mode=""
if [ -f "$netrc" ]; then
  netrc_exists=true
  netrc_mode=$(stat -f '%Lp' "$netrc" 2>/dev/null || stat -c '%a' "$netrc" 2>/dev/null || echo "")
  netrc_mode="0$netrc_mode"
fi

printf '{"netrc_exists": %s, "netrc_mode": "%s"}' "$netrc_exists" "$netrc_mode" | emit_ok CRD002
