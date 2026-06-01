#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD031

v=$(claude_settings_value '.sandbox.enabled')
jq -cn --arg v "$v" '{sandbox_enabled: $v}' | emit_ok CLD031
