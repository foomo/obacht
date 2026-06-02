#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# npm registry transport security.
#
# The default registry should be reached over HTTPS so package
# tarballs, version metadata, and auth tokens are not exposed to
# network-level tampering or interception. An explicit downgrade to
# `http://` in .npmrc disables TLS for every install.

if ! command -v npm >/dev/null 2>&1; then
  emit_skip NPM004 "npm not installed"
  exit 0
fi

value=$(npm config get registry 2>/dev/null | tr -d '[:space:]')

secure=false
case "$value" in
  https://*) secure=true ;;
esac

jq -cn \
  --argjson secure "$secure" \
  --arg raw "$value" \
  '{secure: $secure, raw: $raw}' | emit_ok NPM004
