#!/bin/sh
# shellcheck shell=sh
#
# Shared helpers for bumblebee per-rule input scripts.
#
# Depends on _lib/json.sh for emit_ok / emit_skip.
#
# Usage:
#   bumblebee_skip_if_unavailable BUM001
#   bumblebee_scan_ecosystem BUM001 npm

bumblebee_skip_if_unavailable() {
  rule_id="$1"

  if ! command -v bumblebee >/dev/null 2>&1; then
    emit_skip "$rule_id" "bumblebee not installed (go install github.com/perplexityai/bumblebee/cmd/bumblebee@latest)"
    exit 0
  fi

  if [ -z "${OBACHT_BUMBLEBEE_CATALOG_DIR:-}" ] || [ ! -d "$OBACHT_BUMBLEBEE_CATALOG_DIR" ]; then
    emit_skip "$rule_id" "no exposure catalog available"
    exit 0
  fi
}

bumblebee_scan_ecosystem() {
  rule_id="$1"
  ecosystem="$2"

  out=$(bumblebee scan \
    --profile baseline \
    --exposure-catalog "$OBACHT_BUMBLEBEE_CATALOG_DIR" \
    --ecosystem "$ecosystem" \
    --findings-only 2>/dev/null) || {
    emit_skip "$rule_id" "bumblebee scan failed"
    exit 0
  }

  printf '%s\n' "$out" \
    | jq -s '{findings: map(select(.record_type=="finding"))}' \
    | emit_ok "$rule_id"
}
