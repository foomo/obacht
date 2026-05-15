package obacht.privacy

import rego.v1

findings contains f if {
	not input.password_manager_installed
	f := {
		"rule_id": "PRV001",
		"evidence": "No password manager application detected in /Applications",
	}
}
