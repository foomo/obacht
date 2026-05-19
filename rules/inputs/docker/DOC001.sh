#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Inspect the Docker socket permissions.
socket="/var/run/docker.sock"
socket_exists=false
socket_mode=""

if [ -e "$socket" ]; then
  # Resolve symlinks.
  real_socket=$(readlink -f "$socket" 2>/dev/null || echo "$socket")
  socket_exists=true
  socket_mode=$(stat -f '%04Lp' "$real_socket" 2>/dev/null || stat -c '%04a' "$real_socket" 2>/dev/null || echo "")
fi

printf '{"socket_exists": %s, "socket_mode": "%s"}' \
  "$socket_exists" "$socket_mode" | emit_ok DOC001
