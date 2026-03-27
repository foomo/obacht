#!/bin/sh

# Password manager detection.
password_manager_installed=false
password_manager_name=""
for app in "1Password 7" "1Password" "Bitwarden" "Dashlane" "KeePassXC" "LastPass" "Enpass"; do
  if [ -d "/Applications/$app.app" ]; then
    password_manager_installed=true
    password_manager_name="$app"
    break
  fi
done

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

# DNS-over-HTTPS / encrypted DNS detection.
encrypted_dns=false
dns_config=$(scutil --dns 2>/dev/null || true)
if printf '%s' "$dns_config" | grep -qi "encrypted\|https\|tls"; then
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

printf '{"password_manager_installed": %s, "password_manager_name": "%s", "vpn_configured": %s, "encrypted_dns": %s}' \
  "$password_manager_installed" "$password_manager_name" "$vpn_configured" "$encrypted_dns"
