#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Homebrew auto-update check.
homebrew_installed=false
homebrew_auto_update_disabled=false
if command -v brew >/dev/null 2>&1; then
  homebrew_installed=true
  if [ -n "${HOMEBREW_NO_AUTO_UPDATE:-}" ]; then
    homebrew_auto_update_disabled=true
  fi
fi

printf '{"homebrew_installed": %s, "homebrew_auto_update_disabled": %s}' \
  "$homebrew_installed" "$homebrew_auto_update_disabled" | emit_ok TOL002
