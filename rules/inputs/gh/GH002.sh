#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

if ! command -v gh >/dev/null 2>&1; then
  emit_skip GH002 "gh not installed"
  exit 0
fi

if ! command -v git >/dev/null 2>&1; then
  emit_skip GH002 "git not installed"
  exit 0
fi

github_helpers=$(git config --global --get-all credential.https://github.com.helper 2>/dev/null || true)
gist_helpers=$(git config --global --get-all credential.https://gist.github.com.helper 2>/dev/null || true)

jq -cn \
  --arg gh "$github_helpers" \
  --arg gi "$gist_helpers" \
  '{github_helpers: $gh, gist_helpers: $gi}' \
  | emit_ok GH002
