#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD028

v=$(claude_settings_value '.skipWebFetchPreflight')
jq -cn --arg v "$v" '{settings_skip_web_fetch_preflight: $v}' | emit_ok CLD028
