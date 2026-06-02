#!/bin/sh
# shellcheck shell=sh
#
# Duration parsing for package-manager release-age settings.
#
# parse_duration_to_seconds <value>
#   Echoes the duration in seconds as a non-negative integer.
#   Echoes nothing on parse failure (caller treats as "unsupported").
#
#   Accepted formats:
#     - bare integer:        "604800"      -> 604800
#     - seconds:             "604800s"     -> 604800
#     - minutes:             "10080m"      -> 604800
#     - hours:               "168h"        -> 604800
#     - days:                "7d"          -> 604800
#     - weeks:               "1w"          -> 604800
#     - years (365d):        "1y"          -> 31536000
#
#   Mise's "Nm" = months is intentionally not supported here (ambiguous with
#   minutes used by pnpm/yarn). Callers that may receive month/date values
#   should pre-filter and emit_skip when this helper returns empty.

parse_duration_to_seconds() {
  raw="$1"
  [ -z "$raw" ] && return 0

  # Strip surrounding quotes and whitespace.
  v=$(printf '%s' "$raw" | tr -d '"' | tr -d "'" | tr -d '[:space:]')
  [ -z "$v" ] && return 0

  # Bare non-negative integer -> seconds.
  case "$v" in
    *[!0-9]*) ;;
    *) printf '%s\n' "$v"; return 0 ;;
  esac

  # Split into numeric prefix and unit suffix.
  num=$(printf '%s' "$v" | sed -n 's/^\([0-9][0-9]*\)\([a-zA-Z]*\)$/\1/p')
  unit=$(printf '%s' "$v" | sed -n 's/^\([0-9][0-9]*\)\([a-zA-Z]*\)$/\2/p' | tr '[:upper:]' '[:lower:]')

  [ -z "$num" ] && return 0

  case "$unit" in
    s|sec|secs|second|seconds) printf '%s\n' "$num" ;;
    m|min|mins|minute|minutes) printf '%s\n' "$((num * 60))" ;;
    h|hr|hrs|hour|hours)       printf '%s\n' "$((num * 3600))" ;;
    d|day|days)                printf '%s\n' "$((num * 86400))" ;;
    w|wk|wks|week|weeks)       printf '%s\n' "$((num * 604800))" ;;
    y|yr|yrs|year|years)       printf '%s\n' "$((num * 31536000))" ;;
    *) ;;
  esac
}
