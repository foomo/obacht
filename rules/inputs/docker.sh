#!/bin/sh
# Check if docker is installed.
if ! command -v docker >/dev/null 2>&1; then
  printf '{"installed": false, "socket_exists": false, "socket_mode": "", "user_in_group": false}'
  exit 0
fi

socket="/var/run/docker.sock"
socket_exists=false
socket_mode=""

if [ -e "$socket" ]; then
  # Resolve symlinks.
  real_socket=$(readlink -f "$socket" 2>/dev/null || echo "$socket")
  socket_exists=true
  socket_mode=$(stat -f '%04Lp' "$real_socket" 2>/dev/null || stat -c '%04a' "$real_socket" 2>/dev/null || echo "")
fi

# Check if current user is in the docker group.
user_in_group=false
if id -nG 2>/dev/null | tr ' ' '\n' | grep -qx docker; then
  user_in_group=true
fi

printf '{"installed": true, "socket_exists": %s, "socket_mode": "%s", "user_in_group": %s}' \
  "$socket_exists" "$socket_mode" "$user_in_group"
