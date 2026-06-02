#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# bun install-script execution gating.
#
# With `[install] ignoreScripts = true` in bunfig.toml, bun refuses to
# run lifecycle scripts (preinstall, install, postinstall, prepare) of
# dependencies. This is the primary mitigation against compromised
# packages that deliver payloads via postinstall hooks. Bun has no
# `config get` subcommand, so we parse the resolved bunfig directly.

if ! command -v bun >/dev/null 2>&1; then
  emit_skip BUN002 "bun not installed"
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
    /^[[:space:]]*ignoreScripts[[:space:]]*=/ {
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

ignore_scripts=false
case "$raw" in
  true|True|TRUE) ignore_scripts=true ;;
esac

jq -cn \
  --argjson ignore_scripts "$ignore_scripts" \
  --arg raw "$raw" \
  '{ignore_scripts: $ignore_scripts, raw: $raw}' | emit_ok BUN002
