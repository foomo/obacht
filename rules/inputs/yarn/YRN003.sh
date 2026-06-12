#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/secrets.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/secrets.sh

# Detect plaintext registry auth tokens in the user .yarnrc.yml (yarn berry).
#
# Yarn berry stores `npmAuthToken` / `npmAuthIdent` either at the top
# level or nested under `npmRegistries:` / `npmScopes:` blocks in
# `~/.yarnrc.yml`. Literal values left on disk are recoverable by any
# process that can read the file. Env-var references such as
# `${NPM_TOKEN}` are treated as safe indirection and pass. Yarn
# classic (1.x) does not use this file format and is skipped.

if ! command -v yarn >/dev/null 2>&1; then
  emit_skip YRN003 "yarn not installed"
  exit 0
fi

ver=$(yarn --version 2>/dev/null | tr -d '[:space:]')
case "$ver" in
  1.*|0.*)
    emit_skip YRN003 "yarn classic ($ver) does not use .yarnrc.yml"
    exit 0
    ;;
esac

cfg="$HOME/.yarnrc.yml"

if [ ! -f "$cfg" ]; then
  jq -cn --arg path "$cfg" '{config_path: $path, tokens: []}' | emit_ok YRN003
  exit 0
fi

if [ ! -r "$cfg" ]; then
  emit_skip YRN003 "config not readable: $cfg"
  exit 0
fi

tokens='[]'
while IFS= read -r entry; do
  line_no=${entry%%:*}
  rest=${entry#*:}
  while :; do
    case "$rest" in
      ' '*|"	"*) rest=${rest#?} ;;
      *) break ;;
    esac
  done
  case "$rest" in
    '#'*) continue ;;
  esac
  case "$rest" in
    npmAuthToken:*) key=npmAuthToken; value=${rest#npmAuthToken:} ;;
    npmAuthIdent:*) key=npmAuthIdent; value=${rest#npmAuthIdent:} ;;
    *) continue ;;
  esac
  # Strip leading whitespace from value.
  while :; do
    case "$value" in
      ' '*|"	"*) value=${value#?} ;;
      *) break ;;
    esac
  done
  # Drop inline comment.
  value=${value%%#*}
  # Strip trailing whitespace + CR.
  cr=$(printf '\r')
  value=${value%"$cr"}
  while :; do
    case "$value" in
      *' '|*"	") value=${value%?} ;;
      *) break ;;
    esac
  done
  if is_literal_token "$value"; then
    masked=$(redact "$value")
    tokens=$(printf '%s' "$tokens" | jq -c \
      --arg key "$key" \
      --argjson line "$line_no" \
      --arg masked "$masked" \
      '. + [{key: $key, line: $line, masked: $masked}]')
  fi
done <<EOF
$(grep -n -E '^[[:space:]]*(npmAuthToken|npmAuthIdent):' "$cfg" 2>/dev/null)
EOF

jq -cn \
  --arg path "$cfg" \
  --argjson tokens "$tokens" \
  '{config_path: $path, tokens: $tokens}' | emit_ok YRN003
