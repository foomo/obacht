#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Inspect the kubeconfig file permissions.
config="$HOME/.kube/config"

# Resolve symlinks.
if [ -L "$config" ]; then
  config=$(readlink -f "$config" 2>/dev/null || echo "$config")
fi

if [ ! -f "$config" ]; then
  printf '{"config_exists": false, "config_mode": ""}' | emit_ok KUB001
  exit 0
fi

config_mode=$(stat -f '%04Lp' "$config" 2>/dev/null || stat -c '%04a' "$config" 2>/dev/null || echo "")

printf '{"config_exists": true, "config_mode": "%s"}' "$config_mode" | emit_ok KUB001
