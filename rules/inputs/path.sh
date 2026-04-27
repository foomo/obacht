#!/bin/sh
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

  exists=false
  writable=false
  world_writable=false
  mode=""
  if [ -d "$entry" ]; then
    exists=true
    if [ -w "$entry" ]; then
      writable=true
    fi
    mode=$(stat -f '%Lp' "$entry" 2>/dev/null || stat -c '%a' "$entry" 2>/dev/null || echo "")
    case "$mode" in
      *2|*3|*6|*7) world_writable=true ;;
    esac
  fi

  if [ "$first" = true ]; then first=false; else dirs="$dirs,"; fi
  dirs="$dirs{\"path\":\"$entry\",\"exists\":$exists,\"writable\":$writable,\"world_writable\":$world_writable,\"mode\":\"$mode\",\"is_relative\":$is_relative}"
done

dirs="$dirs]"
printf '{"dirs": %s}' "$dirs"
