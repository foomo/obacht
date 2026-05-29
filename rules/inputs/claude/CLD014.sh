#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD014

v=$(claude_settings_value '.env.DISABLE_FEEDBACK_COMMAND')
jq -cn --arg v "$v" '{env_disable_feedback_command: $v}' | emit_ok CLD014
