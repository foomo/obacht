package obacht.git

import rego.v1

findings contains f if {
	input.installed
	not input.signing_enabled
	f := {
		"rule_id": "GIT002",
		"evidence": "Git commit signing is not enabled",
	}
}
