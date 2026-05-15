#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

set -e
ssh_dir="$HOME/.ssh"

if [ ! -d "$ssh_dir" ]; then
  printf '{"directory_exists": false, "directory_mode": ""}' | emit_ok SSH002
  exit 0
fi

# Resolve symlinks.
ssh_dir=$(cd "$ssh_dir" && pwd -P)

dir_mode=$(stat -f '%Lp' "$ssh_dir" 2>/dev/null || stat -c '%a' "$ssh_dir" 2>/dev/null || echo "")

printf '{"directory_exists": true, "directory_mode": "0%s"}' "$dir_mode" | emit_ok SSH002
