# Privacy Rules

These rules check for privacy and security tooling that protects against common threats like password reuse, network eavesdropping, and DNS snooping.

## PRV001: No password manager application detected

**Severity:** warn

A password manager helps generate and store strong, unique passwords for every account, which is a fundamental cybersecurity practice. Without one, users tend to reuse passwords across services, amplifying the impact of any single breach.

**What it checks:**
- Whether a recognized password manager is installed in `/Applications`
- Checks for 1Password, Bitwarden, Dashlane, KeePassXC, LastPass, and Enpass

**Remediation:**
```bash
# Install a password manager
brew install --cask 1password
# or
brew install --cask bitwarden
# or
brew install --cask keepassxc
```

## PRV002: No VPN configuration detected

**Severity:** info

A VPN encrypts network traffic and protects against eavesdropping on untrusted networks. This is especially important when working from public Wi-Fi or shared networks.

**What it checks:**
- Whether any macOS VPN configuration exists via `scutil --nc list`
- Whether common VPN processes are running (Tailscale, WireGuard, OpenVPN)

**Remediation:**

Configure a VPN in System Settings > Network > VPN, or install a VPN application.

## PRV003: Encrypted DNS is not configured

**Severity:** info

DNS-over-HTTPS (DoH) or DNS-over-TLS (DoT) encrypts DNS queries, preventing network observers from seeing which domains you visit. Standard DNS queries are sent in plaintext.

**What it checks:**
- Whether encrypted DNS is configured in macOS network settings
- Whether DNS encryption tools like `dnscrypt-proxy` or `cloudflared` are running

**Remediation:**
```bash
# Install dnscrypt-proxy
brew install dnscrypt-proxy

# Or configure encrypted DNS in System Settings > Network > DNS
```
