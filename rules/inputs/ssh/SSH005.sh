#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

set -e
ssh_dir="$HOME/.ssh"

if [ ! -d "$ssh_dir" ]; then
  printf '{"keys": []}' | emit_ok SSH005
  exit 0
fi

# Resolve symlinks.
ssh_dir=$(cd "$ssh_dir" && pwd -P)

# Discover private keys with algorithm and bits.
first=true
key_json="["
for f in "$ssh_dir"/id_*; do
  [ -e "$f" ] || continue
  case "$f" in *.pub) continue;; esac
  pub="${f}.pub"
  bits=0
  algo=""
  if [ -f "$pub" ]; then
    line=$(ssh-keygen -l -f "$pub" 2>/dev/null || echo "")
    bits=$(printf '%s' "$line" | awk '{print $1}' | tr -dc '0-9')
    [ -z "$bits" ] && bits=0
    algo=$(printf '%s' "$line" | sed -n 's/.*(\([A-Z0-9]*\)).*/\1/p')
  fi
  if [ "$first" = true ]; then
    first=false
  else
    key_json="$key_json,"
  fi
  key_json="$key_json{\"path\":\"$f\",\"bits\":$bits,\"algorithm\":\"$algo\"}"
done
key_json="$key_json]"

printf '{"keys": %s}' "$key_json" | emit_ok SSH005
