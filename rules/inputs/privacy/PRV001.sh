#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# Password manager detection.
password_manager_installed=false
password_manager_name=""
for app in "1Password 7" "1Password" "Bitwarden" "Dashlane" "KeePassXC" "LastPass" "Enpass"; do
  if [ -d "/Applications/$app.app" ]; then
    password_manager_installed=true
    password_manager_name="$app"
    break
  fi
done

printf '{"password_manager_installed": %s, "password_manager_name": "%s"}' \
  "$password_manager_installed" "$password_manager_name" | emit_ok PRV001
