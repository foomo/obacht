#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# APT metadata freshness. age_days = -1 means unknown (no timestamp
# source available); the rego policy treats that as a data-conditional
# skip.

if ! command -v apt-get >/dev/null 2>&1; then
  emit_skip APT001 "apt-get not installed"
  exit 0
fi

now=$(date +%s)
best=0
for f in /var/lib/apt/periodic/update-success-stamp /var/cache/apt/pkgcache.bin; do
  [ -f "$f" ] || continue
  m=$(stat -f '%m' "$f" 2>/dev/null || stat -c '%Y' "$f" 2>/dev/null || echo 0)
  [ "$m" -gt "$best" ] && best=$m
done

age_days=-1
[ "$best" -gt 0 ] && age_days=$(( (now - best) / 86400 ))

jq -cn --argjson age_days "$age_days" '{age_days: $age_days}' | emit_ok APT001
