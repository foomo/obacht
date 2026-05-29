#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD030

v=$(claude_settings_value '.plansDirectory')
jq -cn --arg v "$v" '{settings_plans_directory: $v}' | emit_ok CLD030
