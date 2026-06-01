#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD011

v=$(claude_settings_value '.env.DISABLE_LOGOUT_COMMAND')
jq -cn --arg v "$v" '{env_disable_logout_command: $v}' | emit_ok CLD011
