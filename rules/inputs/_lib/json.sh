#!/bin/sh
# shellcheck shell=sh
#
# Standard JSON envelope helpers for obacht input scripts.
#
# Usage:
#   collect_data | emit_ok SSH001
#   emit_skip SSH001 "no .ssh directory"

emit_ok() {
  rule_id="$1"
  jq -c --arg id "$rule_id" '{rule_id: $id, status: "ok", data: .}'
}

emit_skip() {
  rule_id="$1"
  reason="$2"
  jq -cn --arg id "$rule_id" --arg r "$reason" \
    '{rule_id: $id, status: "skip", skip_reason: $r}'
}
