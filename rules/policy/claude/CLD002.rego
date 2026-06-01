package obacht.claude

import rego.v1

findings contains f if {
	input.auto_compact_enabled != "false"
	f := {
		"rule_id": "CLD002",
		"evidence": sprintf("autoCompactEnabled is %s, expected false", [input.auto_compact_enabled]),
	}
}
