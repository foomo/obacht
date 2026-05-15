#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

set -e
ssh_dir="$HOME/.ssh"
config_file="$ssh_dir/config"

if [ ! -f "$config_file" ]; then
  printf '{"config_exists": false, "strict_host_key_checking_disabled": false}' | emit_ok SSH003
  exit 0
fi

# Strip comments before matching.
strict_host_key_checking_disabled=false
uncommented=$(sed -e 's/[[:space:]]*#.*$//' "$config_file" 2>/dev/null)
if printf '%s\n' "$uncommented" | grep -qiE '^[[:space:]]*StrictHostKeyChecking[[:space:]]+no([[:space:]]|$)'; then
  strict_host_key_checking_disabled=true
fi

printf '{"config_exists": true, "strict_host_key_checking_disabled": %s}' \
  "$strict_host_key_checking_disabled" | emit_ok SSH003
