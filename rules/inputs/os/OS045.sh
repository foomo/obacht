#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS045
  exit 0
fi

save_to_icloud=false
value=$(defaults read NSGlobalDomain NSDocumentSaveNewDocumentsToCloud 2>/dev/null)
# Default behaviour when key is absent is to save to iCloud (= true)
if [ -z "$value" ] || [ "$value" = "1" ]; then
  save_to_icloud=true
fi

jq -cn --arg os "$os" --argjson e "$save_to_icloud" '{os: $os, save_to_icloud: $e}' | emit_ok OS045
