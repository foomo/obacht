#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS032
  exit 0
fi

timemachine_enabled=false
if tmutil destinationinfo 2>/dev/null | grep -q "Name"; then
  if defaults read /Library/Preferences/com.apple.TimeMachine AutoBackup 2>/dev/null | grep -q "1"; then
    timemachine_enabled=true
  fi
fi

timemachine_destination_connected=false
if [ "$timemachine_enabled" = "true" ]; then
  if tmutil destinationinfo 2>/dev/null | grep -qi "^Mount Point"; then
    timemachine_destination_connected=true
  fi
fi

# Time Machine destination encryption.
# Modern APFS-based TM (Ventura+) reports encryption via `diskutil info` as
# `FileVault: Yes`. Legacy CoreStorage-backed HFS+ destinations report
# `Encrypted: Yes` in `diskutil info`. Legacy sparsebundle destinations report
# `Encrypted: Yes` only in `tmutil destinationinfo`. Try diskutil first (both
# keys), fall back to tmutil.
timemachine_destination_encrypted=true
if [ "$timemachine_enabled" = "true" ] && [ "$timemachine_destination_connected" = "true" ]; then
  tm_mount=$(tmutil destinationinfo 2>/dev/null \
    | awk -F': *' '/^Mount Point/ {sub(/^[[:space:]]+/,"",$2); print $2; exit}')
  encrypted=false
  if [ -n "$tm_mount" ]; then
    diskutil_info=$(diskutil info "$tm_mount" 2>/dev/null)
    if printf '%s\n' "$diskutil_info" | grep -qE '^[[:space:]]*FileVault:[[:space:]]+Yes'; then
      encrypted=true
    elif printf '%s\n' "$diskutil_info" | grep -qE '^[[:space:]]*Encrypted:[[:space:]]+Yes'; then
      encrypted=true
    fi
  fi
  if [ "$encrypted" = "false" ]; then
    if tmutil destinationinfo 2>/dev/null | grep -qiE 'Encrypted[[:space:]]*:[[:space:]]*Yes'; then
      encrypted=true
    fi
  fi
  [ "$encrypted" = "false" ] && timemachine_destination_encrypted=false
fi

jq -cn \
  --arg os "$os" \
  --argjson en "$timemachine_enabled" \
  --argjson con "$timemachine_destination_connected" \
  --argjson enc "$timemachine_destination_encrypted" \
  '{os: $os, timemachine_enabled: $en, timemachine_destination_connected: $con, timemachine_destination_encrypted: $enc}' \
  | emit_ok OS032
