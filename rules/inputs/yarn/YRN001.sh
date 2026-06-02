#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/duration.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/duration.sh

# yarn package-manager release-age gating.
#
# Reads yarn berry's `npmMinimalAgeGate` setting (yarn 4.10+).
# Numeric values are interpreted as minutes; duration strings like "7d"
# are parsed via the shared duration helper. Yarn classic (1.x) lacks
# the setting and is skipped.

if ! command -v yarn >/dev/null 2>&1; then
  emit_skip YRN001 "yarn not installed"
  exit 0
fi

ver=$(yarn --version 2>/dev/null | tr -d '[:space:]')
case "$ver" in
  1.*|0.*)
    emit_skip YRN001 "yarn classic ($ver) does not support release-age gating"
    exit 0
    ;;
esac

raw=$(yarn config get npmMinimalAgeGate 2>/dev/null | tr -d '[:space:]')

configured=false
value_seconds=0
case "$raw" in
  ''|undefined|null|false|0) ;;
  *)
    case "$raw" in
      *[!0-9]*)
        parsed=$(parse_duration_to_seconds "$raw")
        ;;
      *)
        parsed=$((raw * 60))
        ;;
    esac
    case "$parsed" in
      ''|0) ;;
      *[!0-9]*) ;;
      *)
        configured=true
        value_seconds="$parsed"
        ;;
    esac
    ;;
esac

jq -cn \
  --argjson configured "$configured" \
  --argjson value_seconds "$value_seconds" \
  '{configured: $configured, value_seconds: $value_seconds}' | emit_ok YRN001
