#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD025

v=$(claude_settings_value '.disableDeepLinkRegistration')
jq -cn --arg v "$v" '{settings_disable_deep_link_registration: $v}' | emit_ok CLD025
