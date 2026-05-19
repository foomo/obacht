#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

dns_config=$(scutil --dns 2>/dev/null || true)

# Configured nameservers (extract IPs from scutil --dns).
nameservers=$(printf '%s' "$dns_config" | awk '/nameserver\[/ {print $NF}' | sort -u)

# Trusted DNS resolver allowlist (privacy-respecting providers + loopback).
trusted_list="1.1.1.1 1.0.0.1 1.1.1.2 1.0.0.2 1.1.1.3 1.0.0.3 \
9.9.9.9 9.9.9.10 9.9.9.11 149.112.112.112 149.112.112.10 149.112.112.11 \
8.8.8.8 8.8.4.4 \
94.140.14.14 94.140.15.15 94.140.14.140 94.140.14.141 \
194.242.2.2 194.242.2.3 194.242.2.4 194.242.2.5 194.242.2.6 194.242.2.9 \
45.90.28.0 45.90.30.0 \
208.67.222.222 208.67.220.220 \
127.0.0.1 ::1 \
2606:4700:4700::1111 2606:4700:4700::1001 \
2620:fe::fe 2620:fe::9 \
2001:4860:4860::8888 2001:4860:4860::8844"

trusted_dns=true
has_any=false
untrusted_servers=""
for ns in $nameservers; do
  has_any=true
  matched=false
  for trusted in $trusted_list; do
    if [ "$ns" = "$trusted" ]; then
      matched=true
      break
    fi
  done
  if [ "$matched" = false ]; then
    trusted_dns=false
    if [ -z "$untrusted_servers" ]; then
      untrusted_servers="$ns"
    else
      untrusted_servers="$untrusted_servers, $ns"
    fi
  fi
done
if [ "$has_any" = false ]; then
  trusted_dns=false
fi

# Build JSON array of nameservers.
dns_servers_json=$(printf '%s\n' "$nameservers" | awk 'NF{printf "%s\"%s\"", (n++?",":""), $1}')

printf '{"dns_servers": [%s], "trusted_dns": %s, "untrusted_dns_servers": "%s"}' \
  "$dns_servers_json" "$trusted_dns" "$untrusted_servers" | emit_ok PRV004
