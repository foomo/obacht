#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/secrets.sh
# shellcheck source=../_lib/json.sh
# shellcheck source=../_lib/secrets.sh

# Detect plaintext registry auth tokens in the global bunfig.toml.
#
# Bun stores registry credentials in `bunfig.toml`. Token-bearing keys
# (`token`, `username`, `password`) appear under `[install.registry]`,
# `[install.scopes."<scope>"]`, or as inline tables under
# `[install.scopes]`. Literal values left on disk are recoverable by
# any process that can read the file. Env-var references such as
# `$NPM_TOKEN` are treated as safe indirection and pass. Bunfig
# resolution follows BUN_CONFIG_FILE, XDG_CONFIG_HOME, then HOME.

if ! command -v bun >/dev/null 2>&1; then
  emit_skip BUN003 "bun not installed"
  exit 0
fi

resolve_bunfig() {
  if [ -n "${BUN_CONFIG_FILE:-}" ] && [ -f "$BUN_CONFIG_FILE" ]; then
    printf '%s' "$BUN_CONFIG_FILE"
    return
  fi
  if [ -n "${XDG_CONFIG_HOME:-}" ] && [ -f "$XDG_CONFIG_HOME/.bunfig.toml" ]; then
    printf '%s' "$XDG_CONFIG_HOME/.bunfig.toml"
    return
  fi
  if [ -f "$HOME/.config/.bunfig.toml" ]; then
    printf '%s' "$HOME/.config/.bunfig.toml"
    return
  fi
  if [ -f "$HOME/.bunfig.toml" ]; then
    printf '%s' "$HOME/.bunfig.toml"
  fi
}

bunfig=$(resolve_bunfig)

if [ -z "$bunfig" ]; then
  jq -cn --arg path "" '{config_path: $path, tokens: []}' | emit_ok BUN003
  exit 0
fi

if [ ! -r "$bunfig" ]; then
  emit_skip BUN003 "config not readable: $bunfig"
  exit 0
fi

# Extract candidate {line:key:value} tuples via awk. Output one record
# per line: "<line>\t<key>\t<value>" (value preserves quoting for the
# shell-side literal classifier).
raw=$(awk '
  function trim(s) {
    sub(/^[[:space:]]+/, "", s); sub(/[[:space:]]+$/, "", s); return s
  }
  function strip_comment(s,   i,c,inq,q) {
    inq = 0
    for (i = 1; i <= length(s); i++) {
      c = substr(s, i, 1)
      if (inq) {
        if (c == q) inq = 0
        continue
      }
      if (c == "\"" || c == "'\''") { inq = 1; q = c; continue }
      if (c == "#") return substr(s, 1, i-1)
    }
    return s
  }
  function emit(key, value,   v) {
    v = trim(value)
    if (v == "") return
    print NR "\t" key "\t" v
  }
  function scan_inline(line,   body, n, parts, i, kv, k, v) {
    if (match(line, /\{[^}]*\}/) == 0) return
    body = substr(line, RSTART + 1, RLENGTH - 2)
    n = split(body, parts, ",")
    for (i = 1; i <= n; i++) {
      kv = parts[i]
      if (match(kv, /=/) == 0) continue
      k = trim(substr(kv, 1, RSTART - 1))
      v = trim(substr(kv, RSTART + 1))
      if (k == "token" || k == "username" || k == "password") emit(k, v)
    }
  }

  {
    line = strip_comment($0)
    s = trim(line)
    if (s == "") next
    if (substr(s, 1, 1) == "[") {
      section = s
      next
    }
    # Only consider sections under install.registry or install.scopes.
    if (section !~ /^\[install\.registry/ && section !~ /^\[install\.scopes/) next
    # Inline-table assignments under [install.scopes].
    if (index(s, "{") > 0 && index(s, "}") > 0) {
      scan_inline(s)
      next
    }
    # Plain key = value.
    if (match(s, /=/) == 0) next
    k = trim(substr(s, 1, RSTART - 1))
    v = trim(substr(s, RSTART + 1))
    # Strip a possible surrounding quoted-string key indicator (e.g. "token").
    gsub(/^["'\'']|["'\'']$/, "", k)
    if (k == "token" || k == "username" || k == "password") emit(k, v)
  }
' "$bunfig" 2>/dev/null)

tokens='[]'
if [ -n "$raw" ]; then
  while IFS=$(printf '\t') read -r line_no key value; do
    [ -z "$key" ] && continue
    if is_literal_token "$value"; then
      masked=$(redact "$value")
      tokens=$(printf '%s' "$tokens" | jq -c \
        --arg key "$key" \
        --argjson line "$line_no" \
        --arg masked "$masked" \
        '. + [{key: $key, line: $line, masked: $masked}]')
    fi
  done <<EOF
$raw
EOF
fi

jq -cn \
  --arg path "$bunfig" \
  --argjson tokens "$tokens" \
  '{config_path: $path, tokens: $tokens}' | emit_ok BUN003
