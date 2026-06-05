#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS041
  exit 0
fi

env_keep_home=false
for f in /etc/sudoers /etc/sudoers.d/*; do
  [ -r "$f" ] || continue
  # match: Defaults env_keep += "... HOME ..." or = "..."
  if grep -E '^[[:space:]]*Defaults[[:space:]]+env_keep[[:space:]]*[+]?=.*\bHOME\b' "$f" >/dev/null 2>&1; then
    env_keep_home=true
    break
  fi
done

jq -cn --arg os "$os" --argjson e "$env_keep_home" '{os: $os, env_keep_home: $e}' | emit_ok OS041
