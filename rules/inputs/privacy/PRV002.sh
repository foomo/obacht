#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# VPN configuration detection.
vpn_configured=false
if scutil --nc list 2>/dev/null | grep -q '"'; then
  vpn_configured=true
fi
# Also check for common VPN tools.
if [ "$vpn_configured" = false ]; then
  for vpn in tailscaled wireguard-go openvpn; do
    if pgrep -q "$vpn" 2>/dev/null; then
      vpn_configured=true
      break
    fi
  done
fi

printf '{"vpn_configured": %s}' "$vpn_configured" | emit_ok PRV002
