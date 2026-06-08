#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

if ! command -v gh >/dev/null 2>&1; then
  emit_skip GH001 "gh not installed"
  exit 0
fi

# Resolve gh config dir per gh's precedence:
#   1. $GH_CONFIG_DIR
#   2. $XDG_CONFIG_HOME/gh
#   3. $HOME/.config/gh
if [ -n "$GH_CONFIG_DIR" ]; then
  config_dir="$GH_CONFIG_DIR"
elif [ -n "$XDG_CONFIG_HOME" ]; then
  config_dir="$XDG_CONFIG_HOME/gh"
else
  config_dir="$HOME/.config/gh"
fi
config_path="$config_dir/config.yml"

telemetry_value=""
if [ -f "$config_path" ]; then
  telemetry_value=$(grep -E '^telemetry:[[:space:]]*' "$config_path" 2>/dev/null \
    | sed -E 's/^telemetry:[[:space:]]*//' \
    | tr -d '[:space:]' \
    | head -n1)
fi

gh_telemetry_env=$(printf '%s' "${GH_TELEMETRY:-}" | tr '[:upper:]' '[:lower:]')
do_not_track_env=$(printf '%s' "${DO_NOT_TRACK:-}" | tr '[:upper:]' '[:lower:]')

jq -cn \
  --arg t "$telemetry_value" \
  --arg ght "$gh_telemetry_env" \
  --arg dnt "$do_not_track_env" \
  --arg cfg "$config_path" \
  '{telemetry: $t, gh_telemetry_env: $ght, do_not_track_env: $dnt, config_path: $cfg}' \
  | emit_ok GH001
