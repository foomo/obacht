#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD002

v=$(claude_config_value '.autoCompactEnabled')
jq -cn --arg v "$v" '{auto_compact_enabled: $v}' | emit_ok CLD002
