#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS018
  exit 0
fi

edr_deployed=false
for agent in com.crowdstrike.falcon com.sentinelone com.carbon.black com.microsoft.wdav; do
  if launchctl list "$agent" >/dev/null 2>&1; then
    edr_deployed=true
    break
  fi
done

jq -cn --arg os "$os" --argjson e "$edr_deployed" '{os: $os, edr_deployed: $e}' | emit_ok OS018
