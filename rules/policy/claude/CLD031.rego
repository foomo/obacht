package obacht.claude

import rego.v1

findings contains f if {
	input.sandbox_enabled != "true"
	f := {
		"rule_id": "CLD031",
		"evidence": sprintf("sandbox.enabled is %s, expected true", [input.sandbox_enabled]),
	}
}
