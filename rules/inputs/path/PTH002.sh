#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Collect each PATH entry and flag relative entries.
dirs="["
first=true

IFS=':'
for entry in $PATH; do
  [ -z "$entry" ] && continue

  is_relative=false
  case "$entry" in
    /*) ;;
    *) is_relative=true ;;
  esac

  if [ "$first" = true ]; then first=false; else dirs="$dirs,"; fi
  dirs="$dirs{\"path\":\"$entry\",\"is_relative\":$is_relative}"
done

dirs="$dirs]"

printf '{"dirs": %s}' "$dirs" | emit_ok PTH002
