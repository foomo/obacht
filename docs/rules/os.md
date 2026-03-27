# OS Rules

macOS-specific security configuration checks. These rules verify that critical operating system protections are enabled and sharing services are properly restricted.

## OS001: System Integrity Protection is disabled

**Severity:** critical

System Integrity Protection (SIP) prevents modification of protected system files and directories. Disabling SIP significantly weakens macOS security and allows malware to tamper with system components.

**What it checks:**
- Whether SIP is enabled via `csrutil status`

**Remediation:**
```bash
# Boot into Recovery Mode and run:
csrutil enable
```

## OS002: FileVault disk encryption is disabled

**Severity:** critical

FileVault provides full-disk encryption protecting data at rest. Without it, data on the disk can be accessed by removing or mounting the drive externally.

**What it checks:**
- Whether FileVault is enabled via `fdesetup status`

**Remediation:**

Enable FileVault in System Settings > Privacy & Security > FileVault.

## OS003: Application Firewall is disabled

**Severity:** high

The application firewall controls incoming network connections on a per-application basis, blocking unauthorized access attempts.

**What it checks:**
- Whether the macOS application firewall is enabled

**Remediation:**
```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on
```

## OS004: Stealth Mode is disabled

**Severity:** high

Stealth Mode prevents the Mac from responding to probing requests such as ICMP ping, making it harder to discover on the network.

**What it checks:**
- Whether Stealth Mode is enabled in the application firewall

**Remediation:**
```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setstealthmode on
```

## OS005: Gatekeeper is disabled

**Severity:** critical

Gatekeeper ensures only trusted software runs by verifying apps are signed by an identified developer or distributed via the App Store.

**What it checks:**
- Whether Gatekeeper is enabled via `spctl --status`

**Remediation:**
```bash
sudo spctl --master-enable
```

## OS006: Automatic login is enabled

**Severity:** high

Automatic login allows anyone with physical access to use the Mac without authentication.

**What it checks:**
- Whether `autoLoginUser` is set in login window preferences

**Remediation:**

Disable in System Settings > Users & Groups > Automatic login.

## OS007: Guest account is enabled

**Severity:** high

The guest account allows unauthenticated access to the Mac and can be used as an attack vector.

**What it checks:**
- Whether the Guest account is enabled in login window preferences

**Remediation:**

Disable in System Settings > Users & Groups > Guest User.

## OS008: Screen lock timeout exceeds 5 minutes

**Severity:** warn

A long screen lock timeout leaves the Mac accessible when unattended. The timeout should be 300 seconds (5 minutes) or less.

**What it checks:**
- The screensaver idle time setting
- Whether it exceeds 300 seconds

**Remediation:**

Set in System Settings > Lock Screen > Require password after screen saver begins.

## OS009: Automatic OS updates are disabled

**Severity:** high

Automatic OS updates ensure critical security patches are applied promptly without manual intervention.

**What it checks:**
- Whether `AutomaticallyInstallMacOSUpdates` is enabled

**Remediation:**
```bash
softwareupdate --schedule on
```

## OS010: Automatic App Store updates are disabled

**Severity:** warn

Automatic App Store updates keep applications patched against known vulnerabilities.

**What it checks:**
- Whether automatic App Store updates are enabled

**Remediation:**

Enable in System Settings > General > Software Update > Automatic Updates.

## OS011: Rapid Security Responses are disabled

**Severity:** high

Rapid Security Responses deliver urgent security fixes between regular OS updates, providing faster protection against active threats.

**What it checks:**
- Whether `ConfigDataInstall` is enabled in Software Update preferences

**Remediation:**

Enable in System Settings > General > Software Update > Automatic Updates.

## OS013: Screen Sharing is enabled

**Severity:** high

Screen Sharing allows remote access to the Mac desktop and should be disabled unless actively needed.

**What it checks:**
- Whether the `com.apple.screensharing` service is running

**Remediation:**
```bash
sudo launchctl disable system/com.apple.screensharing
```

## OS014: Internet Sharing is enabled

**Severity:** high

Internet Sharing turns the Mac into a network gateway, which can expose internal networks to unauthorized access.

**What it checks:**
- Whether Internet Sharing is enabled in network preferences

**Remediation:**

Disable in System Settings > General > Sharing > Internet Sharing.

## OS015: Printer Sharing is enabled

**Severity:** warn

Printer Sharing exposes the Mac on the network and should be disabled unless required.

**What it checks:**
- Whether CUPS printer sharing is enabled

**Remediation:**

Disable in System Settings > General > Sharing > Printer Sharing.

## OS016: Remote Apple Events are enabled

**Severity:** high

Remote Apple Events allow other computers to send Apple Events to this Mac, which can be used to execute commands remotely.

**What it checks:**
- Whether the `com.apple.AEServer` service is running

**Remediation:**

Disable in System Settings > General > Sharing > Remote Apple Events.

## OS017: AirDrop is set to Everyone

**Severity:** high

Setting AirDrop to Everyone allows any nearby device to send files, increasing the risk of social engineering attacks.

**What it checks:**
- The AirDrop discoverability setting

**Remediation:**

Set AirDrop to Contacts Only or Off in Finder > AirDrop.

## OS018: No EDR agent deployed

**Severity:** critical

An Endpoint Detection & Response agent provides real-time threat monitoring and incident response capabilities.

**What it checks:**
- Whether a recognized EDR agent is running (CrowdStrike, SentinelOne, Carbon Black, Microsoft Defender)

**Remediation:**

Install the organization-approved EDR agent.

## OS019: Legacy kernel extensions are not blocked

**Severity:** high

Legacy kernel extensions (kexts) run with full kernel privileges and should be replaced with System Extensions.

**What it checks:**
- Whether legacy kernel extensions are loaded via `kmutil`

**Remediation:**

Configure system extension policy via MDM or System Settings.

## OS020: Device is not enrolled in MDM

**Severity:** high

MDM enrollment enables centralized security policy enforcement, remote wipe, and compliance monitoring.

**What it checks:**
- Whether the device is enrolled in MDM via `profiles status`

**Remediation:**

Enroll the device via your organization's MDM solution.

## OS021: Rosetta 2 is installed

**Severity:** info

Rosetta 2 enables running Intel binaries on Apple Silicon. If no longer needed, removing it reduces attack surface.

**What it checks:**
- Whether the `oahd` (Rosetta) process is running

**Remediation:**

Remove Rosetta 2 if no Intel-only applications are required.

## OS022: AirDrop is not fully disabled

**Severity:** info

For high-risk roles, consider disabling AirDrop entirely rather than using Contacts Only.

**What it checks:**
- Whether AirDrop is set to any value other than Off

**Remediation:**

Set AirDrop to Off in Finder > AirDrop.

## OS023: Time Machine backup is disabled

**Severity:** warn

Time Machine provides automatic backups that protect against data loss from ransomware, hardware failure, or accidental deletion.

**What it checks:**
- Whether Time Machine is enabled with a configured backup destination

**Remediation:**

Enable Time Machine in System Settings > General > Time Machine.

## OS024: Remote Login (SSH server) is enabled

**Severity:** high

Remote Login runs an SSH server on the Mac, allowing remote shell access. This expands the attack surface and should be disabled unless actively needed.

**What it checks:**
- Whether the Remote Login (SSH) service is enabled via `systemsetup`

**Remediation:**

Disable in System Settings > General > Sharing > Remote Login.

## OS025: Remote Management is enabled

**Severity:** high

Remote Management allows remote control of the Mac via Apple Remote Desktop or VNC. This should be disabled unless required by IT policy.

**What it checks:**
- Whether the `com.apple.RemoteDesktop.agent` service is running

**Remediation:**

Disable in System Settings > General > Sharing > Remote Management.

## OS026: Bluetooth Sharing is enabled

**Severity:** warn

Bluetooth Sharing allows other devices to send files via Bluetooth. This should be disabled to reduce attack surface.

**What it checks:**
- Whether Bluetooth Sharing is enabled in Bluetooth preferences

**Remediation:**

Disable in System Settings > General > Sharing > Bluetooth Sharing.

## OS027: Media Sharing is enabled

**Severity:** warn

Media Sharing exposes media libraries on the local network and should be disabled to reduce attack surface.

**What it checks:**
- Whether home sharing is enabled in media sharing preferences

**Remediation:**

Disable in System Settings > General > Sharing > Media Sharing.

## OS028: File Sharing (SMB) is enabled

**Severity:** warn

File Sharing opens SMB network ports, allowing other devices to access shared folders. This should be disabled unless actively needed.

**What it checks:**
- Whether the `com.apple.smbd` service is running

**Remediation:**

Disable in System Settings > General > Sharing > File Sharing.

## OS029: Content Caching is enabled

**Severity:** warn

Content Caching shares downloaded Apple content with other devices on the network. On developer machines this is unnecessary and increases network exposure.

**What it checks:**
- Whether Content Caching is activated in system preferences

**Remediation:**

Disable in System Settings > General > Sharing > Content Caching.
