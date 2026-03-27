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
printf '{"tools": %s}' "$tools"
