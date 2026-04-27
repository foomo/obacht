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

**Severity:** warn

DNS-over-HTTPS (DoH) or DNS-over-TLS (DoT) encrypts DNS queries, preventing network observers from seeing which domains you visit. Standard DNS queries are sent in plaintext.

**What it checks:**
- Whether `scutil --dns` output contains `dns_over_https` or `dns_over_tls` flags
- Whether DNS encryption tools like `dnscrypt-proxy` or `cloudflared` are running

**Remediation:**
```bash
# Install dnscrypt-proxy
brew install dnscrypt-proxy

# Or configure encrypted DNS in System Settings > Network > DNS
```

## PRV004: Untrusted DNS resolver is configured

**Severity:** warn

ISP- or router-provided DNS resolvers can log and monetize browsing activity, and may inject DNS-level redirects or telemetry. Pinning DNS to a known privacy-respecting provider eliminates the ISP/router as an observer (and pairs with PRV003 to also encrypt the queries).

**What it checks:**
- Configured nameservers from `scutil --dns` against an allowlist of trusted resolvers
- Allowlist: Cloudflare, Quad9, Google, AdGuard, Mullvad, NextDNS, OpenDNS, loopback (for `dnscrypt-proxy`/`cloudflared`)

**Recommended resolver IPs:**

| Provider   | IPv4                          | IPv6                                       |
|------------|-------------------------------|--------------------------------------------|
| Cloudflare | `1.1.1.1`, `1.0.0.1`          | `2606:4700:4700::1111`, `::1001`           |
| Cloudflare (malware-block) | `1.1.1.2`, `1.0.0.2` | — |
| Quad9      | `9.9.9.9`, `149.112.112.112`  | `2620:fe::fe`, `2620:fe::9`                |
| Google     | `8.8.8.8`, `8.8.4.4`          | `2001:4860:4860::8888`, `::8844`           |
| AdGuard    | `94.140.14.14`, `94.140.15.15` | —                                         |
| Mullvad    | `194.242.2.2`                 | —                                          |
| NextDNS    | `45.90.28.0`, `45.90.30.0`    | —                                          |
| OpenDNS    | `208.67.222.222`, `208.67.220.220` | —                                     |

**Remediation:**

System Settings > Network > Wi-Fi/Ethernet > Details > DNS — set primary/secondary to one of the IPs above. For DoH/DoT (also satisfies PRV003), install a config profile (e.g. https://one.one.one.one/dns/) or run `dnscrypt-proxy` locally and point DNS to `127.0.0.1`.
