#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Check whether the current user is in the docker group.
user_in_group=false
if id -nG 2>/dev/null | tr ' ' '\n' | grep -qx docker; then
  user_in_group=true
fi

printf '{"user_in_group": %s}' "$user_in_group" | emit_ok DOC002
