package bouncer.os_test

import data.bouncer.os
import rego.v1

# Helper: a fully compliant macOS input
compliant_macos := {"os": {
	"os": "darwin",
	"arch": "arm64",
	"hostname": "myhost",
	"sip_enabled": true,
	"filevault_enabled": true,
	"firewall_enabled": true,
	"stealth_mode_enabled": true,
	"gatekeeper_enabled": true,
	"auto_login_disabled": true,
	"guest_account_disabled": true,
	"screen_lock_timeout_seconds": 300,
	"os_auto_update_enabled": true,
	"app_auto_update_enabled": true,
	"rsr_enabled": true,
	"screen_sharing_disabled": true,
	"internet_sharing_disabled": true,
	"printer_sharing_disabled": true,
	"remote_apple_events_disabled": true,
	"airdrop_setting": "off",
	"rosetta_installed": false,
	"edr_deployed": true,
	"legacy_kexts_blocked": true,
	"mdm_enrolled": true,
}}

# Fully compliant macOS should produce zero findings
test_compliant_macos if {
	findings := os.findings with input as compliant_macos
	count(findings) == 0
}

# Non-darwin OS should produce zero findings regardless of settings
test_non_darwin_no_findings if {
	findings := os.findings with input as {"os": {
		"os": "linux",
		"arch": "amd64",
		"hostname": "myhost",
	}}
	count(findings) == 0
}

# --- OS001: SIP ---

test_os001_sip_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"sip_enabled": false})})
	some f in findings
	f.rule_id == "OS001"
}

test_os001_sip_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS001")
}

# --- OS002: FileVault ---

test_os002_filevault_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"filevault_enabled": false})})
	some f in findings
	f.rule_id == "OS002"
}

test_os002_filevault_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS002")
}

# --- OS003: Firewall ---

test_os003_firewall_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"firewall_enabled": false})})
	some f in findings
	f.rule_id == "OS003"
}

test_os003_firewall_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS003")
}

# --- OS004: Stealth Mode ---

test_os004_stealth_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"stealth_mode_enabled": false})})
	some f in findings
	f.rule_id == "OS004"
}

test_os004_stealth_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS004")
}

# --- OS005: Gatekeeper ---

test_os005_gatekeeper_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"gatekeeper_enabled": false})})
	some f in findings
	f.rule_id == "OS005"
}

test_os005_gatekeeper_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS005")
}

# --- OS006: Auto Login ---

test_os006_auto_login_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"auto_login_disabled": false})})
	some f in findings
	f.rule_id == "OS006"
}

test_os006_auto_login_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS006")
}

# --- OS007: Guest Account ---

test_os007_guest_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"guest_account_disabled": false})})
	some f in findings
	f.rule_id == "OS007"
}

test_os007_guest_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS007")
}

# --- OS008: Screen Lock Timeout ---

test_os008_timeout_too_long if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"screen_lock_timeout_seconds": 600})})
	some f in findings
	f.rule_id == "OS008"
}

test_os008_timeout_ok if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS008")
}

test_os008_timeout_exact_boundary if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"screen_lock_timeout_seconds": 300})})
	not _has_rule(findings, "OS008")
}

# --- OS009: OS Auto Updates ---

test_os009_updates_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"os_auto_update_enabled": false})})
	some f in findings
	f.rule_id == "OS009"
}

test_os009_updates_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS009")
}

# --- OS010: App Auto Updates ---

test_os010_app_updates_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"app_auto_update_enabled": false})})
	some f in findings
	f.rule_id == "OS010"
}

test_os010_app_updates_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS010")
}

# --- OS011: RSR ---

test_os011_rsr_disabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"rsr_enabled": false})})
	some f in findings
	f.rule_id == "OS011"
}

test_os011_rsr_enabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS011")
}

# --- OS013: Screen Sharing ---

test_os013_screen_sharing_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"screen_sharing_disabled": false})})
	some f in findings
	f.rule_id == "OS013"
}

test_os013_screen_sharing_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS013")
}

# --- OS014: Internet Sharing ---

test_os014_internet_sharing_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"internet_sharing_disabled": false})})
	some f in findings
	f.rule_id == "OS014"
}

test_os014_internet_sharing_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS014")
}

# --- OS015: Printer Sharing ---

test_os015_printer_sharing_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"printer_sharing_disabled": false})})
	some f in findings
	f.rule_id == "OS015"
}

test_os015_printer_sharing_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS015")
}

# --- OS016: Remote Apple Events ---

test_os016_remote_events_enabled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"remote_apple_events_disabled": false})})
	some f in findings
	f.rule_id == "OS016"
}

test_os016_remote_events_disabled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS016")
}

# --- OS017: AirDrop Everyone ---

test_os017_airdrop_everyone if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"airdrop_setting": "everyone"})})
	some f in findings
	f.rule_id == "OS017"
}

test_os017_airdrop_contacts_only if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"airdrop_setting": "contacts_only"})})
	not _has_rule(findings, "OS017")
}

test_os017_airdrop_off if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS017")
}

# --- OS018: EDR ---

test_os018_no_edr if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"edr_deployed": false})})
	some f in findings
	f.rule_id == "OS018"
}

test_os018_edr_deployed if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS018")
}

# --- OS019: Legacy Kexts ---

test_os019_kexts_not_blocked if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"legacy_kexts_blocked": false})})
	some f in findings
	f.rule_id == "OS019"
}

test_os019_kexts_blocked if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS019")
}

# --- OS020: MDM ---

test_os020_not_enrolled if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"mdm_enrolled": false})})
	some f in findings
	f.rule_id == "OS020"
}

test_os020_enrolled if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS020")
}

# --- OS021: Rosetta 2 ---

test_os021_rosetta_installed if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"rosetta_installed": true})})
	some f in findings
	f.rule_id == "OS021"
}

test_os021_rosetta_not_installed if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS021")
}

# --- OS022: AirDrop advisory ---

test_os022_airdrop_contacts_only if {
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"airdrop_setting": "contacts_only"})})
	some f in findings
	f.rule_id == "OS022"
}

test_os022_airdrop_off if {
	findings := os.findings with input as compliant_macos
	not _has_rule(findings, "OS022")
}

test_os022_airdrop_everyone_no_advisory if {
	# OS017 fires for "everyone", OS022 should NOT also fire
	findings := os.findings with input as object.union(compliant_macos, {"os": object.union(compliant_macos.os, {"airdrop_setting": "everyone"})})
	not _has_rule(findings, "OS022")
}

# --- Helper ---

_has_rule(findings, rule_id) if {
	some f in findings
	f.rule_id == rule_id
}
