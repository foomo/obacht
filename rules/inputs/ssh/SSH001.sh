#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

set -e
ssh_dir="$HOME/.ssh"

if [ ! -d "$ssh_dir" ]; then
  printf '{"keys": []}' | emit_ok SSH001
  exit 0
fi

# Resolve symlinks.
ssh_dir=$(cd "$ssh_dir" && pwd -P)

# Discover private keys and capture mode.
first=true
key_json="["
for f in "$ssh_dir"/id_*; do
  [ -e "$f" ] || continue
  case "$f" in *.pub) continue;; esac
  mode=$(stat -f '%Lp' "$f" 2>/dev/null || stat -c '%a' "$f" 2>/dev/null || echo "")
  if [ "$first" = true ]; then
    first=false
  else
    key_json="$key_json,"
  fi
  key_json="$key_json{\"path\":\"$f\",\"mode\":\"0$mode\"}"
done
key_json="$key_json]"

printf '{"keys": %s}' "$key_json" | emit_ok SSH001
