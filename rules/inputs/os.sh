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
stealth=$(bool_check /usr/libexec/ApplicationFirewall/socketfilterfw --getstealthmode "" "" "" "on")
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

# Time Machine backup.
timemachine_enabled=false
if tmutil destinationinfo 2>/dev/null | grep -q "Name"; then
  if defaults read /Library/Preferences/com.apple.TimeMachine AutoBackup 2>/dev/null | grep -q "1"; then
    timemachine_enabled=true
  fi
fi

# Remote Login (SSH server).
remote_login_disabled=true
if systemsetup -getremotelogin 2>/dev/null | grep -qi "on"; then
  remote_login_disabled=false
fi

# Remote Management. The ARD agent may be loaded with OnDemand=true (registered
# but not actively running). Only treat it as enabled when an actual PID is
# attached — `launchctl list <label>` returns success for any loaded job.
remote_management_disabled=true
ard_pid=$(launchctl list 2>/dev/null | awk '$3 == "com.apple.RemoteDesktop.agent" {print $1; exit}')
if [ -n "$ard_pid" ] && [ "$ard_pid" != "-" ]; then
  remote_management_disabled=false
fi

# Bluetooth Sharing.
bluetooth_sharing_disabled=true
bt_sharing=$(defaults read com.apple.Bluetooth PrefKeyServicesEnabled 2>/dev/null || echo "0")
if [ "$bt_sharing" = "1" ]; then
  bluetooth_sharing_disabled=false
fi

# Media Sharing.
media_sharing_disabled=true
media_sharing=$(defaults read com.apple.amp.mediasharingd home-sharing-enabled 2>/dev/null || echo "0")
if [ "$media_sharing" = "1" ]; then
  media_sharing_disabled=false
fi

# File Sharing (SMB).
file_sharing_disabled=true
if launchctl list com.apple.smbd >/dev/null 2>&1; then
  file_sharing_disabled=false
fi

# Content Caching.
content_caching_disabled=true
content_caching=$(defaults read /Library/Preferences/com.apple.AssetCache.plist Activated 2>/dev/null || echo "0")
if [ "$content_caching" = "1" ]; then
  content_caching_disabled=false
fi

# Current user is local admin.
user_is_admin=false
if groups 2>/dev/null | tr ' ' '\n' | grep -qx admin; then
  user_is_admin=true
fi

# Password required immediately after screen lock.
# Modern macOS (Ventura+) uses sysadminctl; pre-Ventura uses legacy defaults keys.
# Sentinels for password_lock_delay_seconds:
#   -1 = unknown (rule skipped)
#   -2 = screen lock disabled (rule fires with different evidence)
#    0 = immediate (rule passes)
#   >0 = delay in seconds (rule fires)
password_required_on_lock=false
password_lock_delay_seconds=-1
sl_out=$(sysadminctl -screenLock status 2>&1 || true)
if printf '%s' "$sl_out" | grep -qi "immediate"; then
  password_required_on_lock=true
  password_lock_delay_seconds=0
elif printf '%s' "$sl_out" | grep -qiE 'set to [0-9]+'; then
  n=$(printf '%s' "$sl_out" | grep -oE 'set to [0-9]+' | grep -oE '[0-9]+' | head -1)
  [ -n "$n" ] && password_lock_delay_seconds=$n
elif printf '%s' "$sl_out" | grep -qiE 'screenLock is OFF|disabled'; then
  password_lock_delay_seconds=-2
fi
# Fallback to legacy keys (pre-Ventura) when sysadminctl gave no answer.
if [ "$password_lock_delay_seconds" = "-1" ]; then
  ask_pw=$(defaults read com.apple.screensaver askForPassword 2>/dev/null | tr -dc '0-9')
  ask_pw_delay=$(defaults read com.apple.screensaver askForPasswordDelay 2>/dev/null | tr -dc '0-9')
  if [ -n "$ask_pw" ] && [ -n "$ask_pw_delay" ]; then
    password_lock_delay_seconds=$ask_pw_delay
    if [ "$ask_pw" = "1" ] && [ "$ask_pw_delay" = "0" ]; then
      password_required_on_lock=true
    fi
  fi
fi

# Time Machine destination encryption.
timemachine_destination_encrypted=true
if [ "$timemachine_enabled" = "true" ]; then
  if ! tmutil destinationinfo 2>/dev/null | grep -qi "Encrypted *: *Yes"; then
    timemachine_destination_encrypted=false
  fi
fi

# Time Machine recent backup (within 14 days).
timemachine_recent_backup=true
if [ "$timemachine_enabled" = "true" ]; then
  latest=$(tmutil latestbackup 2>/dev/null || echo "")
  if [ -z "$latest" ] || [ ! -e "$latest" ]; then
    timemachine_recent_backup=false
  else
    now=$(date +%s)
    mtime=$(stat -f '%m' "$latest" 2>/dev/null || stat -c '%Y' "$latest" 2>/dev/null || echo 0)
    age=$((now - mtime))
    if [ "$age" -gt 1209600 ]; then
      timemachine_recent_backup=false
    fi
  fi
fi

# AirPlay Receiver.
airplay_receiver_enabled=false
airplay_val=$(defaults -currentHost read com.apple.controlcenter AirplayRecieverEnabled 2>/dev/null || echo "")
[ "$airplay_val" = "1" ] && airplay_receiver_enabled=true

# Automatic download of OS updates.
os_auto_download_enabled=$(bool_check defaults read /Library/Preferences/com.apple.SoftwareUpdate AutomaticDownload "" "" "1")

# macOS major version.
macos_major=$(sw_vers -productVersion 2>/dev/null | cut -d. -f1)
macos_major=$(printf '%s' "$macos_major" | tr -dc '0-9')
[ -z "$macos_major" ] && macos_major=0

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
  "mdm_enrolled": $mdm,
  "timemachine_enabled": $timemachine_enabled,
  "remote_login_disabled": $remote_login_disabled,
  "remote_management_disabled": $remote_management_disabled,
  "bluetooth_sharing_disabled": $bluetooth_sharing_disabled,
  "media_sharing_disabled": $media_sharing_disabled,
  "file_sharing_disabled": $file_sharing_disabled,
  "content_caching_disabled": $content_caching_disabled,
  "user_is_admin": $user_is_admin,
  "password_required_on_lock": $password_required_on_lock,
  "password_lock_delay_seconds": $password_lock_delay_seconds,
  "timemachine_destination_encrypted": $timemachine_destination_encrypted,
  "timemachine_recent_backup": $timemachine_recent_backup,
  "airplay_receiver_enabled": $airplay_receiver_enabled,
  "os_auto_download_enabled": $os_auto_download_enabled,
  "macos_major_version": $macos_major
}
EOF
