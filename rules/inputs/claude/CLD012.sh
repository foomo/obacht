#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD012

v=$(claude_settings_value '.env.DISABLE_ERROR_REPORTING')
jq -cn --arg v "$v" '{env_disable_error_reporting: $v}' | emit_ok CLD012
