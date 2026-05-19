#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

installed=false
gitignore_file=""
if command -v git >/dev/null 2>&1; then
  installed=true
  excludes_file=$(git config --global core.excludesfile 2>/dev/null || true)
  if [ -n "$excludes_file" ]; then
    excludes_file=$(eval echo "$excludes_file")
    if [ -f "$excludes_file" ]; then
      gitignore_file="$excludes_file"
    fi
  fi
fi

keys_patterns='id_rsa id_ed25519 id_ecdsa *.pem *.key *.p12 *.pfx *.cer *.crt *.jks *.p8 *.der *.jce *.keystore *.truststore *.p7b *.p7s *.p10 *.csr *.req *.p7c *.p7m *.p7r'

missing="["
first=true
oIFS="$IFS"
IFS=' '
for p in $keys_patterns; do
  if [ -z "$gitignore_file" ] || ! grep -Fxq "$p" "$gitignore_file" 2>/dev/null; then
    if [ "$first" = true ]; then first=false; else missing="$missing,"; fi
    missing="$missing\"$p\""
  fi
done
IFS="$oIFS"
missing="$missing]"

jq -cn \
  --argjson installed "$installed" \
  --argjson missing "$missing" \
  '{installed: $installed, gitignore_missing_keys: $missing}' | emit_ok GIT005
