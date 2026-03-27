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
  if [ -d "$entry" ]; then
    exists=true
    # Test writability by attempting to create a temp file.
    tmpfile=$(mktemp "$entry/.bouncer-check-XXXXXX" 2>/dev/null) && {
      writable=true
      rm -f "$tmpfile"
    }
  fi

  if [ "$first" = true ]; then first=false; else dirs="$dirs,"; fi
  dirs="$dirs{\"path\":\"$entry\",\"exists\":$exists,\"writable\":$writable,\"is_relative\":$is_relative}"
done

dirs="$dirs]"
printf '{"dirs": %s}' "$dirs"
