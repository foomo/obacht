#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS022
  exit 0
fi

airdrop_raw=$(defaults read com.apple.sharingd DiscoverableMode 2>/dev/null || echo "Off")
case "$airdrop_raw" in
  Everyone)        airdrop="everyone" ;;
  "Contacts Only") airdrop="contacts_only" ;;
  *)               airdrop="off" ;;
esac

jq -cn --arg os "$os" --arg a "$airdrop" '{os: $os, airdrop_setting: $a}' | emit_ok OS022
