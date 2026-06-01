#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD015

v=$(claude_settings_value '.env.DISABLE_EXTRA_USAGE_COMMAND')
jq -cn --arg v "$v" '{env_disable_extra_usage_command: $v}' | emit_ok CLD015
