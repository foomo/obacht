#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD016

v=$(claude_settings_value '.env.CLAUDE_CODE_DISABLE_FAST_MODE')
jq -cn --arg v "$v" '{env_claude_code_disable_fast_mode: $v}' | emit_ok CLD016
