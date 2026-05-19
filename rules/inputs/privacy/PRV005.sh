#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# DO_NOT_TRACK env var detection (https://donottrack.sh).
# Value "1" signals opt-out of telemetry to participating CLIs.
value="$(printenv DO_NOT_TRACK 2>/dev/null || true)"
if [ "$value" = "1" ]; then
  enabled=true
else
  enabled=false
fi

jq -cn --argjson enabled "$enabled" --arg value "$value" \
  '{do_not_track_enabled: $enabled, value: $value}' | emit_ok PRV005
