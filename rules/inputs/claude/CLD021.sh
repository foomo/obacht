#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD021

v=$(claude_settings_value '.env.CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS')
jq -cn --arg v "$v" '{env_claude_code_disable_experimental_betas: $v}' | emit_ok CLD021
