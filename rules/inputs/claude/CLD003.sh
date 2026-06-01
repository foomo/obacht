#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD003

v=$(claude_config_value '.prStatusFooterEnabled')
jq -cn --arg v "$v" '{pr_status_footer_enabled: $v}' | emit_ok CLD003
