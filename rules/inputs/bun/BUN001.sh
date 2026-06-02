#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/duration.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/duration.sh

# bun package-manager release-age gating.
#
# Bun (>=1.3) reads `[install].minimumReleaseAge` (seconds) from bunfig.toml.
# Bun has no `config get` subcommand for this, so we read the global bunfig
# directly, honoring $BUN_CONFIG_FILE and $XDG_CONFIG_HOME.

if ! command -v bun >/dev/null 2>&1; then
  emit_skip BUN001 "bun not installed"
  exit 0
fi

resolve_bunfig() {
  if [ -n "${BUN_CONFIG_FILE:-}" ] && [ -f "$BUN_CONFIG_FILE" ]; then
    printf '%s' "$BUN_CONFIG_FILE"
    return
  fi
  if [ -n "${XDG_CONFIG_HOME:-}" ] && [ -f "$XDG_CONFIG_HOME/.bunfig.toml" ]; then
    printf '%s' "$XDG_CONFIG_HOME/.bunfig.toml"
    return
  fi
  if [ -f "$HOME/.config/.bunfig.toml" ]; then
    printf '%s' "$HOME/.config/.bunfig.toml"
    return
  fi
  if [ -f "$HOME/.bunfig.toml" ]; then
    printf '%s' "$HOME/.bunfig.toml"
  fi
}

bunfig=$(resolve_bunfig)
raw=""
if [ -n "$bunfig" ]; then
  raw=$(awk '
    /^[[:space:]]*\[/ {
      gsub(/^[[:space:]]+|[[:space:]]+$/, "")
      section = $0
      next
    }
    /^[[:space:]]*minimumReleaseAge[[:space:]]*=/ {
      if (section == "[install]") {
        sub(/^[^=]*=[[:space:]]*/, "")
        gsub(/[[:space:]]/, "")
        gsub(/^"|"$|^'\''|'\''$/, "")
        print
        exit
      }
    }
  ' "$bunfig" 2>/dev/null)
fi

configured=false
value_seconds=0
if [ -n "$raw" ]; then
  parsed=$(parse_duration_to_seconds "$raw")
  case "$parsed" in
    ''|0) ;;
    *[!0-9]*) ;;
    *)
      configured=true
      value_seconds="$parsed"
      ;;
  esac
fi

jq -cn \
  --argjson configured "$configured" \
  --argjson value_seconds "$value_seconds" \
  '{configured: $configured, value_seconds: $value_seconds}' | emit_ok BUN001
