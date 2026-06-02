#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/duration.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/duration.sh

# mise package-manager release-age gating.
#
# Reads mise's `minimum_release_age` setting. Accepts duration strings
# ("7d", "90d", "1y") or absolute dates ("2024-06-01"). The shared
# duration helper handles the relative form; an absolute date is
# converted to age in seconds against the current time.

if ! command -v mise >/dev/null 2>&1; then
  emit_skip MIS001 "mise not installed"
  exit 0
fi

raw=$(mise settings get minimum_release_age 2>/dev/null | tr -d '[:space:]')

date_to_epoch() {
  # Strip any time component to keep parsing portable.
  d=$(printf '%s' "$1" | cut -d T -f 1)
  # BSD date (macOS) first, then GNU date (Linux).
  date -j -f "%Y-%m-%d" "$d" "+%s" 2>/dev/null \
    || date -d "$d" "+%s" 2>/dev/null
}

configured=false
value_seconds=0
unsupported=false

if [ -n "$raw" ] && [ "$raw" != "null" ] && [ "$raw" != "undefined" ]; then
  case "$raw" in
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]*)
      epoch=$(date_to_epoch "$raw")
      now=$(date +%s)
      if [ -n "$epoch" ] && [ "$epoch" -gt 0 ] && [ "$epoch" -le "$now" ]; then
        configured=true
        value_seconds=$((now - epoch))
      else
        unsupported=true
      fi
      ;;
    *)
      parsed=$(parse_duration_to_seconds "$raw")
      case "$parsed" in
        '')
          unsupported=true
          ;;
        0)
          ;;
        *[!0-9]*)
          unsupported=true
          ;;
        *)
          configured=true
          value_seconds="$parsed"
          ;;
      esac
      ;;
  esac
fi

if [ "$unsupported" = true ]; then
  emit_skip MIS001 "unsupported mise minimum_release_age format: $raw"
  exit 0
fi

jq -cn \
  --argjson configured "$configured" \
  --argjson value_seconds "$value_seconds" \
  '{configured: $configured, value_seconds: $value_seconds}' | emit_ok MIS001
