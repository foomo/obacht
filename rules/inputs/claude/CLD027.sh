#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD027

v=$(claude_settings_value '.respectGitignore')
jq -cn --arg v "$v" '{settings_respect_gitignore: $v}' | emit_ok CLD027
