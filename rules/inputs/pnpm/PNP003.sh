#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# pnpm trust-policy gating.
#
# pnpm 10.21+ supports `trustPolicy`. With `no-downgrade`, pnpm
# refuses to install a package version whose trust evidence (signed
# publisher / provenance) has regressed compared to any earlier
# release — defeating compromised-publisher attacks where the
# attacker drops back to unsigned releases. Default is unset (no
# trust check performed).

if ! command -v pnpm >/dev/null 2>&1; then
  emit_skip PNP003 "pnpm not installed"
  exit 0
fi

value=$(pnpm config get trustPolicy 2>/dev/null | tr -d '[:space:]')

no_downgrade=false
case "$value" in
  no-downgrade) no_downgrade=true ;;
esac

jq -cn \
  --argjson no_downgrade "$no_downgrade" \
  --arg raw "$value" \
  '{no_downgrade: $no_downgrade, raw: $raw}' | emit_ok PNP003
