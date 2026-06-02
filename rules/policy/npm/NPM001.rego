package obacht.npm

import rego.v1

findings contains f if {
	not input.configured
	f := {
		"rule_id": "NPM001",
		"evidence": "npm min-release-age is not configured (recommend >=7 days)",
	}
}

findings contains f if {
	input.configured
	input.value_seconds < 604800
	f := {
		"rule_id": "NPM001",
		"evidence": sprintf("npm min-release-age is %d seconds (<7 days)", [input.value_seconds]),
	}
}
