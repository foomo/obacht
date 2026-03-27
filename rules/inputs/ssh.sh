#!/bin/sh
set -e
home="$HOME"
ssh_dir="$home/.ssh"

# Check if ~/.ssh exists.
if [ ! -d "$ssh_dir" ]; then
  printf '{"directory_exists": false, "directory_mode": "", "keys": [], "config_exists": false, "strict_host_key_checking_disabled": false, "forward_agent_global": false}'
  exit 0
fi

# Resolve symlinks.
ssh_dir=$(cd "$ssh_dir" && pwd -P)

dir_mode=$(stat -f '%Lp' "$ssh_dir" 2>/dev/null || stat -c '%a' "$ssh_dir" 2>/dev/null || echo "")

# Discover private keys.
keys="[]"
first=true
key_json="["
for f in "$ssh_dir"/id_*; do
  [ -e "$f" ] || continue
  case "$f" in *.pub) continue;; esac
  mode=$(stat -f '%Lp' "$f" 2>/dev/null || stat -c '%a' "$f" 2>/dev/null || echo "")
  type=$(basename "$f" | sed 's/^id_//')
  if [ "$first" = true ]; then
    first=false
  else
    key_json="$key_json,"
  fi
  key_json="$key_json{\"path\":\"$f\",\"mode\":\"0$mode\",\"type\":\"$type\"}"
done
key_json="$key_json]"

config_exists=false
[ -f "$ssh_dir/config" ] && config_exists=true

# Parse SSH config for risky settings.
strict_host_key_checking_disabled=false
forward_agent_global=false
if [ "$config_exists" = true ]; then
  # Check for StrictHostKeyChecking no (any Host section).
  if grep -i "StrictHostKeyChecking" "$ssh_dir/config" 2>/dev/null | grep -qi "no"; then
    strict_host_key_checking_disabled=true
  fi
  # Check for ForwardAgent yes in global (Host *) context.
  if awk '/^Host \*/{found=1; next} /^Host /{found=0} found && /ForwardAgent.*yes/i{print; exit}' "$ssh_dir/config" 2>/dev/null | grep -qi "yes"; then
    forward_agent_global=true
  fi
fi

printf '{"directory_exists": true, "directory_mode": "0%s", "keys": %s, "config_exists": %s, "strict_host_key_checking_disabled": %s, "forward_agent_global": %s}' \
  "$dir_mode" "$key_json" "$config_exists" "$strict_host_key_checking_disabled" "$forward_agent_global"
