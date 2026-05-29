#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD039

required='Read(~/.ssh/**)
Read(~/.aws/**)
Read(~/.gnupg/**)
Read(~/.config/gh/**)
Read(~/.kube/**)
Read(~/.docker/config.json)'

settings=$(claude_settings_path)
if [ -n "$settings" ]; then
  deny=$(jq -c '.permissions.deny // []' "$settings" 2>/dev/null || printf '[]')
else
  deny='[]'
fi

missing=""
saved_ifs=$IFS
IFS='
'
for needle in $required; do
  if ! printf '%s' "$deny" | jq -e --arg n "$needle" 'any(.[]; . == $n)' >/dev/null 2>&1; then
    missing="${missing:+$missing }$needle"
  fi
done
IFS=$saved_ifs

jq -cn --arg v "$missing" '{permissions_deny_home_secrets_missing: $v}' | emit_ok CLD039
