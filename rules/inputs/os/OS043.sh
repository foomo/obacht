#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

os=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$os" != "darwin" ]; then
  jq -cn --arg os "$os" '{os: $os}' | emit_ok OS043
  exit 0
fi

plist=/Library/Preferences/SystemConfiguration/com.apple.wifi.known-networks.plist
legacy=/Library/Preferences/SystemConfiguration/com.apple.airport.preferences.plist

src=""
if [ -r "$plist" ]; then
  src="$plist"
elif [ -r "$legacy" ]; then
  src="$legacy"
fi

if [ -z "$src" ]; then
  emit_skip OS043 "Wi-Fi preferences plist not readable"
  exit 0
fi

json=$(plutil -convert json -o - "$src" 2>/dev/null)
if [ -z "$json" ]; then
  emit_skip OS043 "could not parse Wi-Fi preferences plist"
  exit 0
fi

# Two known shapes: top-level dict of networks (modern) or nested KnownNetworks (legacy)
counts=$(printf '%s' "$json" | jq -c '
  def nets:
    if has("KnownNetworks") then .KnownNetworks
    else with_entries(select(.value | type == "object" and (has("PrivateMACAddressModeUserSetting") or has("PrivateMACAddressModeSystemSetting") or has("SSID"))))
    end;
  (nets | to_entries) as $e |
  {
    total: ($e | length),
    disabled: ([$e[] | select(.value.PrivateMACAddressModeUserSetting == 0 or .value.PrivateMACAddressModeSystemSetting == 0)] | length)
  }
')

total=$(printf '%s' "$counts" | jq -r '.total')
disabled=$(printf '%s' "$counts" | jq -r '.disabled')

jq -cn --arg os "$os" --argjson t "$total" --argjson d "$disabled" \
  '{os: $os, known_networks: $t, mac_randomization_disabled_count: $d}' | emit_ok OS043
