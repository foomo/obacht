#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD033

v=$(claude_settings_value '.sandbox.allowUnsandboxedCommands')
jq -cn --arg v "$v" '{sandbox_allow_unsandboxed_commands: $v}' | emit_ok CLD033
