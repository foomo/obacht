#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD035

v=$(claude_settings_value '.permissions.disableBypassPermissionsMode')
jq -cn --arg v "$v" '{permissions_disable_bypass_mode: $v}' | emit_ok CLD035
