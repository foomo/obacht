package obacht.pnpm

import rego.v1

findings contains f if {
	not input.configured
	f := {
		"rule_id": "PNP001",
		"evidence": "pnpm minimumReleaseAge is not configured (recommend >=7 days)",
	}
}

findings contains f if {
	input.configured
	input.value_seconds < 604800
	f := {
		"rule_id": "PNP001",
		"evidence": sprintf("pnpm minimumReleaseAge is %d seconds (<7 days)", [input.value_seconds]),
	}
}
