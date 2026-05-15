#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# DNS-over-HTTPS / encrypted DNS detection.
encrypted_dns=false
dns_config=$(scutil --dns 2>/dev/null || true)
if printf '%s' "$dns_config" | grep -qE "dns_over_https|dns_over_tls"; then
  encrypted_dns=true
fi
# Check for common DNS tools.
if [ "$encrypted_dns" = false ]; then
  for dns in dnscrypt-proxy cloudflared; do
    if pgrep -q "$dns" 2>/dev/null; then
      encrypted_dns=true
      break
    fi
  done
fi

printf '{"encrypted_dns": %s}' "$encrypted_dns" | emit_ok PRV003
