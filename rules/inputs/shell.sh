#!/bin/sh
shell="$SHELL"
home="$HOME"

# Determine history file based on shell.
hist_file=""
case "$(basename "$shell")" in
  bash) hist_file="$home/.bash_history" ;;
  zsh)  hist_file="$home/.zsh_history" ;;
  fish) hist_file="$home/.local/share/fish/fish_history" ;;
esac

hist_mode=""
if [ -n "$hist_file" ] && [ -f "$hist_file" ]; then
  # Resolve symlinks.
  real_hist=$(readlink -f "$hist_file" 2>/dev/null || echo "$hist_file")
  hist_mode=$(stat -f '%04Lp' "$real_hist" 2>/dev/null || stat -c '%04a' "$real_hist" 2>/dev/null || echo "")
  hist_mode="0$hist_mode"
fi

histcontrol="${HISTCONTROL:-}"

printf '{"shell": "%s", "history_file": "%s", "history_file_mode": "%s", "histcontrol": "%s"}' \
  "$shell" "$hist_file" "$hist_mode" "$histcontrol"
