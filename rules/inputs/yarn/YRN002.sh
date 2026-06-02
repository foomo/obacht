#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# yarn install-script execution gating.
#
# Yarn berry: setting `enableScripts: false` in .yarnrc.yml disables
# lifecycle script execution for all installs — the equivalent of
# npm's `ignore-scripts=true`. Yarn classic (1.x) lacks an equivalent
# global switch and is skipped.

if ! command -v yarn >/dev/null 2>&1; then
  emit_skip YRN002 "yarn not installed"
  exit 0
fi

ver=$(yarn --version 2>/dev/null | tr -d '[:space:]')
case "$ver" in
  1.*|0.*)
    emit_skip YRN002 "yarn classic ($ver) has no global enableScripts toggle"
    exit 0
    ;;
esac

value=$(yarn config get enableScripts 2>/dev/null | tr -d '[:space:]')

scripts_disabled=false
case "$value" in
  false|"False"|"FALSE") scripts_disabled=true ;;
esac

jq -cn \
  --argjson scripts_disabled "$scripts_disabled" \
  --arg raw "$value" \
  '{scripts_disabled: $scripts_disabled, raw: $raw}' | emit_ok YRN002
