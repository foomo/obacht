#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/claude.sh

claude_skip_if_missing CLD041

# Native Messaging manifests live under macOS-specific paths only.
case "$(uname -s)" in
  Darwin) ;;
  *)
    emit_skip CLD041 "Native Messaging check is macOS-only"
    exit 0
    ;;
esac

# JSON-escape a string for embedding inside a double-quoted JSON literal.
_cld041_escape() {
  printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}

manifests='[]'
manifest_search_dir="$HOME/Library/Application Support"
if [ -d "$manifest_search_dir" ]; then
  # Bounded depth — manifests live at <browser>/[<profile>/]NativeMessagingHosts/.
  manifest_files=$(find "$manifest_search_dir" -maxdepth 5 -type f \
    -name "com.anthropic.claude_browser_extension.json" 2>/dev/null)
  if [ -n "$manifest_files" ]; then
    saved_ifs=$IFS
    IFS='
'
    json='['
    first=1
    for mf in $manifest_files; do
      stat_out=$(stat -f "%z %Sf" "$mf" 2>/dev/null || printf '0 ')
      mf_size=${stat_out% *}
      mf_flags=${stat_out#* }
      # Mitigated = file ≤1 byte AND user-immutable (uchg) flag set.
      if [ "$mf_size" -le 1 ] && [ "$mf_flags" = uchg ]; then
        continue
      fi
      mf_escaped=$(_cld041_escape "$mf")
      if [ "$first" = 1 ]; then
        json="${json}\"${mf_escaped}\""
        first=0
      else
        json="${json},\"${mf_escaped}\""
      fi
    done
    IFS=$saved_ifs
    manifests="${json}]"
  fi
fi

printf '%s' "$manifests" \
  | jq -c '{claude_desktop_native_messaging_manifests: .}' \
  | emit_ok CLD041
