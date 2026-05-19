#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

set -e
ssh_dir="$HOME/.ssh"
config_file="$ssh_dir/config"

if [ ! -f "$config_file" ]; then
  printf '{"config_exists": false, "forward_agent_global": false}' | emit_ok SSH004
  exit 0
fi

# Strip comments before matching.
forward_agent_global=false
uncommented=$(sed -e 's/[[:space:]]*#.*$//' "$config_file" 2>/dev/null)
# Check for ForwardAgent yes in global (Host *) context.
if printf '%s\n' "$uncommented" | awk 'BEGIN{IGNORECASE=1} /^[[:space:]]*Host[[:space:]]+\*[[:space:]]*$/{found=1; next} /^[[:space:]]*Host[[:space:]]/{found=0} found && /^[[:space:]]*ForwardAgent[[:space:]]+yes([[:space:]]|$)/{print; exit}' | grep -qi "yes"; then
  forward_agent_global=true
fi

printf '{"config_exists": true, "forward_agent_global": %s}' \
  "$forward_agent_global" | emit_ok SSH004
