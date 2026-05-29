#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD004

v=$(claude_config_value '.claudeInChromeDefaultEnabled')
jq -cn --arg v "$v" '{claude_in_chrome_default_enabled: $v}' | emit_ok CLD004
