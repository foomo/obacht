#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD024

v=$(claude_settings_value '.disableAutoMode')
jq -cn --arg v "$v" '{settings_disable_auto_mode: $v}' | emit_ok CLD024
