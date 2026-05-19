#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Package manager metadata freshness (brew, apt). age_days = -1 means
# unknown (no timestamp source available); engine policy treats that as skip.
now=$(date +%s)
package_managers="["
pm_first=true

if command -v brew >/dev/null 2>&1; then
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
  if [ "$pm_first" = true ]; then pm_first=false; else package_managers="$package_managers,"; fi
  package_managers="$package_managers{\"name\":\"brew\",\"age_days\":$age_days}"
fi

if command -v apt-get >/dev/null 2>&1; then
  best=0
  for f in /var/lib/apt/periodic/update-success-stamp /var/cache/apt/pkgcache.bin; do
    [ -f "$f" ] || continue
    m=$(stat -f '%m' "$f" 2>/dev/null || stat -c '%Y' "$f" 2>/dev/null || echo 0)
    [ "$m" -gt "$best" ] && best=$m
  done
  age_days=-1
  [ "$best" -gt 0 ] && age_days=$(( (now - best) / 86400 ))
  if [ "$pm_first" = true ]; then pm_first=false; else package_managers="$package_managers,"; fi
  package_managers="$package_managers{\"name\":\"apt\",\"age_days\":$age_days}"
fi

package_managers="$package_managers]"

printf '{"package_managers": %s}' "$package_managers" | emit_ok TOL003
