#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD026

# Empty-string is the hardened state; "unset" means key missing.
commit=$(claude_settings_value '.attribution.commit')
pr=$(claude_settings_value '.attribution.pr')
jq -cn --arg c "$commit" --arg p "$pr" \
  '{settings_attribution_commit: $c, settings_attribution_pr: $p}' | emit_ok CLD026
