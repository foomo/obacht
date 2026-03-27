package bouncer.os

import rego.v1

# OS001: System Integrity Protection
findings contains f if {
	input.os.os == "darwin"
	not input.os.sip_enabled
	f := {"rule_id": "OS001", "evidence": "System Integrity Protection is disabled"}
}

# OS002: FileVault
findings contains f if {
	input.os.os == "darwin"
	not input.os.filevault_enabled
	f := {"rule_id": "OS002", "evidence": "FileVault disk encryption is disabled"}
}

# OS003: Application Firewall
findings contains f if {
	input.os.os == "darwin"
	not input.os.firewall_enabled
	f := {"rule_id": "OS003", "evidence": "Application Firewall is disabled"}
}

# OS004: Stealth Mode
findings contains f if {
	input.os.os == "darwin"
	not input.os.stealth_mode_enabled
	f := {"rule_id": "OS004", "evidence": "Stealth Mode is disabled"}
}

# OS005: Gatekeeper
findings contains f if {
	input.os.os == "darwin"
	not input.os.gatekeeper_enabled
	f := {"rule_id": "OS005", "evidence": "Gatekeeper is disabled"}
}

# OS006: Automatic Login
findings contains f if {
	input.os.os == "darwin"
	not input.os.auto_login_disabled
	f := {"rule_id": "OS006", "evidence": "Automatic login is enabled"}
}

# OS007: Guest Account
findings contains f if {
	input.os.os == "darwin"
	not input.os.guest_account_disabled
	f := {"rule_id": "OS007", "evidence": "Guest account is enabled"}
}

# OS008: Screen Lock Timeout
findings contains f if {
	input.os.os == "darwin"
	input.os.screen_lock_timeout_seconds > 300
	f := {
		"rule_id": "OS008",
		"evidence": sprintf("Screen lock timeout is %d seconds (maximum 300)", [input.os.screen_lock_timeout_seconds]),
	}
}

# OS009: OS Auto Updates
findings contains f if {
	input.os.os == "darwin"
	not input.os.os_auto_update_enabled
	f := {"rule_id": "OS009", "evidence": "Automatic OS updates are disabled"}
}

# OS010: App Auto Updates
findings contains f if {
	input.os.os == "darwin"
	not input.os.app_auto_update_enabled
	f := {"rule_id": "OS010", "evidence": "Automatic App Store updates are disabled"}
}

# OS011: Rapid Security Response
findings contains f if {
	input.os.os == "darwin"
	not input.os.rsr_enabled
	f := {"rule_id": "OS011", "evidence": "Rapid Security Responses are disabled"}
}

# OS013: Screen Sharing
findings contains f if {
	input.os.os == "darwin"
	not input.os.screen_sharing_disabled
	f := {"rule_id": "OS013", "evidence": "Screen Sharing is enabled"}
}

# OS014: Internet Sharing
findings contains f if {
	input.os.os == "darwin"
	not input.os.internet_sharing_disabled
	f := {"rule_id": "OS014", "evidence": "Internet Sharing is enabled"}
}

# OS015: Printer Sharing
findings contains f if {
	input.os.os == "darwin"
	not input.os.printer_sharing_disabled
	f := {"rule_id": "OS015", "evidence": "Printer Sharing is enabled"}
}

# OS016: Remote Apple Events
findings contains f if {
	input.os.os == "darwin"
	not input.os.remote_apple_events_disabled
	f := {"rule_id": "OS016", "evidence": "Remote Apple Events are enabled"}
}

# OS017: AirDrop
findings contains f if {
	input.os.os == "darwin"
	input.os.airdrop_setting == "everyone"
	f := {"rule_id": "OS017", "evidence": "AirDrop is set to Everyone"}
}

# OS018: EDR
findings contains f if {
	input.os.os == "darwin"
	not input.os.edr_deployed
	f := {"rule_id": "OS018", "evidence": "No Endpoint Detection & Response agent deployed"}
}

# OS019: Legacy Kernel Extensions
findings contains f if {
	input.os.os == "darwin"
	not input.os.legacy_kexts_blocked
	f := {"rule_id": "OS019", "evidence": "Legacy kernel extensions (kexts) are not blocked"}
}

# OS020: MDM Enrollment
findings contains f if {
	input.os.os == "darwin"
	not input.os.mdm_enrolled
	f := {"rule_id": "OS020", "evidence": "Device is not enrolled in MDM"}
}

# OS021: Rosetta 2 (advisory)
findings contains f if {
	input.os.os == "darwin"
	input.os.rosetta_installed
	f := {"rule_id": "OS021", "evidence": "Rosetta 2 is installed; remove if no longer needed to reduce attack surface"}
}

# OS022: AirDrop not off (advisory)
findings contains f if {
	input.os.os == "darwin"
	input.os.airdrop_setting != "off"
	input.os.airdrop_setting != "everyone"
	f := {"rule_id": "OS022", "evidence": sprintf("AirDrop is set to %s; consider disabling entirely", [input.os.airdrop_setting])}
}
