package obacht.claude

import rego.v1

findings contains f if {
	input.sandbox_fail_if_unavailable != "true"
	f := {
		"rule_id": "CLD005",
		"evidence": sprintf("sandbox.failIfUnavailable is %s, expected true", [input.sandbox_fail_if_unavailable]),
	}
}
