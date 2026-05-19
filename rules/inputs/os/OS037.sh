#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS037
  exit 0
fi

in_app_review_disabled=false
value=$(defaults read com.apple.appstore InAppReviewEnabled 2>/dev/null)
if [ "$value" = "0" ]; then
  in_app_review_disabled=true
fi

jq -cn --arg os "$os" --argjson d "$in_app_review_disabled" '{os: $os, in_app_review_disabled: $d}' | emit_ok OS037
