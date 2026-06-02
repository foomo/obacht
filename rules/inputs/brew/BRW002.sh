#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Homebrew metadata freshness. age_days = -1 means unknown (no timestamp
# source available); the rego policy treats that as a data-conditional
# skip.

if ! command -v brew >/dev/null 2>&1; then
  emit_skip BRW002 "brew not installed"
  exit 0
fi

now=$(date +%s)
best=0
cache_dir=$(brew --cache 2>/dev/null || echo "")
repo_dir=$(brew --repository 2>/dev/null || echo "")
for f in "$cache_dir/api/formula.jws.json" "$cache_dir/api/cask.jws.json" "$repo_dir/.git/FETCH_HEAD"; do
  [ -n "$f" ] || continue
  [ -f "$f" ] || continue
  m=$(stat -f '%m' "$f" 2>/dev/null || stat -c '%Y' "$f" 2>/dev/null || echo 0)
  [ "$m" -gt "$best" ] && best=$m
done

age_days=-1
[ "$best" -gt 0 ] && age_days=$(( (now - best) / 86400 ))

jq -cn --argjson age_days "$age_days" '{age_days: $age_days}' | emit_ok BRW002
