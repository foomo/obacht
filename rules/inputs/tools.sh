#!/bin/sh
tools="["
first=true

for tool in git opa gpg ssh-agent; do
  installed=false
  version=""
  path=""

  tool_path=$(command -v "$tool" 2>/dev/null || true)
  if [ -n "$tool_path" ]; then
    installed=true
    path="$tool_path"
    case "$tool" in
      git)       version=$(git --version 2>/dev/null | head -1) ;;
      opa)       version=$(opa version 2>/dev/null | head -1) ;;
      gpg)       version=$(gpg --version 2>/dev/null | head -1) ;;
      ssh-agent) version=$(ssh-agent -V 2>&1 | head -1) ;;
    esac
  fi

  if [ "$first" = true ]; then first=false; else tools="$tools,"; fi
  tools="$tools{\"name\":\"$tool\",\"installed\":$installed,\"version\":\"$version\",\"path\":\"$path\"}"
done

tools="$tools]"

# Homebrew auto-update check.
homebrew_installed=false
homebrew_auto_update_disabled=false
if command -v brew >/dev/null 2>&1; then
  homebrew_installed=true
  if [ -n "${HOMEBREW_NO_AUTO_UPDATE:-}" ]; then
    homebrew_auto_update_disabled=true
  fi
fi

printf '{"tools": %s, "homebrew_installed": %s, "homebrew_auto_update_disabled": %s}' \
  "$tools" "$homebrew_installed" "$homebrew_auto_update_disabled"
