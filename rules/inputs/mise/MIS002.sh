#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# mise not_found_auto_install gating.
#
# When true (mise default), referencing an unknown tool in a mise.toml
# silently installs it on shell activation, which is a supply-chain hazard.
# Setting it to false forces explicit `mise install`.

if ! command -v mise >/dev/null 2>&1; then
  emit_skip MIS002 "mise not installed"
  exit 0
fi

raw=$(mise settings get not_found_auto_install 2>/dev/null | tr -d '[:space:]')

case "$raw" in
  false) value=false ;;
  *)     value=true  ;;
esac

jq -cn --argjson value "$value" '{value: $value}' | emit_ok MIS002
