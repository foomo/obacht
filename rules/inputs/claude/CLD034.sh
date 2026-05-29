#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD034

v=$(claude_settings_value '.sandbox.network.allowManagedDomainsOnly')
jq -cn --arg v "$v" '{sandbox_network_allow_managed_domains_only: $v}' | emit_ok CLD034
