#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# mise status.missing_tools verbosity.
#
# Controls when mise warns about tools listed in config but not installed.
# "always" surfaces drift on every shell activation; the default
# ("if_other_versions_installed") can hide a missing tool entirely.

if ! command -v mise >/dev/null 2>&1; then
  emit_skip MIS003 "mise not installed"
  exit 0
fi

value=$(mise settings get status.missing_tools 2>/dev/null | tr -d '[:space:]"')

jq -cn --arg value "$value" '{value: $value}' | emit_ok MIS003
