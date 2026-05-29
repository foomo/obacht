package obacht.claude

import rego.v1

findings contains f if {
	input.settings_disable_auto_mode != "disable"
	f := {
		"rule_id": "CLD024",
		"evidence": sprintf("disableAutoMode is %s, expected \"disable\"", [input.settings_disable_auto_mode]),
	}
}
