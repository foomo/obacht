#!/bin/sh
# shellcheck shell=sh
#
# Secret-handling helpers for obacht input scripts.
#
# is_literal_token <value>
#   Return 0 if <value> is a non-empty literal token (something a thief
#   could exfiltrate verbatim) and 1 if it is empty or an env-var
#   reference such as ${VAR}, $VAR, "${VAR}", or '${VAR}'. Surrounding
#   single/double quotes are stripped before classification.
#
# redact <value>
#   Print a masked form of <value> suitable for inclusion in evidence:
#   "***" for short values, otherwise first 4 + "***" + last 4.

is_literal_token() {
  v=$1
  v=${v#\"}
  v=${v%\"}
  v=${v#\'}
  v=${v%\'}
  [ -z "$v" ] && return 1
  # shellcheck disable=SC2016
  case "$v" in
    '${'*'}')  return 1 ;;
    '$'[A-Za-z_]*) return 1 ;;
  esac
  return 0
}

redact() {
  v=$1
  v=${v#\"}
  v=${v%\"}
  v=${v#\'}
  v=${v%\'}
  n=$(printf %s "$v" | wc -c | tr -d '[:space:]')
  if [ "$n" -le 8 ]; then
    printf '***'
    return
  fi
  head=$(printf %s "$v" | cut -c1-4)
  tail=$(printf %s "$v" | awk '{ print substr($0, length($0)-3) }')
  printf '%s***%s' "$head" "$tail"
}
