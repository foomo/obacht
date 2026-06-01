#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD032

v=$(claude_settings_value '.sandbox.autoAllowBashIfSandboxed')
jq -cn --arg v "$v" '{sandbox_auto_allow_bash_if_sandboxed: $v}' | emit_ok CLD032
