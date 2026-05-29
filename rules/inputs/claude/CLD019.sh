#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD019

v=$(claude_settings_value '.env.CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY')
jq -cn --arg v "$v" '{env_claude_code_disable_feedback_survey: $v}' | emit_ok CLD019
