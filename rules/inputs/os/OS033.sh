#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS033
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

# Time Machine recent backup (within 14 days).
# `tmutil latestbackup` requires Full Disk Access. The TimeMachine preferences
# expose a top-level `LastBackupActivity` key in `YYYY-MM-DD-HHMMSS` format
# that `defaults read` can access without elevated permissions.
timemachine_recent_backup=true
if [ "$timemachine_enabled" = "true" ] && [ "$timemachine_destination_connected" = "true" ]; then
  stamp=$(defaults read /Library/Preferences/com.apple.TimeMachine LastBackupActivity 2>/dev/null \
    | grep -oE '^[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{6}' | head -1)
  if [ -z "$stamp" ]; then
    timemachine_recent_backup=false
  else
    backup_epoch=$(date -j -f "%Y-%m-%d-%H%M%S" "$stamp" +%s 2>/dev/null || echo 0)
    now=$(date +%s)
    if [ "$backup_epoch" -le 0 ]; then
      timemachine_recent_backup=false
    else
      age=$((now - backup_epoch))
      if [ "$age" -gt 1209600 ]; then
        timemachine_recent_backup=false
      fi
    fi
  fi
fi

jq -cn \
  --arg os "$os" \
  --argjson en "$timemachine_enabled" \
  --argjson con "$timemachine_destination_connected" \
  --argjson rec "$timemachine_recent_backup" \
  '{os: $os, timemachine_enabled: $en, timemachine_destination_connected: $con, timemachine_recent_backup: $rec}' \
  | emit_ok OS033
