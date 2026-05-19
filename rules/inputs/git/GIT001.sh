#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

credential_helper=""
if command -v git >/dev/null 2>&1; then
  credential_helper=$(git config --global credential.helper 2>/dev/null || true)
fi

jq -cn --arg ch "$credential_helper" '{credential_helper: $ch}' | emit_ok GIT001
