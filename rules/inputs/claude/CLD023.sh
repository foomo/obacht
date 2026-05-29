#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD023

v=$(claude_settings_value '.env.IS_DEMO')
jq -cn --arg v "$v" '{env_is_demo: $v}' | emit_ok CLD023
