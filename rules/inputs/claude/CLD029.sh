#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD029

v=$(claude_settings_value '.autoMemoryDirectory')
jq -cn --arg v "$v" '{settings_auto_memory_directory: $v}' | emit_ok CLD029
