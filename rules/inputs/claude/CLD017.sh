#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD017

v=$(claude_settings_value '.env.DISABLE_INSTALL_GITHUB_APP_COMMAND')
jq -cn --arg v "$v" '{env_disable_install_github_app_command: $v}' | emit_ok CLD017
