#!/bin/sh
config="$HOME/.kube/config"

# Resolve symlinks.
if [ -L "$config" ]; then
  config=$(readlink -f "$config" 2>/dev/null || echo "$config")
fi

if [ ! -f "$config" ]; then
  printf '{"config_exists": false, "config_mode": "", "contexts": []}'
  exit 0
fi

config_mode=$(stat -f '%04Lp' "$config" 2>/dev/null || stat -c '%04a' "$config" 2>/dev/null || echo "")

# Parse contexts from kubeconfig using awk (no yq dependency).
# Extract context name and cluster from YAML.
contexts="["
first=true
in_contexts=false
current_name=""
current_cluster=""

while IFS= read -r line; do
  # Detect the contexts section.
  case "$line" in
    "contexts:"*) in_contexts=true; continue ;;
  esac

  if [ "$in_contexts" = true ]; then
    # End of contexts section (new top-level key).
    case "$line" in
      [a-z]*:*) in_contexts=false; continue ;;
    esac

    case "$line" in
      *"- name: "*)
        # Flush previous context.
        if [ -n "$current_name" ]; then
          if [ "$first" = true ]; then first=false; else contexts="$contexts,"; fi
          contexts="$contexts{\"name\":\"$current_name\",\"cluster\":\"$current_cluster\"}"
        fi
        current_name=$(printf '%s' "$line" | sed 's/.*- name: *//' | tr -d '"')
        current_cluster=""
        ;;
      *"cluster: "*)
        current_cluster=$(printf '%s' "$line" | sed 's/.*cluster: *//' | tr -d '"')
        ;;
    esac
  fi
done < "$config"

# Flush last context.
if [ -n "$current_name" ]; then
  if [ "$first" = true ]; then first=false; else contexts="$contexts,"; fi
  contexts="$contexts{\"name\":\"$current_name\",\"cluster\":\"$current_cluster\"}"
fi

contexts="$contexts]"

printf '{"config_exists": true, "config_mode": "%s", "contexts": %s}' \
  "$config_mode" "$contexts"
