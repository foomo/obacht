package obacht.mise

import rego.v1

findings contains f if {
	not input.configured
	f := {
		"rule_id": "MIS001",
		"evidence": "mise minimum_release_age is not configured (recommend >=7 days)",
	}
}

findings contains f if {
	input.configured
	input.value_seconds < 604800
	f := {
		"rule_id": "MIS001",
		"evidence": sprintf("mise minimum_release_age is %d seconds (<7 days)", [input.value_seconds]),
	}
}
