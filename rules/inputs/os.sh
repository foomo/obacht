#!/bin/sh
os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
hostname=$(hostname 2>/dev/null || echo "")

# Default all booleans for non-macOS.
if [ "$os" != "darwin" ]; then
  printf '{"os":"%s","arch":"%s","hostname":"%s"}' "$os" "$arch" "$hostname"
  exit 0
fi

# macOS-specific checks.
cmd_contains() {
  output=$("$@" 2>&1) || return 1
  return 0
}

bool_check() {
  output=$($1 ${2:+"$2"} ${3:+"$3"} ${4:+"$4"} ${5:+"$5"} 2>&1)
  if [ $? -ne 0 ]; then echo "false"; return; fi
  if [ -z "$6" ]; then echo "true"; return; fi
  if printf '%s' "$output" | grep -q "$6"; then echo "true"; else echo "false"; fi
}

sip=$(bool_check csrutil status "" "" "" "" "enabled")
filevault=$(bool_check fdesetup status "" "" "" "" "On")
firewall=$(bool_check /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate "" "" "" "enabled")
stealth=$(bool_check /usr/libexec/ApplicationFirewall/socketfilterfw --getstealthmode "" "" "" "enabled")
gatekeeper=$(bool_check spctl --status "" "" "" "enabled")

# AutoLogin: if defaults read succeeds, auto-login IS enabled (bad).
auto_login_disabled=true
if defaults read /Library/Preferences/com.apple.loginwindow autoLoginUser >/dev/null 2>&1; then
  auto_login_disabled=false
fi

guest_disabled=$(bool_check defaults read /Library/Preferences/com.apple.loginwindow GuestEnabled "" "" "0")

# Screen lock timeout.
timeout_raw=$(defaults -currentHost read com.apple.screensaver idleTime 2>/dev/null || echo "0")
screen_lock_timeout=$(printf '%s' "$timeout_raw" | tr -dc '0-9')
[ -z "$screen_lock_timeout" ] && screen_lock_timeout=0

os_auto_update=$(bool_check defaults read /Library/Preferences/com.apple.SoftwareUpdate AutomaticallyInstallMacOSUpdates "" "" "1")
app_auto_update=$(bool_check defaults read /Library/Preferences/com.apple.commerce AutoUpdate "" "" "1")
rsr=$(bool_check defaults read /Library/Preferences/com.apple.SoftwareUpdate ConfigDataInstall "" "" "1")

# Screen sharing: if launchctl list succeeds, it's enabled (bad).
screen_sharing_disabled=true
if launchctl list com.apple.screensharing >/dev/null 2>&1; then
  screen_sharing_disabled=false
fi

# Internet sharing.
internet_sharing_disabled=true
if defaults read /Library/Preferences/SystemConfiguration/com.apple.nat Enabled 2>/dev/null | grep -q 1; then
  internet_sharing_disabled=false
fi

printer_sharing_disabled=$(bool_check cupsctl "" "" "" "" "_share_printers=0")

# Remote Apple Events.
remote_apple_events_disabled=true
if launchctl list com.apple.AEServer >/dev/null 2>&1; then
  remote_apple_events_disabled=false
fi

# AirDrop setting.
airdrop_raw=$(defaults read com.apple.sharingd DiscoverableMode 2>/dev/null || echo "Off")
case "$airdrop_raw" in
  Everyone)       airdrop="everyone" ;;
  "Contacts Only") airdrop="contacts_only" ;;
  *)              airdrop="off" ;;
esac

# Rosetta.
rosetta=false
if pgrep -q oahd 2>/dev/null; then rosetta=true; fi

# EDR detection.
edr=false
for agent in com.crowdstrike.falcon com.sentinelone com.carbon.black com.microsoft.wdav; do
  if launchctl list "$agent" >/dev/null 2>&1; then edr=true; break; fi
done

# Legacy kexts.
legacy_kexts_blocked=true
if kmutil showloaded --list-only 2>/dev/null | grep -q "com.apple"; then
  legacy_kexts_blocked=false
fi

# MDM.
mdm=false
if profiles status -type enrollment 2>/dev/null | grep -q "MDM enrollment: Yes"; then
  mdm=true
fi

cat <<EOF
{
  "os": "$os",
  "arch": "$arch",
  "hostname": "$hostname",
  "sip_enabled": $sip,
  "filevault_enabled": $filevault,
  "firewall_enabled": $firewall,
  "stealth_mode_enabled": $stealth,
  "gatekeeper_enabled": $gatekeeper,
  "auto_login_disabled": $auto_login_disabled,
  "guest_account_disabled": $guest_disabled,
  "screen_lock_timeout_seconds": $screen_lock_timeout,
  "os_auto_update_enabled": $os_auto_update,
  "app_auto_update_enabled": $app_auto_update,
  "rsr_enabled": $rsr,
  "screen_sharing_disabled": $screen_sharing_disabled,
  "internet_sharing_disabled": $internet_sharing_disabled,
  "printer_sharing_disabled": $printer_sharing_disabled,
  "remote_apple_events_disabled": $remote_apple_events_disabled,
  "airdrop_setting": "$airdrop",
  "rosetta_installed": $rosetta,
  "edr_deployed": $edr,
  "legacy_kexts_blocked": $legacy_kexts_blocked,
  "mdm_enrolled": $mdm
}
EOF
