package obacht.yarn

import rego.v1

findings contains f if {
	not input.configured
	f := {
		"rule_id": "YRN001",
		"evidence": "yarn npmMinimalAgeGate is not configured (recommend >=7 days)",
	}
}

findings contains f if {
	input.configured
	input.value_seconds < 604800
	f := {
		"rule_id": "YRN001",
		"evidence": sprintf("yarn npmMinimalAgeGate is %d seconds (<7 days)", [input.value_seconds]),
	}
}
