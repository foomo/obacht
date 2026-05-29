#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD005

v=$(claude_config_value '.sandbox.failIfUnavailable')
jq -cn --arg v "$v" '{sandbox_fail_if_unavailable: $v}' | emit_ok CLD005
