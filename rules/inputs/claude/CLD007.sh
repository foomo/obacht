#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD007

v=$(claude_settings_value '.env.DISABLE_TELEMETRY')
jq -cn --arg v "$v" '{env_disable_telemetry: $v}' | emit_ok CLD007
