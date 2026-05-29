#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD020

v=$(claude_settings_value '.env.CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING')
jq -cn --arg v "$v" '{env_claude_code_disable_file_checkpointing: $v}' | emit_ok CLD020
