#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Collect each PATH entry and flag world-writable directories.
dirs="["
first=true

IFS=':'
for entry in $PATH; do
  [ -z "$entry" ] && continue

  exists=false
  world_writable=false
  mode=""
  if [ -d "$entry" ]; then
    exists=true
    mode=$(stat -f '%Lp' "$entry" 2>/dev/null || stat -c '%a' "$entry" 2>/dev/null || echo "")
    case "$mode" in
      *2|*3|*6|*7) world_writable=true ;;
    esac
  fi

  if [ "$first" = true ]; then first=false; else dirs="$dirs,"; fi
  dirs="$dirs{\"path\":\"$entry\",\"exists\":$exists,\"world_writable\":$world_writable,\"mode\":\"$mode\"}"
done

dirs="$dirs]"

printf '{"dirs": %s}' "$dirs" | emit_ok PTH001
