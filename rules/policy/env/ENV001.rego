package obacht.env

import rego.v1

findings contains f if {
	some var in input.suspicious_vars
	f := {
		"rule_id": "ENV001",
		"evidence": sprintf("Suspicious env var: %s (matched pattern: %s)", [var.name, var.pattern]),
	}
}
