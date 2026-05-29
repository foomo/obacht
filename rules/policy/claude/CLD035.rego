package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_disable_bypass_mode != "disable"
	f := {
		"rule_id": "CLD035",
		"evidence": sprintf("permissions.disableBypassPermissionsMode is %s, expected \"disable\"", [input.permissions_disable_bypass_mode]),
	}
}
