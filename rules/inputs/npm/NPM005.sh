#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/secrets.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/secrets.sh

# Detect plaintext registry auth tokens in the user .npmrc.
#
# npm stores auth credentials in INI-style `.npmrc` files. Literal
# values for `_authToken`, `_auth`, `_password`, and `_authIdent`
# are recoverable by any process that can read the file. Env-var
# references such as `${NPM_TOKEN}` are treated as safe indirection
# and pass.

if ! command -v npm >/dev/null 2>&1; then
  emit_skip NPM005 "npm not installed"
  exit 0
fi

cfg="${NPM_CONFIG_USERCONFIG:-$HOME/.npmrc}"

if [ ! -f "$cfg" ]; then
  jq -cn --arg path "$cfg" '{config_path: $path, tokens: []}' | emit_ok NPM005
  exit 0
fi

if [ ! -r "$cfg" ]; then
  emit_skip NPM005 "config not readable: $cfg"
  exit 0
fi

tokens='[]'
# Match lines like:
#   _authToken=value
#   _auth=value
#   _password=value
#   _authIdent=value
#   //registry.example.com/:_authToken=value
# Skip comment lines starting with # or ;.
while IFS= read -r entry; do
  line_no=${entry%%:*}
  rest=${entry#*:}
  # Strip leading whitespace.
  while :; do
    case "$rest" in
      ' '*|"	"*) rest=${rest#?} ;;
      *) break ;;
    esac
  done
  case "$rest" in
    '#'*|';'*) continue ;;
  esac
  case "$rest" in
    *=*) ;;
    *) continue ;;
  esac
  key_part=${rest%%=*}
  value=${rest#*=}
  # Strip trailing CR / whitespace from value.
  cr=$(printf '\r')
  value=${value%"$cr"}
  while :; do
    case "$value" in
      *' '|*"	") value=${value%?} ;;
      *) break ;;
    esac
  done
  # Determine the bare key (after any `//host/:` prefix).
  case "$key_part" in
    *:*) bare=${key_part##*:} ;;
    *)   bare=$key_part ;;
  esac
  case "$bare" in
    _authToken|_auth|_password|_authIdent) ;;
    *) continue ;;
  esac
  if is_literal_token "$value"; then
    masked=$(redact "$value")
    tokens=$(printf '%s' "$tokens" | jq -c \
      --arg key "$bare" \
      --argjson line "$line_no" \
      --arg masked "$masked" \
      '. + [{key: $key, line: $line, masked: $masked}]')
  fi
done <<EOF
$(grep -n -E '(^|[^[:alnum:]_])(_authToken|_auth|_password|_authIdent)[[:space:]]*=' "$cfg" 2>/dev/null)
EOF

jq -cn \
  --arg path "$cfg" \
  --argjson tokens "$tokens" \
  '{config_path: $path, tokens: $tokens}' | emit_ok NPM005
