#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD022

v=$(claude_settings_value '.env.FORCE_AUTOUPDATE_PLUGINS')
jq -cn --arg v "$v" '{env_force_autoupdate_plugins: $v}' | emit_ok CLD022
