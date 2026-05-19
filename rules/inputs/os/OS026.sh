#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS026
  exit 0
fi

bluetooth_sharing_disabled=true
bt_sharing=$(defaults read com.apple.Bluetooth PrefKeyServicesEnabled 2>/dev/null || echo "0")
if [ "$bt_sharing" = "1" ]; then
  bluetooth_sharing_disabled=false
fi

jq -cn --arg os "$os" --argjson e "$bluetooth_sharing_disabled" '{os: $os, bluetooth_sharing_disabled: $e}' | emit_ok OS026
