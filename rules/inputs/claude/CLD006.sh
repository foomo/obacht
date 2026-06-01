#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD006

v=$(claude_settings_value '.env.DISABLE_COMPACT')
jq -cn --arg v "$v" '{env_disable_compact: $v}' | emit_ok CLD006
