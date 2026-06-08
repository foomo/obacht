#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# mise github.credential_command sourcing.
#
# When set, mise invokes the command to authenticate GitHub API requests
# (release listings, asset downloads). "gh auth token" reuses the local gh
# CLI session, raising rate limits and avoiding unauthenticated traffic.

if ! command -v mise >/dev/null 2>&1; then
  emit_skip MIS004 "mise not installed"
  exit 0
fi

value=$(mise settings get github.credential_command 2>/dev/null \
  | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/^"//' -e 's/"$//')

jq -cn --arg value "$value" '{value: $value}' | emit_ok MIS004
